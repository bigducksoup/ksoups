package api

import (
	"config-manager/center/apiserver/response"
	"config-manager/center/server"
	"config-manager/common/message"
	"config-manager/common/message/data"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func DirRead(ctx *gin.Context) {

	path, pok := ctx.GetQuery("path")

	probeId, aok := ctx.GetQuery("probeId")

	if !aok || !pok {
		ctx.JSON(200, response.ParamsError())
		return
	}

	fileOnly := ctx.DefaultQuery("fileOnly", "false")

	read := data.DirRead{
		Path:     path,
		FileOnly: fileOnly == "true",
	}

	probe, err := server.Ctx.GetProbe(probeId)
	if err != nil {
		ctx.JSON(200, response.Fail(err))
		return
	}

	bytes, err := server.Ctx.SendMsgExpectRes(probe.Id, read, message.READDIR)

	if err != nil {
		ctx.JSON(200, response.Fail(err))
		return
	}

	resp := data.DirResponse{}

	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		ctx.JSON(200, response.Fail(err))
		return
	}

	ctx.JSON(200, response.Success[data.DirResponse](resp))

}
