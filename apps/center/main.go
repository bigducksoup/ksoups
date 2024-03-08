package main

import (
	"apps/center/api"
	"apps/center/config"
	"apps/center/db"
	"apps/center/global"
	"apps/center/model"
	"apps/center/server"
	"apps/center/server/core"
	"apps/center/service"
	"apps/common/message/data"
	"apps/common/utils"
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	path := flag.String("c", "./center/center.yaml", "config file path")
	flag.Parse()

	//解析配置文件
	conf, err := config.LoadConf(*path)
	if err != nil {
		panic(err)
	}
	global.Conf = *conf

	hp := global.Conf.WebApi.Port
	sp := global.Conf.Center.Port
	db.InitDB()

	ctx, cancel := context.WithCancel(context.Background())

	InitCenterServer(sp, context.WithoutCancel(ctx))

	service.Init()
	api.InitApiServer(hp, context.WithoutCancel(ctx))

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown everything ...")
	cancel()

}

func InitCenterServer(port int, ctx context.Context) {
	centerServer := server.MakeServer(core.CenterServerOptions{
		RegisterTimeOut: 30 * time.Second,
		ResponseTimeOut: 10 * time.Minute,
		Port:            port,
		Context:         ctx,
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
	global.CenterServer = centerServer
}
