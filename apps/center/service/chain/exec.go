package chain

import (
	"apps/center/action"
	"apps/center/api/vo"
	"apps/center/model"
	"apps/common/utils"
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

type ExecService struct {
	ChainCRUD *CRUDService
	Log       *LogService
	// key: dispatchId value: *dispatcher
	dispatcherMap sync.Map
}

// Exec 执行链
func (e *ExecService) Exec(chainId string) (dispatchId string, error error) {

	dispatcher, err := e.NewDispatcher(chainId)

	if err != nil {
		return "", err
	}

	go e.DoDispatch(dispatcher)

	dispatchId = dispatcher.Id
	error = nil
	return
}

func (e *ExecService) GetDispatcherFromMap(dispatchId string) (dispatcher *action.Dispatcher, ok bool) {
	dp, ok := e.dispatcherMap.Load(dispatchId)
	if ok {
		dispatcher = dp.(*action.Dispatcher)
	}
	return
}

// NewDispatch 单步调度前的准备
func (e *ExecService) NewDispatch(chainId string) (*action.Dispatcher, error) {
	dispatcher, err := e.NewDispatcher(chainId)
	if err == nil {
		_, err = e.Log.LoadElseNewDispatchLog(dispatcher)
	}
	return dispatcher, err
}

func (e *ExecService) DoSingleStepDispatch(dispatchId string) error {

	dispatcher, ok := e.dispatcherMap.Load(dispatchId)

	if !ok {
		return errors.New("no such dispatcher or dispatcher has been done")
	}

	d := dispatcher.(*action.Dispatcher)

	curNode, ok := d.GetCurNode()
	shortcut, sok := d.CurShortcut()

	err := d.Next()

	if errors.Is(err, action.ErrNoMoreNode) {
		// 更新调度日志
		e.Log.UpdateExecStatus(dispatchId, model.DispatchStatusDone)
		e.dispatcherMap.Delete(dispatchId)
		return nil
	}

	if errors.Is(err, action.ErrAbort) {
		// 更新调度日志
		e.Log.UpdateExecStatus(dispatchId, model.DispatchStatusAborted)
		e.dispatcherMap.Delete(dispatchId)
		return nil
	}

	if ok {

		nodeExecLog := model.NodeExecLog{
			Id:         utils.UUID(),
			DispatchId: d.Id,
			Ok:         d.Ok,
			Out:        d.Out,
			NodeName:   &curNode.Name,
			ChainId:    d.ChainId,
			NodeId:     curNode.Id,
			CreateTime: time.Now(),
		}

		if sok {
			nodeExecLog.ShortcutId = &shortcut.Id
			nodeExecLog.ShortcutName = &shortcut.Name
			nodeExecLog.ShortcutType = &shortcut.Type
			nodeExecLog.ProbeId = &shortcut.ProbeId
			nodeExecLog.Payload = &shortcut.Payload
		}

		e.Log.Db.Create(&nodeExecLog)
	}

	if d.Status == action.Done {
		// 更新调度日志
		e.Log.UpdateExecStatus(dispatchId, model.DispatchStatusDone)
		e.dispatcherMap.Delete(dispatchId)
		return nil
	}

	return nil

}

func (e *ExecService) NewDispatcher(chainId string) (dispatcher *action.Dispatcher, err error) {
	// 1. 从数据库中获取所有的节点、边、快捷方式、绑定关系
	nodes, err := e.ChainCRUD.GetNodesByChainId(chainId)
	if err != nil {
		return nil, err
	}

	edges, err := e.ChainCRUD.GetEdgesByChainId(chainId)
	if err != nil {
		return nil, err
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
		ChainId:   chainId,
	}

	// 3. 创建调度器
	dispatcher, err = action.NewDispatcher(params)

	if err == nil {
		e.dispatcherMap.Store(dispatcher.Id, dispatcher)
	}

	return
}

// DoDispatch 执行调度
// 接收一个调度器，执行调度器的Next方法，直到调度器没有下一个节点
func (e *ExecService) DoDispatch(dispatcher *action.Dispatcher) error {

	dispatchLog, err := e.Log.LoadElseNewDispatchLog(dispatcher)

	if err != nil {
		return err
	}

	for {

		curNode, ok := dispatcher.GetCurNode()
		shortcut, sok := dispatcher.CurShortcut()

		err := dispatcher.Next()
		if errors.Is(err, action.ErrNoMoreNode) {
			// 更新调度日志
			e.Log.UpdateExecStatus(dispatchLog.Id, model.DispatchStatusDone)
			e.dispatcherMap.Delete(dispatcher.Id)
			break
		}

		if errors.Is(err, action.ErrAbort) {
			// 更新调度日志
			e.Log.UpdateExecStatus(dispatchLog.Id, model.DispatchStatusAborted)
			e.dispatcherMap.Delete(dispatcher.Id)
			break
		}
		if ok {

			nodeExecLog := model.NodeExecLog{
				Id:         utils.UUID(),
				DispatchId: dispatcher.Id,
				Ok:         dispatcher.Ok,
				Out:        dispatcher.Out,
				NodeId:     curNode.Id,
				ChainId:    dispatcher.ChainId,
				NodeName:   &curNode.Name,
				CreateTime: time.Now(),
			}

			if sok {
				nodeExecLog.ShortcutId = &shortcut.Id
				nodeExecLog.ShortcutName = &shortcut.Name
				nodeExecLog.ShortcutType = &shortcut.Type
				nodeExecLog.ProbeId = &shortcut.ProbeId
				nodeExecLog.Payload = &shortcut.Payload
			}

			e.Log.Db.Create(&nodeExecLog)
		}
	}

	return nil

}

// GetDispatchStatus 获取调度情况
func (e *ExecService) GetDispatchStatus(dispatchId string) (d *vo.DispatchStatus, err error) {

	d = &vo.DispatchStatus{}

	dispatchLog, err := e.Log.GetDispatchLog(dispatchId)

	if err != nil {
		return nil, err
	}
	d.DispatchLog = dispatchLog

	nodeExecDetails, err := e.Log.GetNodeExecDetail(dispatchId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	d.PreNodes = nodeExecDetails

	dp, ok := e.dispatcherMap.Load(dispatchId)

	if !ok {
		d.CurNode = nil
		d.SuccessThen = nil
		d.FailThen = nil
		return d, nil
	}
	dispatcher := dp.(*action.Dispatcher)

	shortcut, err := e.ChainCRUD.GetShortcutByNodeId(dispatcher.CurNode.Id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	d.CurNode = &vo.NodeDetail{
		Node:     dispatcher.CurNode,
		Shortcut: shortcut,
	}

	successThenNode, ok := dispatcher.SuccessThenNode()

	if ok {

		shortcut, err := e.ChainCRUD.GetShortcutByNodeId(successThenNode.Id)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		d.SuccessThen = &vo.NodeDetail{
			Node:     successThenNode,
			Shortcut: shortcut,
		}

	}

	failThenNode, ok := dispatcher.FailedThenNode()

	if ok {

		shortcut, err := e.ChainCRUD.GetShortcutByNodeId(failThenNode.Id)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		d.FailThen = &vo.NodeDetail{
			Node:     failThenNode,
			Shortcut: shortcut,
		}

	}
	return
}
