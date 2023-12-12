package server

import (
	"bufio"
	"config-manager/core/message"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

func init() {

	fmt.Println("center server initializing")

	//tcp 监听
	listener, err := net.Listen("tcp", ":9999")

	if err != nil {
		panic(err)
	}

	go func() {
		for true {
			conn, accErr := listener.Accept()
			if accErr != nil {
				continue
			}
			go handleConn(&conn)
			go scanDeadProbe()

		}
	}()

}

// 处理连接
func handleConn(conn *net.Conn) {

	//创建探针
	probe := CreateProbe(conn)

	fmt.Printf("handling connection, addr:%s, id:%s \n", probe.Addr, probe.Id)

	//添加探针
	Ctx.AddProbe(probe.Id, &probe)

	reader := bufio.NewReader(*conn)

	for true {

		//读取消息
		decode, err := message.DecodedToBytes(reader)

		//如果连接断开，删除探针
		if err == io.EOF || len(decode) == 0 {
			Ctx.RemoveProbe(probe.Id)
			break
		}

		response := message.Msg{}

		err = json.Unmarshal(decode, &response)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = handleMessage(response)
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

}

var handlePolicy map[message.Type]func(msg message.Msg) = map[message.Type]func(msg message.Msg){
	message.RESPONSE:  HandleRESPONSE,
	message.HEARTBEAT: HandleHEARTBEAT,
}

// 处理message.Msg
func handleMessage(msg message.Msg) error {

	policy, ok := handlePolicy[msg.Type]

	if !ok {
		return errors.New("unknown message type , could not find handle policy")
	}

	policy(msg)

	return nil
}

// 扫描心跳消失的探针，超过8秒则关闭探针
func scanDeadProbe() {

	//定时器
	ticker := time.NewTicker(8 * time.Second)

	for true {
		select {
		case <-ticker.C:
			//遍历所有探针
			Ctx.Probes.Range(func(key, value any) bool {

				probe := value.(*Probe)
				lastPingTime := probe.LastPingTime
				duration := time.Since(lastPingTime)

				//如果超过8秒，关闭探针
				if duration > 8*time.Second {
					Ctx.RemoveProbe(probe.Id)
				}

				return true
			})
		}
	}

}
