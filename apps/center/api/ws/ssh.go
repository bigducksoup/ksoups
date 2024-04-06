package ws

import (
	"apps/center/api/ws/base"
	"apps/center/global"
	"apps/center/model"
	"apps/center/ssh"
	"bufio"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"log"
)

func DoSSH(client *base.Client, c *gin.Context) {

	id, a := c.GetQuery("sshInfoId")

	if !(a) {
		Ctx.DeRegChan <- client
		return
	}

	var sshInfo model.SSHInfo
	rows := global.DB.Where("id = ?", id).First(&sshInfo).RowsAffected

	if rows == 0 {
		Ctx.DeRegChan <- client
		return
	}

	// 创建ssh session
	session, err := ssh.NewSession(sshInfo.AddrPort, sshInfo.Username, sshInfo.Password)

	if err != nil {
		Ctx.DeRegChan <- client
		return
	}

	// 为shell准备标准输出和标准错误的 pipe
	or, ow := io.Pipe()
	er, ew := io.Pipe()

	// 开启一个shell
	w, err := session.OpenShell(ow, ew)

	if err != nil {
		Ctx.DeRegChan <- client
		return
	}

	// 用于控制client收发信息的context
	ctx, cancel := context.WithCancel(context.TODO())

	client.SetMessageHandleFunc(func(messageType int, bytes []byte, err error) {

		if err != nil || !(messageType == websocket.TextMessage || messageType == websocket.BinaryMessage) {
			log.Println("err occur !!!" + err.Error())
			cancel()
			Ctx.DeRegChan <- client
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
			Ctx.DeRegChan <- client
			session.Close()
			ow.Close()
			or.Close()
			er.Close()
			ew.Close()
		}

	})
	err = client.Setup(ctx)

	if err != nil {
		cancel()
		Ctx.DeRegChan <- client
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
		Ctx.DeRegChan <- client
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
		Ctx.DeRegChan <- client
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
