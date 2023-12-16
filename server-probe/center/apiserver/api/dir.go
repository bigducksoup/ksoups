package api

import (
	"config-manager/center/apiserver/response"
	"config-manager/center/server"
	"config-manager/common/message"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func DirRead(ctx *gin.Context) {

	path, pok := ctx.GetQuery("path")

	addr, aok := ctx.GetQuery("address")

	if !aok || !pok {
		ctx.JSON(200, response.ParamsError())
		return
	}

	fileOnly := ctx.DefaultQuery("fileOnly", "false")

	read := message.DirRead{
		Path:     path,
		FileOnly: fileOnly == "true",
	}

	probe, err := server.Ctx.GetProbeByAddr(addr)
	if err != nil {
		ctx.JSON(200, response.Fail(err))
		return
	}

	bytes, err := server.Ctx.SendMsgExpectRes(probe.Id, read, message.READDIR)

	if err != nil {
		ctx.JSON(200, response.Fail(err))
		return
	}

	resp := message.DirResponse{}

	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		ctx.JSON(200, response.Fail(err))
		return
	}

	ctx.JSON(200, response.Success[message.DirResponse](resp))

}
