package shutdown

import "time"

type blankLog struct {}

func (*blankLog) Errorf(format string, args ...interface{}) {}
func (*blankLog) Infof(format string, args ...interface{}) {}
func (*blankLog) Debugf(format string, args ...interface{}) {}
func (*blankLog) Warnf(format string, args ...interface{}) {}

type serviceWithStop struct {
	stopped bool
}

func (s* serviceWithStop) serve(timer time.Duration) {
	for {
		if s.stopped {
			return
		}
		time.Sleep(timer)
	}
}

func (s *serviceWithStop) Stop() error {
	s.stopped = true
	return nil
}