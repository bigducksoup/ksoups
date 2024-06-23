package ws

import (
	"apps/center/api/ws/base"
	"apps/center/service"
	"apps/center/ssh"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"log"
)

func DoSSH(client *base.Client, c *gin.Context) {

	id, a := c.GetQuery("sshInfoId")

	if !(a) {
		return
	}

	sshInfo, ok := service.SSHCRUD.GetSSHInfo(id)

	if !ok {
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

		_, err = (*session.StdinPipe).Write(bytes)

		if err != nil {
			session.Close()
			client.Close()
		}
	})

	handle := func(reader *io.PipeReader) {
		for {
			b := make([]byte, 1024)
			n, rErr := reader.Read(b)
			if rErr != nil {
				session.Close()
				client.Close()
				return
			}
			client.Send(websocket.BinaryMessage, b[:n])
		}
	}

	go handle(session.StderrPipe)
	go handle(session.StdoutPipe)

	session.Wait()

	log.Println("session bye")

}
