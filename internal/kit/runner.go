package kit

import (
	"sync"
)

type RunnerContext struct {
	error     chan error
	errorOnce sync.Once
}

func NewRunnerContext() *RunnerContext {
	return &RunnerContext{
		error: make(chan error, 1),
	}
}

func (rc *RunnerContext) Error(err error) {
	rc.errorOnce.Do(func() {
		rc.error <- err
	})
}

func (rc *RunnerContext) Errored() chan error {
	return rc.error
}

type Runner interface {
	Start(ctx *RunnerContext)
	Stop(ctx *RunnerContext)
}

type RunnerSlice struct {
	runners []Runner
}

func NewRunnerSlice() *RunnerSlice {
	return &RunnerSlice{
		runners: make([]Runner, 0),
	}
}

func (rs *RunnerSlice) WithRunner(r Runner) {
	rs.runners = append(rs.runners, r)
}

func (rs *RunnerSlice) Start(ctx *RunnerContext) {
	for _, r := range rs.runners {
		r.Start(ctx)
	}
}

func (rs *RunnerSlice) Stop(ctx *RunnerContext) {
	for i := len(rs.runners) - 1; i >= 0; i-- {
		rs.runners[i].Stop(ctx)
	}
}
