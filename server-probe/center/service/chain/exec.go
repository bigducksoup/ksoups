package chain

import (
	"config-manager/center/action"
	"config-manager/center/model"
	"config-manager/common/utils"
	"errors"
	"time"
)

type ExecService struct {
	ChainCRUD *CRUDService
}

// Exec 执行链
// 1. 从数据库中获取所有的节点、边、快捷方式、绑定关系
// 2. 创建调度器参数
// 3. 创建调度器
// 4. 创建调度日志
// 5. 执行调度器
// 6. 创建节点执行日志
func (e *ExecService) Exec(chainId string) (dispatchId string, error error) {

	// 1. 从数据库中获取所有的节点、边、快捷方式、绑定关系
	nodes, err := e.ChainCRUD.GetNodesByChainId(chainId)
	if err != nil {
		return "", err
	}

	edges, err := e.ChainCRUD.GetEdgesByChainId(chainId)
	if err != nil {
		return "", err
	}

	var nodeIds []string
	for _, node := range nodes {
		nodeIds = append(nodeIds, node.Id)
	}

	var bindings []model.ShortcutNodeBinding
	e.ChainCRUD.Db.Where("node_id in ?", nodeIds).Find(&bindings)

	var shortcutIds []string
	for _, binding := range bindings {
		shortcutIds = append(shortcutIds, binding.ShortcutId)
	}
	var shortcuts []model.Shortcut
	e.ChainCRUD.Db.Where("id in ?", shortcutIds).Find(&shortcuts)

	// 2. 创建调度器参数
	params := action.ChainExecParams{
		Shortcuts: shortcuts,
		Nodes:     nodes,
		Edges:     edges,
		Bindings:  bindings,
	}

	// 3. 创建调度器
	dispatcher, err := action.NewDispatcher(params)

	if err != nil {
		return "", err
	}

	go func() {
		// 4. 创建调度日志
		dispatchLog := model.DispatchLog{
			Id:         dispatcher.Id,
			ChainId:    chainId,
			Status:     model.DispatchStatusRunning,
			Done:       false,
			CreateTime: time.Now(),
		}
		e.ChainCRUD.Db.Create(&dispatchLog)

		// 5. 执行调度器
		for {

			curNodeId, ok := dispatcher.GetCurNodeId()

			err := dispatcher.Next()
			if errors.Is(err, action.ErrNoMoreNode) {

				// 更新调度日志
				e.ChainCRUD.Db.Model(&dispatchLog).Where("id = ? ", dispatchLog.Id).Updates(model.DispatchLog{
					Status: model.DispatchStatusDone,
					Done:   true,
				})

				break
			}

			// 6. 创建节点执行日志
			nodeExecLog := model.NodeExecLog{
				Id:         utils.UUID(),
				DispatchId: dispatcher.Id,
				Ok:         dispatcher.Ok,
				Out:        dispatcher.Out,
			}

			if ok {
				nodeExecLog.NodeId = *curNodeId
			}

			e.ChainCRUD.Db.Create(&nodeExecLog)
		}
	}()

	dispatchId = dispatcher.Id
	error = nil
	return
}

type DispatchStatus struct {
	DispatchLog model.DispatchLog              `json:"dispatchLog"`
	NodeExecLog map[string][]model.NodeExecLog `json:"nodeExecLog"`
}

// GetDispatchStatus 获取调度情况
func (e *ExecService) GetDispatchStatus(dispatchId string) (d DispatchStatus, err error) {

	err = e.ChainCRUD.Db.Where("id = ?", dispatchId).First(&d.DispatchLog).Error

	if err != nil {
		return
	}

	var nodeExecLogs []model.NodeExecLog
	err = e.ChainCRUD.Db.Where("dispatch_id = ?", dispatchId).Find(&nodeExecLogs).Error

	if err != nil {
		return
	}

	nodeExecLogMap, err := utils.SliceGroupBy[string, model.NodeExecLog](nodeExecLogs, "NodeId")
	d.NodeExecLog = nodeExecLogMap
	return
}
