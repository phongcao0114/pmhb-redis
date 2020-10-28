package klog

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// logrusLogger is the wrapper of logrus logger to satisfy the klog.Logger interface.
type logrusLogger struct {
	prefix string // Holds the string prefix of logger
	*logrus.Entry
}

// newLogrusLogger creates a new logrus logger based on provided configurations.
func newLogrusLogger(cfg *Config) (Logger, error) {
	// Log level
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(lvl)

	// Log format
	switch strings.ToLower(cfg.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
		})
	}

	// Log output
	output := strings.ToLower(cfg.Output)
	if output == "stdout" {
		logger.SetOutput(os.Stdout)
	} else if output == "discard" {
		logger.SetOutput(ioutil.Discard)
	} else if strings.HasPrefix(output, "file://") {
		if isWriterAlreadyOpened(output) {
			return nil, fmt.Errorf("log writer already opened on %s", output)
		}
		f, err := os.OpenFile(strings.TrimPrefix(output, "file://"), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(f)
		trackOpeningLogWriter(output, f)
	} else {
		logger.SetOutput(os.Stderr)
	}
	return &logrusLogger{Entry: logrus.NewEntry(logger)}, nil
}

// WithFields allows to add additional fields to every log record written by returning logger.
func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		prefix: l.prefix,
		Entry:  logrus.NewEntry(l.Entry.Logger).WithFields(l.Entry.Data).WithFields(fields),
	}
}

// WithPrefix allows to add prefix to every log record written by the returning logger.
func (l *logrusLogger) WithPrefix(prefix string) Logger {
	if l.prefix != "" {
		prefix = l.prefix + "/" + prefix
	}
	return &logrusLogger{
		prefix: prefix,
		Entry:  l.Entry,
	}
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Debugf(format, args...)
}
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Infof(format, args...)
}
func (l *logrusLogger) Printf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Printf(format, args...)
}
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Warnf(format, args...)
}
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Errorf(format, args...)
}
func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	if l.prefix != "" {
		format = l.prefix + ": " + format
	}
	l.Entry.Panicf(format, args...)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Debug(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Debug(args...)
}
func (l *logrusLogger) Info(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Info(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Info(args...)
}
func (l *logrusLogger) Print(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Print(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Print(args...)
}
func (l *logrusLogger) Warn(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Warn(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Warn(args...)
}
func (l *logrusLogger) Error(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Error(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Error(args...)
}
func (l *logrusLogger) Panic(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Panic(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Panic(args...)
}

func (l *logrusLogger) Debugln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Debugln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Debugln(args...)
}
func (l *logrusLogger) Infoln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Infoln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Infoln(args...)

}
func (l *logrusLogger) Println(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Println(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Println(args...)
}
func (l *logrusLogger) Warnln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Warnln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Warnln(args...)
}
func (l *logrusLogger) Errorln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Errorln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Errorln(args...)
}
func (l *logrusLogger) Panicln(args ...interface{}) {
	if l.prefix != "" {
		l.Entry.Panicln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Entry.Panicln(args...)
}

func (l *logrusLogger) KDebugf(ctx context.Context, format string, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Debugf(l.getPrefixedFormat(format), args...)
		return
	}
	l.Debugf(l.getPrefixedFormat(format), args...)
}
func (l *logrusLogger) KInfof(ctx context.Context, format string, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Infof(l.getPrefixedFormat(format), args...)
		return
	}
	l.Infof(l.getPrefixedFormat(format), args...)
}
func (l *logrusLogger) KPrintf(ctx context.Context, format string, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Printf(l.getPrefixedFormat(format), args...)
		return
	}
	l.Printf(l.getPrefixedFormat(format), args...)
}
func (l *logrusLogger) KWarnf(ctx context.Context, format string, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Warnf(l.getPrefixedFormat(format), args...)
		return
	}
	l.Warnf(l.getPrefixedFormat(format), args...)
}
func (l *logrusLogger) KErrorf(ctx context.Context, format string, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Errorf(l.getPrefixedFormat(format), args...)
		return
	}
	l.Errorf(l.getPrefixedFormat(format), args...)
}
func (l *logrusLogger) KPanicf(ctx context.Context, format string, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Panicf(l.getPrefixedFormat(format), args...)
		return
	}
	l.Panicf(l.getPrefixedFormat(format), args...)
}

func (l *logrusLogger) KDebug(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Debug(prefixHelper(l.prefix, args)...)
		return
	}
	l.Debug(args...)
}
func (l *logrusLogger) KInfo(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Info(prefixHelper(l.prefix, args)...)
		return
	}
	l.Info(args...)
}
func (l *logrusLogger) KPrint(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Print(prefixHelper(l.prefix, args)...)
		return
	}
	l.Print(args...)
}
func (l *logrusLogger) KWarn(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Warn(prefixHelper(l.prefix, args)...)
		return
	}
	l.Warn(args...)
}
func (l *logrusLogger) KError(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Error(prefixHelper(l.prefix, args)...)
		return
	}
	l.Error(args...)
}
func (l *logrusLogger) KPanic(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Panic(prefixHelper(l.prefix, args)...)
		return
	}
	l.Panic(args...)
}

func (l *logrusLogger) KDebugln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Debugln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Debugln(args...)
}
func (l *logrusLogger) KInfoln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Infoln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Infoln(args...)

}
func (l *logrusLogger) KPrintln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Println(prefixHelper(l.prefix, args)...)
		return
	}
	l.Println(args...)

}
func (l *logrusLogger) KWarnln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Warnln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Warnln(args...)
}
func (l *logrusLogger) KErrorln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Errorln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Errorln(args...)
}
func (l *logrusLogger) KPanicln(ctx context.Context, args ...interface{}) {
	if lg := l.getLoggerFromContext(ctx); lg != nil {
		lg.Panicln(prefixHelper(l.prefix, args)...)
		return
	}
	l.Panicln(args...)
}

func (l *logrusLogger) getLoggerFromContext(ctx context.Context) *logrusLogger {
	if ctx == nil {
		return nil
	}

	lg := ctx.Value(KLogger)
	if lg == nil {
		return nil
	}
	if _, ok := lg.(*logrusLogger); !ok {
		return nil
	}
	logger := lg.(*logrusLogger)
	logger = logger.WithFields(l.Entry.Data).WithFields(logger.Entry.Data).(*logrusLogger)

	return logger
}

func (l *logrusLogger) getPrefixedFormat(format string) string {
	if l.prefix != "" {
		return l.prefix + ": " + format
	}
	return format
}

func prefixHelper(prefix interface{}, s []interface{}) []interface{} {
	if len(s) == 0 {
		return []interface{}{prefix}
	}
	s[0] = fmt.Sprintf("%s: %v", prefix, s[0])
	return s
}

// NewContext return new context with logger inside
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, KLogger, logger)
}
