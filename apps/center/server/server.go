package server

import (
	"apps/center/global"
	"apps/center/model"
	"apps/center/server/core"
	"apps/center/service"
	"apps/common/message"
	"apps/common/message/data"
	"apps/common/utils"
	"context"
	"errors"
	"time"
)

func MakeServer(options core.CenterServerOptions) *core.CenterServer {
	s := core.CreateCenterServer(options)
	return s
}

func InitCenterServer(port int, ctx context.Context) *core.CenterServer {
	centerServer := MakeServer(core.CenterServerOptions{
		RegisterTimeOut: 30 * time.Second,
		ResponseTimeOut: 10 * time.Minute,
		Port:            port,
		Context:         ctx,
		HandlePolicy: map[message.Type]func(msg message.Msg, serverContext *core.Context) error{
			message.RESPONSE:       HandleResponse,
			message.PROACTIVE_PUSH: HandleProActivePush,
		},
	})

	centerServer.SetAuthenticateMethod(func(info data.RegisterInfo) error {
		//1. find private key
		//2. try to decrypt registerInfo.encryptedName
		//3. check whether registerInfo.Name == decryptedData

		registerKey, err := service.CENTER_INFO.GetKeyPairByPublicKeyMd5(info.PublicKeyMd5)

		if err != nil {
			return err
		}

		privateKeyBytes, err := utils.DecodeBase64ToKey(registerKey.PrivateKey)

		if err != nil {
			return err
		}

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

	})

	centerServer.SetOnProbeRegister(func(probe *core.Probe, info data.RegisterInfo) {
		var ct int64
		global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probe.Id).Count(&ct)

		if ct == 0 {
			registerKey, _ := service.CENTER_INFO.GetKeyPairByPublicKeyMd5(info.PublicKeyMd5)

			n := model.ProbeInfo{
				Id:      probe.Id,
				Name:    probe.Id,
				Address: probe.Addr,
				RegTime: time.Now(),
				Online:  true,
				KeyId:   registerKey.ID,
			}

			global.DB.Create(&n)
		} else {
			global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probe.Id).Updates(map[string]any{
				"address": probe.Addr,
				"online":  1,
			})
		}
	})

	centerServer.SetOnProbeOffLine(func(probe *core.Probe) {
		global.DB.Model(&model.ProbeInfo{}).Where("id = ?", probe.Id).Update("online", false)
	})

	err := centerServer.Start()

	if err != nil {
		panic(err)
	}

	return centerServer
}
