package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	// 日志接口
	*zap.SugaredLogger
}

func (l *ZapLogger) setLogger(level zapcore.Level) {
	zapLogger, err := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			//CallerKey:    "caller",
			//EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()

	if err == nil {
		l.SugaredLogger = zapLogger.Sugar()
	}
}

func (l *ZapLogger) Close() {
	l.SugaredLogger.Sync()
}

func NewZapLogger(logLevel string) *ZapLogger {
	l := &ZapLogger{}

	var level zapcore.Level
	if level.UnmarshalText([]byte(logLevel)) != nil {
		return nil
	}

	l.setLogger(level)
	if l.SugaredLogger == nil {
		return nil
	}

	return l
}
