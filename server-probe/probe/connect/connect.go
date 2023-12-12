package connect

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

func InitConnect(addr string) {

	//create a connection
	connection := createConnection(addr)
	//goroutine for heartbeat
	go heartBeat(&connection)

	//receive message
	for true {
		msg, err := connection.Receive()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println(err)
		}
		HandleMessage(msg, &connection)
	}

}

// do heart beat
// if beat did not receive response, try to reconnect
func heartBeat(connection *Connection) {

	//create a ticker
	ticker := time.NewTicker(3 * time.Second)

	//define a func for reconnecting
	tryReconnect := func(n int) {
		for i := 0; i < n; i++ {
			err := connection.reconnect()
			if err == nil {
				break
			}
		}
		log.Fatal(errors.New("could not reconnect to center"))
	}

	//tick
	for true {
		select {
		case <-ticker.C:
			//di ping
			ping := connection.ping()
			if ping {
				continue
			}
			// if ping == false, try to reconnect
			tryReconnect(10)
		}
	}

}
