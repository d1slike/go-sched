package log

import (
	"log"
	"os"
)

var (
	logLevel = LevelWarn
	logger   = NewStdErrLogger("go-sched")
)

type Level int

const (
	LevelError = Level(iota)
	LevelWarn
	LevelInfo
	LevelDebug
)

type Logger interface {
	Log(level Level, entry interface{})
	Logf(level Level, format string, args ...interface{})
}

type stdErrLogger struct {
	logger *log.Logger
}

func (l *stdErrLogger) Log(level Level, entry interface{}) {
	if logLevel <= level {
		l.logger.Println(entry)
	}
}

func (l *stdErrLogger) Logf(level Level, format string, args ...interface{}) {
	if logLevel <= level {
		l.logger.Printf(format, args...)
	}
}

func NewStdErrLogger(name string) Logger {
	return &stdErrLogger{
		logger: log.New(os.Stderr, name+":", log.LstdFlags),
	}
}

func SetLogger(log Logger) {
	logger = log
}

func SetLogLevel(level Level) {
	logLevel = level
}
