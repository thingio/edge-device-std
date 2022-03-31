package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/thingio/edge-device-std/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func NewLogger(options *config.LogOptions) (*Logger, error) {
	root := logrus.NewEntry(logrus.New())

	logLevel, err := logrus.ParseLevel(options.Level)
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	root.Logger.Level = logLevel

	var logWriter io.Writer
	if options.Console {
		logWriter = os.Stdout
	} else {
		logWriter = ioutil.Discard
		if options.Path != "" {
			if err := os.Mkdir(filepath.Dir(options.Path), 0664); err != nil {
				return nil, err
			}
			if hook, err := newFileHook(fileConfig{
				Filename:   options.Path,
				MaxSize:    options.Size.Max,
				MaxAge:     options.Age.Max,
				MaxBackups: options.Backup.Max,
				Compress:   true,
				Level:      logLevel,
				Formatter:  newFormatter(options.Format, false),
			}); err != nil {
				return nil, err
			} else {
				root.Logger.Hooks.Add(hook)
			}
		}
	}


	root.Logger.SetOutput(logWriter)
	root.Logger.SetFormatter(newFormatter("text", true))

	return &Logger{logger: root}, nil
}

type Logger struct {
	logger *logrus.Entry
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

type fileConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
	Level      logrus.Level
	Formatter  logrus.Formatter
}

type fileHook struct {
	config fileConfig
	writer io.Writer
}

func newFileHook(config fileConfig) (logrus.Hook, error) {
	hook := fileHook{
		config: config,
	}

	var zeroLevel logrus.Level
	if hook.config.Level == zeroLevel {
		hook.config.Level = logrus.InfoLevel
	}
	var zeroFormatter logrus.Formatter
	if hook.config.Formatter == zeroFormatter {
		hook.config.Formatter = new(logrus.TextFormatter)
	}

	hook.writer = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		LocalTime:  config.LocalTime,
		Compress:   config.Compress,
	}

	return &hook, nil
}

// Levels Levels
func (hook *fileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

// Fire Fire
func (hook *fileHook) Fire(entry *logrus.Entry) (err error) {
	if hook.config.Level < entry.Level {
		return nil
	}
	b, err := hook.config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.writer.Write(b)
	return nil
}

func newFormatter(format string, color bool) logrus.Formatter {
	var formatter logrus.Formatter
	if strings.ToLower(format) == "json" {
		formatter = &logrus.JSONFormatter{}
	} else {
		if runtime.GOOS == "windows" {
			color = false
		}
		formatter = &logFormatter{logrus.TextFormatter{FullTimestamp: true, DisableColors: !color, ForceColors: color}}
	}
	return formatter
}
