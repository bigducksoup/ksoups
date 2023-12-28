package server

import (
	"bufio"
	"config-manager/center/global"
	"config-manager/common/message"
	"config-manager/common/model/node"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

func Start(port string) {

	log.Println("center server initializing")

	//tcp 监听
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		panic(err)
	}

	log.Println("center listening port", port)

	go func() {
		for {
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

	//处理probe注册
	probe, err := handleRegister(conn, &Ctx, 20*time.Second)

	if err != nil {
		log.Println(err)
		(*conn).Close()
		return
	}

	reader := bufio.NewReader(*conn)

	for {

		//读取消息
		decode, err := message.DecodedToBytes(reader)

		//如果连接断开，删除探针
		if err == io.EOF || len(decode) == 0 {
			offLineNode(probe)
			break
		}

		response := message.Msg{}

		err = json.Unmarshal(decode, &response)
		if err != nil {
			log.Println(err)
			continue
		}

		err = handleMessage(response)
		if err != nil {
			log.Println(err)
			continue
		}

	}

}

var handlePolicy = map[message.Type]func(msg message.Msg){
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

	for {
		select {
		case <-ticker.C:
			//遍历所有探针
			Ctx.Probes.Range(func(key, value any) bool {

				probe := value.(*Probe)
				lastPingTime := probe.LastPingTime
				duration := time.Since(lastPingTime)

				//如果超过8秒，关闭探针
				if duration > 8*time.Second {
					offLineNode(probe)
				}

				return true
			})
		}
	}

}

func offLineNode(probe *Probe) {
	Ctx.RemoveProbe(probe.Id)
	global.DB.Model(&node.Node{}).Where("id = ?", probe.Id).Update("online", false)
}

// 将probe注册到ctx
func handleRegister(conn *net.Conn, serverCtx *Context, timeout time.Duration) (*Probe, error) {

	//超时定时器
	timer := time.NewTimer(timeout)

	msgChan := make(chan message.Msg)
	errChan := make(chan error)

	//开启协程等待注册消息
	go func() {
		reader := bufio.NewReader(*conn)
		decode, err := message.DecodedToBytes(reader)

		if err == io.EOF || len(decode) == 0 {
			errChan <- errors.New("connection lost")
			return
		}

		registerMsg := message.Msg{}

		err = json.Unmarshal(decode, &registerMsg)
		if err != nil {
			errChan <- err
			return
		}
		msgChan <- registerMsg
	}()

	select {
	case <-timer.C:
		//超时
		return nil, errors.New("register time out,no message send")
	case err := <-errChan:
		//发生错误
		return nil, err
	case msg := <-msgChan:
		//创建探针
		probeId := string(msg.Data)
		_, err := serverCtx.GetProbe(probeId)

		//err == nil 说明已存在id = probeId
		if err == nil {
			m := message.Msg{Id: msg.Id, ErrMark: true, Data: []byte("There is already a probe with id :" + probeId)}
			bytes, _ := json.Marshal(m)
			encode, _ := message.Encode(bytes)
			(*conn).Write(encode)
			return nil, errors.New("There is already a probe with id :" + probeId)
		}

		//创建probe
		probe := CreateProbe(conn, probeId)
		log.Printf("handling connection, addr:%s, id:%s \n", probe.Addr, probe.Id)

		//添加探针
		serverCtx.AddProbe(probe.Id, &probe)
		var ct int64
		global.DB.Model(&node.Node{}).Where("id = ?", probeId).Count(&ct)

		if ct == 0 {
			n := node.Node{
				Id:      probeId,
				Name:    probeId,
				Address: probe.Addr,
				RegTime: time.Now(),
				Online:  true,
			}

			global.DB.Create(&n)
		} else {
			global.DB.Model(&node.Node{}).Where("id = ?", probeId).Updates(map[string]any{
				"address": probe.Addr,
				"online":  1,
			})
		}

		//返回注册成功
		err = serverCtx.SendMsg(probe.Id, message.Msg{
			Id:      msg.Id,
			ErrMark: false,
		})

		if err != nil {
			return nil, err
		}

		return &probe, nil
	}
}
