package process

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"sync/atomic"
	"syscall"
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
	Error      error
	State      ProcessState
	StoppedAt  time.Time
	errHandler func(error)
}

func New(cmd string, prepareFn func(*exec.Cmd), errHandler func(error)) *Process {
	command := NewCommand(cmd)
	command.PrepareFunc = prepareFn
	return &Process{
		Command:    command,
		errHandler: errHandler,
	}
}

func (p *Process) SetError(err error) {
	p.State = StateError
	p.Error = err
	p.StoppedAt = time.Now()
	if p.errHandler != nil {
		p.errHandler(err)
	}
}

type ProcessManager struct {
	Strategy    RestartStrategy
	RestartWait time.Duration
	processes   []*Process
	numActives  int32
	doneCh      chan bool
	stopCh      chan bool
	sigCh       chan os.Signal
}

func NewProcessManager(strategy RestartStrategy, restartWait time.Duration) *ProcessManager {
	m := &ProcessManager{
		Strategy:    strategy,
		RestartWait: restartWait,
		doneCh:      make(chan bool, 1),
		stopCh:      make(chan bool, 1),
		sigCh:       make(chan os.Signal, 1),
	}

	signal.Notify(m.sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range m.sigCh {
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				m.Stop()
			default:
				m.SignalAll(sig)
			}
		}
	}()

	return m
}

func (m *ProcessManager) Add(p *Process) {
	m.processes = append(m.processes, p)
}

func (m *ProcessManager) Start() {
	// Increment the counter to avoid 0 active processes and abnormal exiting during start-up.
	// It's occurred due to super short running command.
	atomic.AddInt32(&m.numActives, 1)
	for _, p := range m.processes {
		m.startProcess(p)
	}
	m.maybeDone()
}

func (m *ProcessManager) Wait() {
	<-m.doneCh
	signal.Stop(m.sigCh)
	close(m.sigCh)
}

func (m *ProcessManager) Run() {
	m.Start()
	m.Wait()
}

func (m *ProcessManager) Stop() {
	close(m.stopCh)
	m.SignalAll(syscall.SIGTERM)
	select {
	case <-m.doneCh:
	case <-time.After(10 * time.Second):
		m.SignalAll(syscall.SIGKILL)
	}
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
	defer m.maybeDone()

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
}

func (m *ProcessManager) maybeDone() {
	if atomic.AddInt32(&m.numActives, -1) == 0 {
		close(m.doneCh)
	}
}

func (m *ProcessManager) SignalAll(sig os.Signal) {
	for _, p := range m.processes {
		p.Cmd.Process.Signal(sig)
	}
}
