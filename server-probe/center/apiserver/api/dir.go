package api

import (
	"config-manager/center/server"
	"config-manager/core/message"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func DirRead(ctx *gin.Context) {

	path, pok := ctx.GetQuery("path")

	addr, aok := ctx.GetQuery("address")

	if !aok || !pok {
		ctx.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	fileOnly := ctx.DefaultQuery("fileOnly", "false")

	read := message.DirRead{
		Path:     path,
		FileOnly: fileOnly == "true",
	}

	probe, err := server.Ctx.GetProbeByAddr(addr)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": "addr不存在",
		})
		return
	}

	bytes, err := server.Ctx.SendMsgExpectRes(probe.Id, read, message.READDIR)

	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err,
		})
		return
	}

	response := message.DirResponse{}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    500,
			"message": err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})

}
