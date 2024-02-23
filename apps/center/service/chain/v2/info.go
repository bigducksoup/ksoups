package v2

import (
	"apps/center/api/vo"
	"apps/center/model"
	"gorm.io/gorm"
)

type InfoService struct {
	Db *gorm.DB
}

func (i *InfoService) ExecInfoList(chainId string) ([]model.DispatchLog, error) {
	var res []model.DispatchLog
	err := i.Db.Debug().Where("chain_id = ?", chainId).
		Select([]string{"id",
			"chain_id",
			"create_time",
			"status",
			"done"}).
		Order("create_time desc").
		Find(&res).Error
	return res, err
}

func (i *InfoService) ExecInfoDetail(dispatchId string) (*vo.ExecDetail, error) {
	var dispatch model.DispatchLog

	err := i.Db.Where("id = ?", dispatchId).First(&dispatch).Error

	if err != nil {
		return nil, err
	}

	var nodeExecs []model.NodeExecLog

	err = i.Db.Where("dispatch_id = ?", dispatchId).Find(&nodeExecs).Error

	if err != nil {
		return nil, err
	}

	return &vo.ExecDetail{
		Info:    &dispatch,
		Results: nodeExecs,
	}, nil
}
