package server

import (
	"apps/center/global"
	"apps/center/model"
	"apps/center/server/ServerContext"
	"apps/center/service"
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/utils"
	"bufio"

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
			//go scanDeadProbe()

		}
	}()

}

// 处理连接
func handleConn(conn *net.Conn) {

	//处理probe注册
	probe, err := handleRegister(conn, &ServerContext.Ctx, 20*time.Second)

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
			ServerContext.Ctx.Probes.Range(func(key, value any) bool {

				probe := value.(*ServerContext.Probe)
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

func offLineNode(probe *ServerContext.Probe) {
	ServerContext.Ctx.RemoveProbe(probe.Id)
	global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probe.Id).Update("online", false)
}

func checkRegisterInfo(info data.RegisterInfo) error {

	//1. find private key
	//2. try to decrypt registerInfo.encryptedName
	//3. check whether registerInfo.Name == decryptedData

	registerKey, err := service.CENTER_INFO.GetKeyPairByPublicKeyMd5(info.PublicKeyMd5)

	if err != nil {
		return err
	}

	privateKeyBytes, err := utils.DecodeBase64ToKey(registerKey.PrivateKey)
	privateKey, err := utils.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return err
	}

	decryptData, err := utils.DecryptData([]byte(info.EncryptedName), privateKey)
	if err != nil {
		return err
	}

	if string(decryptData) != info.Name {
		return errors.New("register invalid")
	}

	return nil

}

// 将probe注册到ctx
func handleRegister(conn *net.Conn, serverCtx *ServerContext.Context, timeout time.Duration) (*ServerContext.Probe, error) {

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

		registerInfo := data.RegisterInfo{}

		err := json.Unmarshal(msg.Data, &registerInfo)

		if err != nil {
			return nil, err
		}

		err = checkRegisterInfo(registerInfo)

		if err != nil {
			return nil, err
		}

		//创建探针
		probeId := registerInfo.Name
		_, err = serverCtx.GetProbe(probeId)

		//err == nil 说明已存在id = probeId
		if err == nil {
			m := message.Msg{Id: msg.Id, ErrMark: true, Data: []byte("There is already a probe with id :" + probeId)}
			bytes, _ := json.Marshal(m)
			encode, _ := message.Encode(bytes)
			(*conn).Write(encode)
			return nil, errors.New("There is already a probe with id :" + probeId)
		}

		//创建probe
		probe := ServerContext.CreateProbe(conn, probeId)
		log.Printf("handling connection, addr:%s, id:%s \n", probe.Addr, probe.Id)

		//添加探针
		serverCtx.AddProbe(probe.Id, &probe)
		var ct int64
		global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probeId).Count(&ct)

		if ct == 0 {
			n := model.ProbeInfo{
				Id:      probeId,
				Name:    probeId,
				Address: probe.Addr,
				RegTime: time.Now(),
				Online:  true,
			}

			global.DB.Create(&n)
		} else {
			global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probeId).Updates(map[string]any{
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
