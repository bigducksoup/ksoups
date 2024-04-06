package base

import (
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/utils"
	"apps/probe/common"
	"apps/probe/config"
	"apps/probe/connect"
	"apps/probe/handlers"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"
	"time"
)

func InitConnect(addr string, ctx context.Context) {
	InitInfo()

	//create a connection
	connection := common.CreateConnection(addr)

	//注册到 center
	err := register(&connection)

	if err != nil {
		log.Fatal(err)
	}

	//notifyChan := make(chan any, 1)

	//goroutine for heartbeat
	//go heartBeat(&connection, core, &notifyChan)

	connect.Connection = &connection

	for {
		select {
		case <-ctx.Done():
			log.Println("stop accept message")
			log.Println(ctx.Err())
			return

		default:
			msg, err := connection.Receive()
			if err == io.EOF {
				flag := false
				log.Println("connection lost")
				for i := 0; i < 100; i++ {
					time.Sleep(5 * time.Second)
					log.Printf("trying reconnected to center,%d time", i)
					err := connection.Reconnect()
					if err != nil {
						continue
					}
					err = register(&connection)
					if err != nil {
						continue
					}
					flag = true
					break
				}

				if !flag {
					panic(errors.New("failed to reconnected to center after 100 times retry"))
				}
				continue
				//notify := <-notifyChan
				//
				//if notify == 1 {
				//	continue
				//}
				//
				//if notify == 0 {
				//	log.Fatal(errors.New("connection lost,reconnect failed after 10 times retry"))
				//}
			}
			if err != nil {
				log.Println(err)
				continue
			}
			go handlers.HandleMessage(msg, &connection)
		}

	}

}

// do heart beat
// if beat did not receive response, try to reconnect
func heartBeat(connection *common.Connection, ctx context.Context, notifyChan *chan any) {

	//create a ticker
	ticker := time.NewTicker(3 * time.Second)

	//define a func for reconnecting
	tryReconnect := func(n int) {

		for i := 0; i < n; i++ {
			time.Sleep(5 * time.Second)
			err := connection.Reconnect()
			if err != nil {
				continue
			}
			err = register(connection)
			if err != nil {
				continue
			}

			*notifyChan <- 1
			return

		}
		*notifyChan <- 0
		log.Println("reconnected to center failed")
	}

	//tick
	for {
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

var registerInfo data.RegisterInfo

func register(connection *common.Connection) error {

	if config.Conf.Name == "" {
		config.Conf.Name = strings.Split(connection.LocalAddr, ":")[0]
	}

	marshaledRegisterInfo, err := json.Marshal(registerInfo)

	if err != nil {
		panic(err)
	}

	msg := message.Msg{
		Id:       utils.UUID(),
		Type:     message.REGISTER,
		Data:     marshaledRegisterInfo,
		DataType: message.DEFAULT,
		ErrMark:  false,
	}

	err = connection.SendMessage(msg)

	if err != nil {
		return err
	}

	//超时控制
	ctx, cancelFunc := context.WithTimeout(context.TODO(), 20*time.Second)

	defer cancelFunc()
	resp, err := connection.ReceiveWithCtx(ctx)

	if err != nil {
		return err
	}

	if resp.Id == msg.Id && resp.ErrMark == false {
		return nil
	}

	return errors.New(string(resp.Data))

}

func InitInfo() {
	bytes, err := utils.DecodeBase64ToKey(config.Conf.PublicKey)

	if err != nil {
		panic(err)
	}

	publicKey, err := utils.ParsePublicKey(bytes)

	if err != nil {
		panic(err)
	}

	encryptData, err := utils.EncryptData([]byte(config.Conf.Name), publicKey)

	if err != nil {
		panic(err)
	}

	md5 := utils.Md5([]byte(config.Conf.PublicKey))

	registerInfo = data.RegisterInfo{
		Name:          config.Conf.Name,
		EncryptedName: encryptData,
		PublicKeyMd5:  md5,
	}

}
