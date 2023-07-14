package mrapp

type LoggerStub struct { }

func NewLoggerStub() *LoggerStub {
    return &LoggerStub{}
}

// Make sure the LoggerStub conforms with the Logger interface
var _ Logger = (*LoggerStub)(nil)

func (l *LoggerStub) WithContext(name string) Logger {
    return l
}

func (l *LoggerStub) Fatal(message any, args ...any) {
}

func (l *LoggerStub) Error(message any, args ...any) {
}

func (l *LoggerStub) Warn(message string, args ...any) {
}

func (l *LoggerStub) Info(message string, args ...any) {
}

func (l *LoggerStub) Debug(message any, args ...any) {
}

func (l *LoggerStub) Event(message string, args ...any) {
}
