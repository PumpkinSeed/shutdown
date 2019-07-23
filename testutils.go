package shutdown

type blankLog struct {}

func (*blankLog) Errorf(format string, args ...interface{}) {}
func (*blankLog) Infof(format string, args ...interface{}) {}
func (*blankLog) Debugf(format string, args ...interface{}) {}
func (*blankLog) Warnf(format string, args ...interface{}) {}

type serviceWithStop struct {}

func (*serviceWithStop) Stop() error {
	return nil
}