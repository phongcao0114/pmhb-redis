package klog

import (
	"context"
	"log"

	"github.com/spf13/viper"
)

const KLogger = "klogger"

type (
	// Logger is the API wrapper for underlying logging libraries.
	Logger interface {
		WithFields(fields map[string]interface{}) Logger
		WithPrefix(prefix string) Logger

		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Printf(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Panicf(format string, args ...interface{})

		Debug(args ...interface{})
		Info(args ...interface{})
		Print(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Panic(args ...interface{})

		Debugln(args ...interface{})
		Infoln(args ...interface{})
		Println(args ...interface{})
		Warnln(args ...interface{})
		Errorln(args ...interface{})
		Panicln(args ...interface{})

		KDebugf(ctx context.Context, format string, args ...interface{})
		KInfof(ctx context.Context, format string, args ...interface{})
		KPrintf(ctx context.Context, format string, args ...interface{})
		KWarnf(ctx context.Context, format string, args ...interface{})
		KErrorf(ctx context.Context, format string, args ...interface{})
		KPanicf(ctx context.Context, format string, args ...interface{})

		KDebug(ctx context.Context, args ...interface{})
		KInfo(ctx context.Context, args ...interface{})
		KPrint(ctx context.Context, args ...interface{})
		KWarn(ctx context.Context, args ...interface{})
		KError(ctx context.Context, args ...interface{})
		KPanic(ctx context.Context, args ...interface{})

		KDebugln(ctx context.Context, args ...interface{})
		KInfoln(ctx context.Context, args ...interface{})
		KPrintln(ctx context.Context, args ...interface{})
		KWarnln(ctx context.Context, args ...interface{})
		KErrorln(ctx context.Context, args ...interface{})
		KPanicln(ctx context.Context, args ...interface{})
	}

	// Config is the configuration of logger.
	Config struct {
		// Log level.
		// Can be one of: debug, info, warn, error, panic
		Level string `mapstructure:"level"`
		// Where log will be written to.
		// Can be one of: stdout, stderr, discard, file://path/to/log/file
		Output string `mapstructure:"output"`
		// Log output format.
		// Can be one of: text, json
		Format string `mapstructure:"format"`
	}
)

// Standard logger that will be used as the default logger for this package.
var std Logger

func init() {
	cfg, err := readLogConfig()
	if err != nil {
		log.Println("klog: failed to load log configs, initialize using default values")
	}
	log.Printf("klog: config: %#v", cfg)

	std, err = New(cfg) // Create default logger
	if err != nil {
		panic(err)
	}
}

// readLogConfig tries to parses klog configurations from config files and environment variables.
func readLogConfig() (*Config, error) {
	vp := viper.New()
	vp.SetDefault("level", "debug")
	vp.SetDefault("format", "text")
	vp.SetDefault("output", "stderr")

	vp.AddConfigPath(".")
	vp.AddConfigPath("config")
	vp.AddConfigPath("configs")
	vp.SetConfigName("log")
	vp.SetEnvPrefix("KLOG")
	vp.AutomaticEnv()
	_ = vp.ReadInConfig()

	cfg := &Config{}
	if err := vp.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// New creates new logger by provided configurations.
// Currently creating logrus logger by default.
func New(cfg *Config) (Logger, error) {
	return newLogrusLogger(cfg)
}

// WithFields allows to add additional fields to every log record written by the returning logger.
func WithFields(fields map[string]interface{}) Logger {
	return std.WithFields(fields)
}

// WithPrefix allows to add prefix to every log record written by the returning logger.
func WithPrefix(prefix string) Logger {
	return std.WithPrefix(prefix)
}

func WithMetricType(metricType string) Logger {
	return std.WithFields(map[string]interface{}{
		"metric_type": metricType,
	})
}

// Normal logging funcs
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	std.Debug(args...)
}
func Info(args ...interface{}) {
	std.Info(args...)
}
func Print(args ...interface{}) {
	std.Print(args...)
}
func Warn(args ...interface{}) {
	std.Warn(args...)
}
func Error(args ...interface{}) {
	std.Error(args...)
}
func Panic(args ...interface{}) {
	std.Panic(args...)
}

func Debugln(args ...interface{}) {
	std.Debugln(args...)
}
func Infoln(args ...interface{}) {
	std.Infoln(args...)
}
func Println(args ...interface{}) {
	std.Println(args...)
}
func Warnln(args ...interface{}) {
	std.Warnln(args...)
}
func Errorln(args ...interface{}) {
	std.Errorln(args...)
}
func Panicln(args ...interface{}) {
	std.Panicln(args...)
}

// Logging with context
func KDebugf(ctx context.Context, format string, args ...interface{}) {
	std.KDebugf(ctx, format, args...)
}
func KInfof(ctx context.Context, format string, args ...interface{}) {
	std.KInfof(ctx, format, args...)
}
func KPrintf(ctx context.Context, format string, args ...interface{}) {
	std.KPrintf(ctx, format, args...)
}
func KWarnf(ctx context.Context, format string, args ...interface{}) {
	std.KWarnf(ctx, format, args...)
}
func KErrorf(ctx context.Context, format string, args ...interface{}) {
	std.KErrorf(ctx, format, args...)
}
func KPanicf(ctx context.Context, format string, args ...interface{}) {
	std.KPanicf(ctx, format, args...)
}

func KDebug(ctx context.Context, args ...interface{}) {
	std.KDebug(ctx, args...)
}
func KInfo(ctx context.Context, args ...interface{}) {
	std.KInfo(ctx, args...)
}
func KPrint(ctx context.Context, args ...interface{}) {
	std.KPrint(ctx, args...)
}
func KWarn(ctx context.Context, args ...interface{}) {
	std.KWarn(ctx, args...)
}
func KError(ctx context.Context, args ...interface{}) {
	std.KError(ctx, args...)
}
func KPanic(ctx context.Context, args ...interface{}) {
	std.KPanic(ctx, args...)
}

func KDebugln(ctx context.Context, args ...interface{}) {
	std.KDebugln(ctx, args...)
}
func KInfoln(ctx context.Context, args ...interface{}) {
	std.KInfoln(ctx, args...)
}
func KPrintln(ctx context.Context, args ...interface{}) {
	std.KPrintln(ctx, args...)
}
func KWarnln(ctx context.Context, args ...interface{}) {
	std.KWarnln(ctx, args...)
}
func KErrorln(ctx context.Context, args ...interface{}) {
	std.KErrorln(ctx, args...)
}
func KPanicln(ctx context.Context, args ...interface{}) {
	std.KPanicln(ctx, args...)
}
