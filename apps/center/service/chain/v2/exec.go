package v2

import (
	"apps/center/action"
	"apps/center/model"
	"apps/center/service/chain"
	"apps/common/utils"
	"errors"
	"sync"
	"time"
)

type ExecService struct {
	machineMap  sync.Map
	CRUDService *chain.CRUDService
}

func (e *ExecService) LoadMachine(machineId string) (*action.Machine, error) {

	m, ok := e.machineMap.Load(machineId)
	if !ok {
		return nil, errors.New("could not find machine with id = " + machineId)
	}
	return m.(*action.Machine), nil
}

func (e *ExecService) StoreMachine(machine *action.Machine) {
	e.machineMap.Store(machine.Id, machine)
}

func (e *ExecService) DeleteMachine(machineId string) {
	e.machineMap.Delete(machineId)
}

func (e *ExecService) CreateMachine(chainId string) (*action.Machine, error) {

	var c model.Chain

	err := e.CRUDService.Db.Where("id = ?", chainId).Select("origin_data").Take(&c).Error

	if err != nil {
		return nil, err
	}

	nodes, err := e.CRUDService.GetNodesByChainId(chainId)
	if err != nil {
		return nil, err
	}

	edges, err := e.CRUDService.GetEdgesByChainId(chainId)
	if err != nil {
		return nil, err
	}

	var shortcuts []model.Shortcut

	for _, node := range nodes {
		shortcut, err := e.CRUDService.GetShortcutByNodeId(node.Id)
		if err != nil {
			continue
		}
		shortcuts = append(shortcuts, *shortcut)
	}

	bindings, err := e.CRUDService.GetBindingsByChainID(chainId)

	if err != nil {
		return nil, err
	}

	var runner action.ShortcutRunner = action.NewNormalShortcutRunner()
	var scheduler action.Scheduler = action.NewNormalScheduler()
	var seeker action.Seeker = action.NewNormalSeeker(nodes, edges, shortcuts, bindings)
	machine := action.NewMachine(&scheduler, &seeker, &runner)

	//create a record
	dispatch := model.DispatchLog{
		Id:         machine.Id,
		ChainId:    chainId,
		CreateTime: time.Now(),
		Status:     model.DispatchStatusReady,
		Done:       false,
		ChainData:  &c.OriginData,
	}

	err = e.CRUDService.Db.Create(&dispatch).Error

	if err != nil {
		return nil, err
	}

	e.StoreMachine(machine)

	return machine, nil
}

func (e *ExecService) MachineDoNextOne(machineId string) (*action.NodeRunResult, error) {
	machine, err := e.LoadMachine(machineId)

	if err != nil {
		return nil, err
	}

	result, err := machine.NextOne()

	if err != nil {
		if errors.Is(err, action.ErrJobQueueEmpty) {
			machine.Status = action.MachineJobDone
			e.DeleteMachine(machineId)
			e.CRUDService.Db.Model(&model.DispatchLog{}).Where("id = ?", machineId).Updates(map[string]any{
				"status": model.DispatchStatusDone,
				"done":   true,
			})
		}

		return nil, err
	}

	e.CRUDService.Db.Create(&model.NodeExecLog{
		Id:           utils.UUID(),
		NodeId:       result.Node.Id,
		ChainId:      result.Node.ChainId,
		DispatchId:   machineId,
		NodeName:     &result.Node.Name,
		ShortcutId:   &result.Shortcut.Id,
		ShortcutName: &result.Shortcut.Name,
		ProbeId:      &result.Shortcut.ProbeId,
		ShortcutType: &result.Shortcut.Type,
		Payload:      &result.Shortcut.Payload,
		Ok:           result.RunResult.Ok,
		StdOut:       result.RunResult.StdOut,
		StdErr:       result.RunResult.StdErr,
		CreateTime:   time.Now(),
	})

	return result, err
}

func (e *ExecService) MachineDoNextAll(machineId string) ([]*action.NodeRunResult, error) {
	machine, err := e.LoadMachine(machineId)

	if err != nil {
		return nil, err
	}

	results, err := machine.NextAll()

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		e.CRUDService.Db.Create(&model.NodeExecLog{
			Id:           utils.UUID(),
			NodeId:       result.Node.Id,
			ChainId:      result.Node.ChainId,
			DispatchId:   machineId,
			NodeName:     &result.Node.Name,
			ShortcutId:   &result.Shortcut.Id,
			ShortcutName: &result.Shortcut.Name,
			ProbeId:      &result.Shortcut.ProbeId,
			ShortcutType: &result.Shortcut.Type,
			Payload:      &result.Shortcut.Payload,
			Ok:           result.RunResult.Ok,
			StdOut:       result.RunResult.StdOut,
			StdErr:       result.RunResult.StdErr,
			CreateTime:   time.Now(),
		})
	}

	machine.Status = action.MachineJobDone
	e.DeleteMachine(machineId)
	e.CRUDService.Db.Model(&model.DispatchLog{}).Where("id = ?", machineId).Updates(map[string]any{
		"status": model.DispatchStatusDone,
		"done":   true,
	})

	return results, err
}

func (e *ExecService) MachineCancel(machineId string) error {

	machine, err := e.LoadMachine(machineId)

	if err != nil {
		return err
	}

	err = machine.Cancel()

	if err != nil {
		return err
	}

	e.DeleteMachine(machineId)
	return nil
}
