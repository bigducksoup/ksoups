package api

import (
	"config-manager/center/apiserver/params"
	"config-manager/center/apiserver/response"
	"config-manager/center/server"
	"config-manager/common/message"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func FileRead(c *gin.Context) {

	path, pathOk := c.GetQuery("path")
	address, addrOk := c.GetQuery("address")

	if !addrOk || !pathOk {
		c.JSON(200, response.ParamsError())
		return
	}

	probe, err := server.Ctx.GetProbeByAddr(address)
	if err != nil {
		c.JSON(200, response.StringFail(err.Error()))
		return
	}

	read := message.FileRead{
		Path: path,
	}

	bytes, err := server.Ctx.SendMsgExpectRes(probe.Id, read, message.READFILE)

	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	resp := message.FileReadResponse{}

	err = json.Unmarshal(bytes, &resp)

	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	c.JSON(200, response.Success[message.FileReadResponse](resp))

}

func FileModify(c *gin.Context) {

	param := params.ModifyFileParams{}

	err := c.ShouldBindJSON(&param)

	if err != nil {
		c.JSON(200, response.Fail(err))
	}

	probe, err := server.Ctx.GetProbeByAddr(param.Addr)

	if err != nil {
		c.JSON(200, response.Fail(err))
		return
	}

	fileMReq := message.FileModify{
		Path:    param.Path,
		Changes: []message.Change{},
	}

	for _, v := range param.Changes {
		fileMReq.Changes = append(fileMReq.Changes, message.Change{
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

	resp := message.FileModifyResponse{}

	json.Unmarshal(bytes, &resp)

	c.JSON(200, response.Success[message.FileModifyResponse](resp))

}
