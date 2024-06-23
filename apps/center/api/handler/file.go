package handler

import (
	"apps/center/api/param"
	"apps/center/api/response"
	"apps/center/global"
	"apps/common/message"
	"apps/common/message/data"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileRead(c *gin.Context) {

	path, pathOk := c.GetQuery("path")
	id, ok := c.GetQuery("probeId")

	if !ok || !pathOk {
		c.JSON(200, response.ParamsError())
		return
	}

	probe, err := global.CenterServer.Ctx.GetProbe(id)
	if err != nil {
		c.JSON(200, response.StringFail(err.Error()))
		return
	}

	read := data.FileRead{
		Path: path,
	}

	bytes, err := global.CenterServer.Ctx.SendMsgExpectRes(probe.Id, read, message.READ_FILE)

	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	resp := data.FileReadResponse{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	c.JSON(200, response.Success[data.FileReadResponse](resp))

}

func FileModify(c *gin.Context) {

	p := param.ModifyFileParams{}

	err := c.ShouldBindJSON(&p)

	if err != nil {
		c.JSON(200, response.Fail(err))
	}

	probe, err := global.CenterServer.Ctx.GetProbe(p.ProbeId)

	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	fileMReq := data.FileModify{
		Path:    p.Path,
		Changes: []data.Change{},
	}

	for _, v := range p.Changes {
		fileMReq.Changes = append(fileMReq.Changes, data.Change{
			Count:     v.Count,
			Operation: v.Operation,
			Value:     v.Value,
		})
	}

	bytes, err := global.CenterServer.Ctx.SendMsgExpectRes(probe.Id, fileMReq, message.MODIFY_FILE)

	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	resp := data.FileModifyResponse{}

	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	c.JSON(200, response.Success[data.FileModifyResponse](resp))

}

func FileCreate(c *gin.Context) {

	fileCreateParams := param.FileCreateParams{}

	err := c.ShouldBindJSON(&fileCreateParams)

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	probe, err := global.CenterServer.Ctx.GetProbe(fileCreateParams.ProbeId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	fileCreate := data.FileCreate{
		Path:       fileCreateParams.Path,
		Permission: fileCreateParams.Permission,
	}

	bytes, err := global.CenterServer.Ctx.SendMsgExpectRes(probe.Id, fileCreate, message.CREATE_FILE)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	resp := data.FileCreateResponse{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success[data.FileCreateResponse](resp))

}
