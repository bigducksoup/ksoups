package handler

import (
	"apps/center/api/param"
	"apps/center/api/response"
	"apps/center/model"
	"apps/center/service"
	chainService "apps/center/service/chain"
	"apps/common/utils"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, response.Success(ch))
}

// ChainDelete 删除链
func ChainDelete(c *gin.Context) {

	id, ok := c.GetQuery("chainId")
	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err := service.ChainCRUD.DeleteChain(id)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
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

func NodeDelete(c *gin.Context) {

	id, ok := c.GetQuery("nodeId")
	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err := service.ChainCRUD.DeleteNode(id)

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

func UnBindShortcut(c *gin.Context) {

	p := param.BindShortcutToNodeParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err = service.ChainCRUD.UnbindShortcut(p.NodeId, p.ShortcutId)

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
		return
	}
	c.JSON(http.StatusOK, response.SuccessWithNoData())
}

// UnlinkNodes 解除两个节点的连接

func UnLinkNodes(c *gin.Context) {

	p := param.ConnectNodesParams{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err = service.ChainCRUD.UnlinkNode(chainService.ConnectTwoNodesParams{
		SourceId: p.SourceId,
		TargetId: p.TargetId,
		ChainId:  p.ChainId,
		Type:     p.Type,
	})

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
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

func ChainExecLogHistory(c *gin.Context) {

	chainId, ok := c.GetQuery("chainId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
	}

	history, err := service.ChainLOG.GetChainExecHistory(chainId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(history))

}

func ChainLoadFromAllData(c *gin.Context) {

	p := param.ChainAllDataParams{}
	err := c.ShouldBindJSON(&p)

	if err != nil {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	err = service.ChainCRUD.LoadFromAllData(&p)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(true))
}

func NewExecMachine(c *gin.Context) {

	chainId, ok := c.GetQuery("chainId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	machine, err := service.ChainEXECV2.CreateMachine(chainId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}
	var dispatch model.DispatchLog

	err = service.ChainCRUD.Db.Where("id = ?", machine.Id).
		Select("id", "chain_id", "create_time", "status", "done").
		Find(&dispatch).Error

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(dispatch))
}

func MachineExecAll(c *gin.Context) {

	machineId, ok := c.GetQuery("machineId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	_, err := service.ChainEXECV2.MachineDoNextAll(machineId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	logs, err := service.ChainLOG.GetNodeExecLogs(machineId)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(logs))

}

func MachineExecOne(c *gin.Context) {
	machineId, ok := c.GetQuery("machineId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	_, err := service.ChainEXECV2.MachineDoNextOne(machineId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	logs, err := service.ChainLOG.GetNodeExecLogs(machineId)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(logs))
}

func ExecList(c *gin.Context) {

	chainId, ok := c.GetQuery("chainId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	infoList, err := service.ChainINFO.ExecInfoList(chainId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(infoList))

}

func ExecDetail(c *gin.Context) {

	dispatchId, ok := c.GetQuery("dispatchId")

	if !ok {
		c.JSON(http.StatusOK, response.ParamsError())
		return
	}

	detail, err := service.ChainINFO.ExecInfoDetail(dispatchId)

	if err != nil {
		c.JSON(http.StatusOK, response.Fail(err))
		return
	}

	c.JSON(http.StatusOK, response.Success(detail))

}
