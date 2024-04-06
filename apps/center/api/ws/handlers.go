package ws

import (
	"apps/center/api/session"
	"apps/center/api/ws/base"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var AppMap = map[string]func(client *base.Client, c *gin.Context){
	"ssh": DoSSH,
	"msg": DoMessagePush,
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

	ip := c.RemoteIP()

	if ip != "127.0.0.1" {

		sid, ok := c.GetQuery("sid")

		if !ok {
			c.Status(http.StatusForbidden)
			return
		}

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
	client := base.NewClient(conn)

	// 注册到Context
	Ctx.RegChan <- client

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	f(client, c)

}
