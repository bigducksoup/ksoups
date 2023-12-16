package connect

import (
	"config-manager/probe/common"
	"config-manager/probe/handlers"
	"context"
	"errors"
	"io"
	"log"
	"time"
)

func InitConnect(addr string, ctx context.Context) {

	//create a connection
	connection := common.CreateConnection(addr)
	//goroutine for heartbeat
	go heartBeat(&connection, ctx)

	for true {
		select {
		case <-ctx.Done():
			log.Println("stop accept message")
			log.Println(ctx.Err())
			return

		default:
			msg, err := connection.Receive()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Println(err)
			}
			go handlers.HandleMessage(msg, &connection)
		}

	}

}

// do heart beat
// if beat did not receive response, try to reconnect
func heartBeat(connection *common.Connection, ctx context.Context) {

	//create a ticker
	ticker := time.NewTicker(3 * time.Second)

	//define a func for reconnecting
	tryReconnect := func(n int) {
		for i := 0; i < n; i++ {
			err := connection.Reconnect()
			if err == nil {
				break
			}
		}
		log.Fatal(errors.New("could not reconnect to center"))
	}

	//tick
	for true {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			//di ping
			ping := connection.Ping()
			if ping {
				continue
			}
			// if ping == false, try to reconnect
			tryReconnect(10)
		}
	}

}
