package ws

import (
	"apps/center/api/ws/base"
	"apps/center/global"
	"apps/center/model"
	"apps/center/ssh"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

func DoSSH(client *base.Client, c *gin.Context) {

	id, a := c.GetQuery("sshInfoId")

	if !(a) {
		return
	}

	// find ssh info detail
	var sshInfo model.SSHInfo
	rows := global.DB.Where("id = ?", id).First(&sshInfo).RowsAffected

	if rows == 0 {
		return
	}

	// 创建ssh session
	session, err := ssh.NewSession(sshInfo.AddrPort, sshInfo.Username, sshInfo.Password)

	if err != nil {
		return
	}

	// open shell
	err = session.OpenShell()
	if err != nil {
		return
	}

	client.Setup(func(messageType int, bytes []byte, err error) {
		if err != nil {
			session.Close()
			client.Close()
			return
		}

		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			session.Close()
			client.Close()
			return
		}

		n, err := (*session.StdinPipe).Write(bytes)

		log.Println(n)

		if err != nil {
			session.Close()
			client.Close()
		}
	})

	go func() {
		for {
			b := make([]byte, 1024)
			n, rErr := session.StderrPipe.Read(b)
			if rErr != nil {
				session.Close()
				client.Close()
				return
			}
			client.Send(websocket.BinaryMessage, b[:n])
		}
	}()

	go func() {
		for {
			b := make([]byte, 1024)
			n, rErr := session.StdoutPipe.Read(b)
			if rErr != nil {
				session.Close()
				client.Close()
				return
			}
			client.Send(websocket.BinaryMessage, b[:n])
		}
	}()

	session.Wait()

	log.Println("session bye")

}
