package logger

// 日志服务伪实现（主要用于单元测试）
type DummyLogger struct {
}

func NewDummyLogger(logLevel string) *DummyLogger {
	return &DummyLogger{}
}

func (l *DummyLogger) Debug(args ...interface{}) {
}

func (l *DummyLogger) Debugf(template string, args ...interface{}) {
}

func (l *DummyLogger) Info(args ...interface{}) {
}

func (l *DummyLogger) Infof(template string, args ...interface{}) {
}

func (l *DummyLogger) Warn(args ...interface{}) {
}

func (l *DummyLogger) Warnf(template string, args ...interface{}) {
}

func (l *DummyLogger) Error(args ...interface{}) {
}

func (l *DummyLogger) Errorf(template string, args ...interface{}) {
}

func (l *DummyLogger) Panic(args ...interface{}) {
}

func (l *DummyLogger) Panicf(template string, args ...interface{}) {
}

func (l *DummyLogger) Fatal(args ...interface{}) {
}

func (l *DummyLogger) Fatalf(template string, args ...interface{}) {
}
