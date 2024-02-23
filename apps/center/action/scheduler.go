package action

import (
	"errors"
	"sync"
)

// Scheduler orchestrate and run
type Scheduler interface {
	RunOne(seeker *Seeker, runner *ShortcutRunner) (*NodeRunResult, error)
	RunAll(seeker *Seeker, runner *ShortcutRunner) ([]*NodeRunResult, error)
	Cancel() error
}

type NormalScheduler struct {
	lock       *sync.Mutex
	jobQueue   []*NodeRunJob
	cancelChan chan struct{}
	preResults []*NodeRunResult
	CancelMark bool
}

func NewNormalScheduler() *NormalScheduler {
	return &NormalScheduler{
		lock:       &sync.Mutex{},
		jobQueue:   make([]*NodeRunJob, 0),
		cancelChan: make(chan struct{}, 1),
		preResults: make([]*NodeRunResult, 0),
		CancelMark: false,
	}
}

func (n *NormalScheduler) TakeThenRunShortcut(runner *ShortcutRunner) (*NodeRunResult, error) {

	select {
	case <-n.cancelChan:
		n.CancelMark = true
		return nil, ErrCanceled
	default:
		if len(n.jobQueue) == 0 {
			return nil, ErrJobQueueEmpty
		}

		// take first job in job queue
		job := n.jobQueue[0]
		n.jobQueue = n.jobQueue[1:]

		runResult, err := (*runner).Run(*job.Shortcut)

		return &NodeRunResult{
			RunResult: &runResult,
			Node:      job.Node,
			Shortcut:  job.Shortcut,
		}, err
	}
}

// TODO TEST
func (n *NormalScheduler) RunOne(seeker *Seeker, runner *ShortcutRunner) (*NodeRunResult, error) {

	if n.CancelMark {
		return nil, ErrCanceled
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	runJobs, err := (*seeker).SeekNext(n.preResults)
	n.jobQueue = append(n.jobQueue, runJobs...)

	result, err := n.TakeThenRunShortcut(runner)

	if err != nil {
		return nil, err
	}

	n.preResults = nil
	n.preResults = append(n.preResults, result)

	return result, nil
}

// TODO TEST
func (n *NormalScheduler) RunAll(seeker *Seeker, runner *ShortcutRunner) ([]*NodeRunResult, error) {

	var res []*NodeRunResult

	for {
		result, err := n.RunOne(seeker, runner)

		if errors.Is(err, ErrJobQueueEmpty) {
			return res, nil
		}

		if err != nil {
			return res, err
		}
		res = append(res, result)
	}
}

func (n *NormalScheduler) Cancel() error {
	n.cancelChan <- struct{}{}
	return nil
}
