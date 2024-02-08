package handler

import (
	"apps/center/api/param"
	"apps/center/api/response"
	"apps/center/server"
	"apps/common/message"
	"apps/common/message/data"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FileRead(c *gin.Context) {

	path, pathOk := c.GetQuery("path")
	id, ok := c.GetQuery("probeId")

	if !ok || !pathOk {
		c.JSON(200, response.ParamsError())
		return
	}

	probe, err := server.Ctx.GetProbe(id)
	if err != nil {
		c.JSON(200, response.StringFail(err.Error()))
		return
	}

	read := data.FileRead{
		Path: path,
	}

	bytes, err := server.Ctx.SendMsgExpectRes(probe.Id, read, message.READFILE)

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

	param := param.ModifyFileParams{}

	err := c.ShouldBindJSON(&param)

	if err != nil {
		c.JSON(200, response.Fail(err))
	}

	probe, err := server.Ctx.GetProbe(param.ProbeId)

	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	fileMReq := data.FileModify{
		Path:    param.Path,
		Changes: []data.Change{},
	}

	for _, v := range param.Changes {
		fileMReq.Changes = append(fileMReq.Changes, data.Change{
			Count:     v.Count,
			Operation: v.Operation,
			Value:     v.Value,
		})
	}

	bytes, err := server.Ctx.SendMsgExpectRes(probe.Id, fileMReq, message.MODIFYFILE)

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

	probe, err := server.Ctx.GetProbe(fileCreateParams.ProbeId)

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

	bytes, err := server.Ctx.SendMsgExpectRes(probe.Id, fileCreate, message.CREATEFILE)
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
