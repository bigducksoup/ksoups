package action

import (
	"apps/center/model"
	"apps/common/utils"
	"errors"
)

type RunResult struct {
	Ok     bool   `json:"ok"`
	StdErr string `json:"stdErr"`
	StdOut string `json:"stdOut"`
}

type NodeRunResult struct {
	RunResult *RunResult      `json:"runResult"`
	Node      *model.Node     `json:"node,omitempty"`
	Shortcut  *model.Shortcut `json:"shortcut"`
}

type NodeRunJob struct {
	Node     *model.Node     `json:"node,omitempty"`
	Shortcut *model.Shortcut `json:"shortcut,omitempty"`
}

var (
	ErrScheduleConflict = errors.New("scheduler is doing it's job right now")
	ErrJobQueueEmpty    = errors.New("job queue is empty, unable to run")
	ErrCanceled         = errors.New("current schedule has been canceled")
)

type MachineStatus uint8

const (
	MachineReady      MachineStatus = 0
	MachineJobDone    MachineStatus = 1
	MachineJobCancel  MachineStatus = 2
	MachineUnknownErr MachineStatus = 3
)

type Machine struct {
	Id             string
	shortcutRunner *ShortcutRunner
	seeker         *Seeker
	scheduler      *Scheduler
	Status         MachineStatus
}

func NewMachine(scheduler *Scheduler, seeker *Seeker, runner *ShortcutRunner) *Machine {
	return &Machine{
		Id:             utils.UUID(),
		seeker:         seeker,
		shortcutRunner: runner,
		scheduler:      scheduler,
		Status:         MachineReady,
	}
}

func (m *Machine) NextOne() (*NodeRunResult, error) {

	if m.Status != MachineReady {
		return nil, ErrCanceled
	}

	result, err := (*m.scheduler).RunOne(m.seeker, m.shortcutRunner)

	if errors.Is(err, ErrCanceled) {
		m.Status = MachineJobCancel
		return nil, err
	}

	if errors.Is(err, ErrJobQueueEmpty) {
		m.Status = MachineJobDone
		return nil, err
	}

	if err != nil {
		m.Status = MachineUnknownErr
		return nil, err
	}

	return result, nil
}

func (m *Machine) NextAll() ([]*NodeRunResult, error) {

	if m.Status != MachineReady {
		return nil, ErrCanceled
	}

	results, err := (*m.scheduler).RunAll(m.seeker, m.shortcutRunner)

	if errors.Is(err, ErrCanceled) {
		m.Status = MachineJobCancel
		return nil, err
	}

	if err != nil {
		m.Status = MachineUnknownErr
		return nil, err
	}

	m.Status = MachineJobDone
	return results, nil
}

func (m *Machine) Cancel() error {
	err := (*m.scheduler).Cancel()

	if err != nil {
		m.Status = MachineUnknownErr
		return err
	}

	m.Status = MachineJobCancel
	return err
}
