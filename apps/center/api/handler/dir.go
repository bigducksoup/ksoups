package handler

import (
	"apps/center/api/param"
	"apps/center/api/response"
	"apps/center/server/ServerContext"
	"apps/common/message"
	"apps/common/message/data"
	"encoding/json"
	"net/http"

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

	probe, err := ServerContext.Ctx.GetProbe(probeId)
	if err != nil {
		ctx.JSON(200, response.Fail(err))
		return
	}

	bytes, err := ServerContext.Ctx.SendMsgExpectRes(probe.Id, read, message.READDIR)

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

func DirCreate(c *gin.Context) {

	p := param.DirCreateParams{}

	err := c.ShouldBindJSON(&p)

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	probe, err := ServerContext.Ctx.GetProbe(p.ProbeId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	dirCreate := data.DirCreate{
		Path:       p.Path,
		Permission: p.Permission,
	}

	bytes, err := ServerContext.Ctx.SendMsgExpectRes(probe.Id, dirCreate, message.CREATE_DIR)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	resp := data.DirCreateResponse{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(resp))

}
