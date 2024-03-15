package ws

import (
	"apps/center/api/session"
	"apps/center/global"
	"apps/center/model"
	"apps/center/ssh"
	"bufio"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

var AppMap = map[string]func(client *Client, c *gin.Context){
	"ssh": DoSSH,
}

var updater = websocket.Upgrader{
	ReadBufferSize:   2048,
	WriteBufferSize:  2048,
	HandshakeTimeout: 100 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWS(c *gin.Context) {

	sid, ok := c.GetQuery("sid")

	if !ok {
		c.Status(http.StatusForbidden)
		return
	}

	ip := c.RemoteIP()

	if ip != "127.0.0.1" {
		//检查是否存在对应 session
		_, ok = session.GetSession(sid)

		if !ok {
			c.Status(http.StatusForbidden)
			return
		}
	}

	app := c.Param("app")

	//寻找对应策略
	f, ok := AppMap[app]

	if !ok {
		c.Status(http.StatusForbidden)
		return
	}

	// 处理请求升级为websocket
	conn, err := updater.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// 创建客户端
	client := NewClient(conn)

	// 注册到Context
	Ctx.regChan <- client

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	f(client, c)

}

func DoSSH(client *Client, c *gin.Context) {

	id, a := c.GetQuery("sshInfoId")

	if !(a) {
		Ctx.deRegChan <- client
		return
	}

	var sshInfo model.SSHInfo
	rows := global.DB.Where("id = ?", id).First(&sshInfo).RowsAffected

	if rows == 0 {
		Ctx.deRegChan <- client
		return
	}

	// 创建ssh session
	session, err := ssh.NewSession(sshInfo.AddrPort, sshInfo.Username, sshInfo.Password)

	if err != nil {
		Ctx.deRegChan <- client
		return
	}

	// 为shell准备标准输出和标准错误的 pipe
	or, ow := io.Pipe()
	er, ew := io.Pipe()

	// 开启一个shell
	w, err := session.OpenShell(ow, ew)

	if err != nil {
		Ctx.deRegChan <- client
		return
	}

	// 用于控制client收发信息的context
	ctx, cancel := context.WithCancel(context.TODO())

	client.SetMessageHandleFunc(func(messageType int, bytes []byte, err error) {

		if err != nil || !(messageType == websocket.TextMessage || messageType == websocket.BinaryMessage) {
			log.Println("err occur !!!" + err.Error())
			cancel()
			Ctx.deRegChan <- client
			session.Close()
			ow.Close()
			or.Close()
			er.Close()
			ew.Close()
		}

		_, err = (*w).Write(bytes)
		if err != nil {
			log.Println("err occur !!!" + err.Error())
			cancel()
			Ctx.deRegChan <- client
			session.Close()
			ow.Close()
			or.Close()
			er.Close()
			ew.Close()
		}

	})
	err = client.setup(ctx)

	if err != nil {
		cancel()
		Ctx.deRegChan <- client
		session.Close()
		ow.Close()
		or.Close()
		er.Close()
		ew.Close()
		log.Println(err)
	}

	go handleOut(ctx, or, func(b []byte) {
		client.Send(websocket.BinaryMessage, b)
	}, func(err error) {
		cancel()
		session.Close()
		Ctx.deRegChan <- client
		ow.Close()
		or.Close()
		er.Close()
		ew.Close()
	})

	handleOut(ctx, er, func(b []byte) {
		client.Send(websocket.BinaryMessage, b)
	}, func(err error) {
		cancel()
		session.Close()
		Ctx.deRegChan <- client
		ow.Close()
		or.Close()
		er.Close()
		ew.Close()
	})

}

func handleOut(ctx context.Context, or *io.PipeReader, f func(b []byte), eh func(err error)) {

	reader := bufio.NewReader(or)

	for {
		select {
		case <-ctx.Done():
			or.Close()
			return
		default:

			b := make([]byte, 1024)
			n, err := reader.Read(b)

			if err != nil {
				eh(err)
				break
			}
			f(b[:n])
		}
	}

}

func DataPush(client *Client, c *gin.Context) {
	//TODO push data to client, data may be probe performance data, interactive shortcut output, etc.
}
