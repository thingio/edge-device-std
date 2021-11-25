package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "errors"
	FatalLevel LogLevel = "fatal"
	PanicLevel LogLevel = "panic"
)

func NewLogger() *Logger {
	logger := &Logger{
		logger: logrus.NewEntry(logrus.New()),
	}
	_ = logger.SetLevel(InfoLevel)
	return logger
}

type Logger struct {
	logger *logrus.Entry
}

func (l Logger) SetLevel(level LogLevel) error {
	lvl, err := logrus.ParseLevel(string(level))
	if err != nil {
		return err
	}
	l.logger.Logger.SetLevel(lvl)
	l.logger.Logger.SetOutput(os.Stdout)
	l.logger.Logger.SetFormatter(&logFormatter{logrus.TextFormatter{FullTimestamp: true, ForceColors: true}})
	return nil
}

// WithFields adds a map of fields to the Entry.
func (l Logger) WithFields(vs ...string) *logrus.Entry {
	fs := logrus.Fields{}
	for index := 0; index < len(vs)-1; index = index + 2 {
		fs[vs[index]] = vs[index+1]
	}
	return l.logger.WithFields(fs)
}

// WithError adds an error as single field (using the key defined in ErrorKey) to the Entry.
func (l Logger) WithError(err error) *logrus.Entry {
	if err == nil {
		return l.logger
	}

	return l.logger.WithField(logrus.ErrorKey, err.Error())
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}
func (l Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l Logger) Warnf(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}
func (l Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
func (l Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
func (l Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}
func (l Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}
func (l Logger) Warn(args ...interface{}) {
	l.logger.Info(args...)
}
func (l Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}
func (l Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
func (l Logger) Panic(args ...interface{}) {
	l.logger.Panic(args)
}

type logFormatter struct {
	logrus.TextFormatter
}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data, err := f.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}
	return data, nil
}
