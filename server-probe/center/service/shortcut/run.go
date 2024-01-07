package shortcut

import (
	"config-manager/center/action"
	"config-manager/center/model"
	"config-manager/common/utils"
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
