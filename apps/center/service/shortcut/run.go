package shortcut

import (
	"apps/center/action"
	"apps/center/model"
	"apps/common/message/data"
	"apps/common/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RUNService struct {
	Runner action.Runner
	Db     *gorm.DB
}

func (r *RUNService) Run(id string) (stdout string, stderr string, state model.ScriptRunState, err error) {

	execTime := time.Now()
	sc := model.Shortcut{}
	tx := r.Db.First(&sc, "id = ?", id)

	if tx.Error != nil {
		return "", "", 0, tx.Error
	}

	stdout, stderr, err = r.Runner.Run(sc)

	go func() {
		if r.Db != nil {
			state := model.SCRIPT_RUN_SUCCESS
			if err != nil {
				state = model.SCRIPT_RUN_FAIL
			} else if len(stderr) == 0 {
				state = model.SCRIPT_RUN_ERR
			}
			r.Db.Create(&model.ShortcutExecLog{
				Id:          utils.UUID(),
				ShortcutId:  sc.Id,
				StdOut:      stdout,
				StdErr:      stderr,
				State:       state,
				ExecuteTime: execTime,
				RunByChain:  false,
				ChainId:     nil,
				NodeId:      nil,
			})
		}
	}()

	return
}

func (r *RUNService) AsyncRun(id string) (streams chan data.AsyncShortCutRespStream, err error) {
	// execTime := time.Now()
	sc := model.Shortcut{}
	tx := r.Db.First(&sc, "id = ?", id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	streams, err = r.Runner.RealTimeRun(sc)

	if err != nil {
		return nil, err
	}

	if r.Db != nil {
		// TODO fix this
		// r.Db.Create(&model.ShortcutExecLog{
		// 	Id:          runId,
		// 	ShortcutId:  id,
		// 	Out:         "",
		// 	OK:          false,
		// 	CreateTime:  time.Now(),
		// 	ExecuteTime: execTime,
		// })
	}

	return streams, nil
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
