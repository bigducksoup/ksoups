package api

import (
	"config-manager/center/server"
	"config-manager/core/message"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func FileRead(c *gin.Context) {

	path, pathOk := c.GetQuery("path")
	address, addrOk := c.GetQuery("address")

	if !addrOk || !pathOk {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	probe, err := server.Ctx.GetProbeByAddr(address)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "addr不存在",
		})
		return
	}

	read := message.FileRead{
		Path: path,
	}

	bytes, err := server.Ctx.SendMsgExpectRes(probe.Id, read, message.READFILE)

	if err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	resp := message.FileReadResponse{}

	err = json.Unmarshal(bytes, &resp)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    resp,
	})

}
