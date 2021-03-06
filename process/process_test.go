package process

import (
	"bytes"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRestartStrategy(t *testing.T) {
	var s RestartStrategy
	s.UnmarshalFlag("none")
	assert.Equal(t, StrategyNoRestart, s)

	s.UnmarshalFlag("always")
	assert.Equal(t, StrategyRestartAlways, s)

	s.UnmarshalFlag("error")
	assert.Equal(t, StrategyRestartOnError, s)

	err := s.UnmarshalFlag("other")
	assert.Equal(t, "Unsupported restart strategy", err.Error())
}

func TestProcessManager(t *testing.T) {
	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	pm := NewProcessManager(StrategyNoRestart, 0)

	pm.Add(New(`echo -n foo`, func(cmd *exec.Cmd) { cmd.Stdout = b1 }, nil))
	pm.Add(New(`echo -n bar`, func(cmd *exec.Cmd) { cmd.Stdout = b2 }, nil))
	pm.Run()
	assert.Equal(t, "foo", b1.String())
	assert.Equal(t, "bar", b2.String())
}

func TestProcessManagerWithRestartAlways(t *testing.T) {
	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	pm := NewProcessManager(StrategyRestartAlways, 100*time.Millisecond)

	pm.Add(New(`echo -n foo`, func(cmd *exec.Cmd) { cmd.Stdout = b1 }, nil))
	pm.Add(New(`echo -n bar`, func(cmd *exec.Cmd) { cmd.Stdout = b2 }, nil))
	pm.Start()
	go func() {
		<-time.After(150 * time.Millisecond)
		pm.Stop()
	}()
	pm.Wait()
	assert.Equal(t, "foofoo", b1.String())
	assert.Equal(t, "barbar", b2.String())
}

func TestProcessManagerStop(t *testing.T) {
	b1 := new(bytes.Buffer)
	pm := NewProcessManager(StrategyNoRestart, 0)

	pm.Add(New(`sleep 1`, func(cmd *exec.Cmd) { cmd.Stdout = b1 }, nil))
	start := time.Now()
	pm.Start()
	pm.Stop()
	pm.Wait()
	assert.True(t, time.Since(start).Seconds() < 0.1)
}
