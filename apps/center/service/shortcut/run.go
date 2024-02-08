package shortcut

import (
	"apps/center/action"
	"apps/center/model"
	"apps/common/utils"
	"errors"
	"gorm.io/gorm"
	"time"
)

type RUNService struct {
	Runner action.Runner
	Db     *gorm.DB
}

func (r *RUNService) Run(id string) (out string, error error) {

	execTime := time.Now()
	sc := model.Shortcut{}
	tx := r.Db.First(&sc, "id = ?", id)

	if tx.Error != nil {
		return "", error
	}

	out, ok := r.Runner.Run(sc)

	go func() {
		if r.Db != nil {
			r.Db.Create(&model.ShortcutExecLog{
				Id:          utils.UUID(),
				ShortcutId:  id,
				Out:         out,
				OK:          ok,
				CreateTime:  time.Now(),
				ExecuteTime: execTime,
			})
		}
	}()

	if !ok {
		error = errors.New(out)
		return
	}
	return
}

func (r *RUNService) RunHistory(id string) (runHistory []model.ShortcutExecLog, error error) {

	if r.Db == nil {
		return nil, errors.New("db is nil")
	}

	tx := r.Db.Order("execute_time desc").Find(&runHistory, "shortcut_id = ?", id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return
}
