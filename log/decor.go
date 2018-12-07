package log

func Error(v interface{}) {
	logger.Log(LevelError, v)
}

func Errorf(format string, args ...interface{}) {
	logger.Logf(LevelError, format, args...)
}

func Warn(v interface{}) {
	logger.Log(LevelWarn, v)
}

func Warnf(format string, args ...interface{}) {
	logger.Logf(LevelWarn, format, args...)
}

func Info(v interface{}) {
	logger.Log(LevelInfo, v)
}

func Infof(format string, args ...interface{}) {
	logger.Logf(LevelInfo, format, args...)
}

func Debug(v interface{}) {
	logger.Log(LevelDebug, v)
}

func Debugf(format string, args ...interface{}) {
	logger.Logf(LevelDebug, format, args...)
}
