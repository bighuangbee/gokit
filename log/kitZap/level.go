package kitZap

import "go.uber.org/zap/zapcore"

func ToZapLevel(lv string) (r zapcore.Level) {
	switch lv {
	case "debug":
		r = zapcore.DebugLevel
	case "warn":
		r = zapcore.WarnLevel
	case "error":
		r = zapcore.ErrorLevel
	default:
		// default info
		r = zapcore.InfoLevel
	}
	return
}
