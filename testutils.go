package shutdown

import (
	"errors"
	"log"
	"testing"
	"time"
)

var (
	errStoppedService = errors.New("service stopped")
)

type blankLog struct {}

func (*blankLog) Errorf(format string, args ...interface{}) {}
func (*blankLog) Infof(format string, args ...interface{}) {}
func (*blankLog) Debugf(format string, args ...interface{}) {}
func (*blankLog) Warnf(format string, args ...interface{}) {}

type testLog struct {
	t *testing.T
}

func (t *testLog) Errorf(format string, args ...interface{}) {
	t.t.Logf("[error] "+format, args...)
}
func (t *testLog) Infof(format string, args ...interface{}) {
	t.t.Logf("[info] "+format, args...)
}
func (t *testLog) Debugf(format string, args ...interface{}) {
	t.t.Logf("[debug] "+format, args...)
}
func (t *testLog) Warnf(format string, args ...interface{}) {
	t.t.Logf("[warning] "+format, args...)
}

type serviceWithStop struct {
	serving bool
	stopped bool
}

func (s* serviceWithStop) serve(timer time.Duration) {
	s.serving = true
	for {
		if s.stopped {
			s.serving = false
			return
		}
		time.Sleep(timer)
	}
}

func (s *serviceWithStop) Stop() error {
	s.stopped = true
	return nil
}

func (s *serviceWithStop) Ping() error {
	log.Print("ping happened")
	if !s.serving {
		return errStoppedService
	}
	return nil
}

func (s *serviceWithStop) Reconnect() error {
	return nil
}