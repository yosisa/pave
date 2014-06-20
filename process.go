package main

import (
	"errors"
	"sync/atomic"
	"time"
)

type ProcessState uint
type RestartStrategy uint

const (
	StateStopped ProcessState = iota
	StateRunning
	StateError
)

const (
	StrategyNoRestart RestartStrategy = iota
	StrategyRestartAlways
	StrategyRestartOnError
)

func (s *RestartStrategy) UnmarshalFlag(v string) error {
	switch v {
	case "none":
		*s = StrategyNoRestart
	case "always":
		*s = StrategyRestartAlways
	case "error":
		*s = StrategyRestartOnError
	default:
		return errors.New("Unsupported restart strategy")
	}
	return nil
}

type Process struct {
	*Command
	Error     error
	State     ProcessState
	StoppedAt time.Time
}

func (p *Process) SetError(err error) {
	p.State = StateError
	p.Error = err
	p.StoppedAt = time.Now()
}

type ProcessManager struct {
	Strategy    RestartStrategy
	RestartWait time.Duration
	processes   []*Process
	numActives  int32
	doneCh      chan bool
	stopCh      chan bool
}

func NewProcessManager(strategy RestartStrategy, restartWait time.Duration) *ProcessManager {
	return &ProcessManager{
		Strategy:    strategy,
		RestartWait: restartWait,
		doneCh:      make(chan bool, 1),
		stopCh:      make(chan bool, 1),
	}
}

func (m *ProcessManager) Add(cmd *Command) {
	p := &Process{Command: cmd}
	m.processes = append(m.processes, p)
}

func (m *ProcessManager) Start() {
	for _, p := range m.processes {
		m.startProcess(p)
	}
}

func (m *ProcessManager) Wait() {
	<-m.doneCh
}

func (m *ProcessManager) Run() {
	m.Start()
	m.Wait()
}

func (m *ProcessManager) Stop() {
	close(m.stopCh)
}

func (m *ProcessManager) startProcess(p *Process) {
	if err := p.Start(); err != nil {
		p.SetError(err)
		return
	}
	p.State = StateRunning
	atomic.AddInt32(&m.numActives, 1)
	go m.sentinel(p)
}

func (m *ProcessManager) sentinel(p *Process) {
	err := p.Cmd.Wait()
	if err != nil {
		p.SetError(err)
	}

	if m.Strategy == StrategyRestartAlways || m.Strategy == StrategyRestartOnError && err != nil {
		select {
		case <-m.stopCh:
		case <-time.After(m.RestartWait):
			m.startProcess(p)
		}
	}

	if atomic.AddInt32(&m.numActives, -1) == 0 {
		close(m.doneCh)
	}
}
