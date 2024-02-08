package chain

import (
	"apps/center/action"
	"apps/center/model"
	"errors"
	"gorm.io/gorm"
	"time"
)

type LogService struct {
	Db *gorm.DB
}

// LoadElseNewDispatchLog 获取或新建调度日志
func (l *LogService) LoadElseNewDispatchLog(dispatcher *action.Dispatcher) (dispatchLog *model.DispatchLog, err error) {

	dispatchLog = &model.DispatchLog{}

	err = l.Db.Where("id = ? ", dispatcher.Id).First(dispatchLog).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {

		dispatchLog = &model.DispatchLog{
			Id:         dispatcher.Id,
			ChainId:    dispatcher.ChainId,
			Status:     model.DispatchStatusRunning,
			Done:       false,
			CreateTime: time.Now(),
		}

		err = l.Db.Create(dispatchLog).Error

		if err != nil {
			return nil, err
		}

		return dispatchLog, nil

	}

	if err != nil {
		return nil, err
	}

	return dispatchLog, nil

}

func (l *LogService) UpdateExecStatus(dispatcherId string, status model.DispatchStatus) error {
	done := status != model.DispatchStatusRunning
	return l.Db.Model(&model.DispatchLog{}).
		Where("id = ?", dispatcherId).
		Updates(map[string]interface{}{"status": status, "done": done}).
		Error
}

func (l *LogService) GetDispatchLog(dispatchId string) (dispatchLog *model.DispatchLog, err error) {
	dispatchLog = &model.DispatchLog{}
	err = l.Db.Where("id = ?", dispatchId).First(dispatchLog).Error
	return dispatchLog, err
}

func (l *LogService) GetNodeExecLogs(dispatchId string) (nodeExecLogs []*model.NodeExecLog, err error) {
	err = l.Db.Where("dispatch_id = ?", dispatchId).Order("create_time").Find(&nodeExecLogs).Error
	return nodeExecLogs, err
}

// GetNodeExecDetail
// SQL ：
//
//		SELECT
//		ne.create_time,
//		ne.ok,
//		ne.out,
//		nodes.name as 'node_name',
//		shortcuts.name as 'shortcut_name',
//		shortcuts.payload,
//		shortcuts.type as 'shortcut_type'
//	 shortcuts.probe_id as 'probeId'
//
// FROM
//
//	node_exec_logs ne,
//	nodes,
//	shortcut_node_bindings bd,
//	shortcuts
//
// WHERE
//
//	ne.node_id = nodes.id
//	AND nodes.id = bd.node_id
//	AND bd.shortcut_id = shortcuts.id
//	AND ne.dispatch_id = '814570a6-5f38-459d-894e-5fac4c91ca82'
//	ORDER BY ne.create_time DESC
func (l *LogService) GetNodeExecDetail(dispatchId string) (nodeExecLog []*model.NodeExecLog, err error) {

	//err = l.Db.Table("node_exec_logs ne").
	//	Select("ne.create_time,ne.ok,ne.out,nodes.name as 'node_name',shortcuts.name as 'shortcut_name',shortcuts.payload,shortcuts.type as 'shortcut_type',shortcuts.probe_id as 'probe_id'").
	//	Joins("inner join nodes on ne.node_id = nodes.id").
	//	Joins("left join shortcut_node_bindings bd on nodes.id = bd.node_id").
	//	Joins("left join shortcuts on bd.shortcut_id = shortcuts.id").
	//	Where("ne.dispatch_id = ?", dispatchId).
	//	Order("ne.create_time desc").
	//	Find(&nodeExecDetail).Error

	err = l.Db.Where("dispatch_id = ?", dispatchId).Order("create_time desc").Find(&nodeExecLog).Error
	return nodeExecLog, err

}

func (l *LogService) GetChainExecHistory(chainId string) (chainExecHistory []*model.DispatchLog, err error) {
	err = l.Db.Where("chain_id = ?", chainId).Order("create_time desc").Find(&chainExecHistory).Error
	return chainExecHistory, err
}
