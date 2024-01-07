package handler

import (
	"config-manager/center/model"
	"config-manager/center/service"
	chainService "config-manager/center/service/chain"
	"config-manager/center/webapi/param"
	"config-manager/center/webapi/response"
	"config-manager/common/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// ChainCreate 创建链
func ChainCreate(c *gin.Context) {

	p := param.ChainCreateParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
	}

	ch := model.Chain{
		Id:          utils.UUID(),
		Name:        p.Name,
		CreateTime:  time.Now(),
		Description: p.Description,
	}

	err = service.ChainCRUD.SaveChain(&ch)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
	}

	c.JSON(http.StatusOK, response.SuccessWithNoData())
}

func NodeCreate(c *gin.Context) {

	p := param.NodeCreateParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	node := model.Node{
		Id:          utils.UUID(),
		Name:        p.Name,
		ChainId:     p.ChainId,
		Description: p.Description,
		Root:        false,
	}

	err = service.ChainCRUD.SaveNode(&node)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithNoData())

}

// BindShortcut 绑定快捷方式到节点
func BindShortcut(c *gin.Context) {

	p := param.BindShortcutToNodeParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err = service.ChainCRUD.BindShortcut(p.NodeId, p.ShortcutId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithNoData())
}

// LinkNodes 连接两个节点
func LinkNodes(c *gin.Context) {

	p := param.ConnectNodesParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	if p.Type != model.SUCCESS && p.Type != model.FAILED {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err = service.ChainCRUD.LinkNode(chainService.ConnectTwoNodesParams{
		SourceId: p.SourceId,
		TargetId: p.TargetId,
		ChainId:  p.ChainId,
		Type:     p.Type,
	})

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
	}
	c.JSON(http.StatusOK, response.SuccessWithNoData())
}

// SetChainRoot 设置链的根节点
func SetChainRoot(c *gin.Context) {

	nodeId, ok := c.GetQuery("nodeId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err := service.ChainCRUD.SetRoot(nodeId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithNoData())
}

// ChainInfo 获取链信息
func ChainInfo(c *gin.Context) {

	id, ok := c.GetQuery("chainId")
	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	chainInfo, err := service.ChainCRUD.ChainInfo(id)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(chainInfo))

}

// ChainList 获取链列表
func ChainList(c *gin.Context) {

	list, err := service.ChainCRUD.ChainList()

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(list))
}

// ChainExec 执行链
func ChainExec(c *gin.Context) {

	chainId, ok := c.GetQuery("chainId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	dispatchId, err := service.ChainEXEC.Exec(chainId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(dispatchId))
}

// ChainExecLog 获取链执行日志
func ChainExecLog(c *gin.Context) {

	dispatchId, ok := c.GetQuery("dispatchId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	log, err := service.ChainEXEC.GetDispatchStatus(dispatchId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(log))

}
