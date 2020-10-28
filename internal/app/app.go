package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/app/utils"
	"pmhb-redis/internal/pkg/klog"

	"syscall"
	"time"

	"github.com/gomodule/redigo/redis"
)

type (
	// Option return function to modify Options value
	Option func(o *Options) error

	// OptionList interface for App
	OptionList interface {
		// Apply app options
		Use(option ...Option)
		// Start app
		Start()
	}

	optionList struct {
		options []Option
	}

	// Options model
	Options struct {
		Config *config.Configs
		Redis  *redis.Pool
		Server *http.Server
	}
)

func (oL *optionList) Use(option ...Option) {
	oL.options = append(oL.options, option...)
}

func (oL *optionList) Start() {
	// Write initializer here

	// 1. Safely close all log file writers (if exists)
	defer klog.Close()

	opts, err := GetOptions(oL)
	if err != nil {
		panic(err)
	}

	// 2. Prepare logger
	KLogger := klog.WithPrefix("main")
	KLogger.WithFields(map[string]interface{}{
		"state":  opts.Config.Stage,
		"port":   opts.Config.HTTPServerPort,
		"config": opts.Config,
	}).Info("starting server")

	// 3. Set singleton
	utils.ResponseAppID = opts.Config.AppID
	config.InitRandomProfileUserID()

	// 4. Load location
	bkkLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		KLogger.Panicf("[main] Failed to load location, err: %v", err)
	}
	utils.BKKLocation = bkkLocation

	srvCtx, srvCancel := context.WithCancel(context.Background())

	// If Server exist
	if opts.Server != nil {
		go func() {
			if err := opts.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				KLogger.Errorf("listen: %s\n", err)
			}
		}()
	}

	// 12. Graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-signals
	KLogger.Info("shutting down http server...")
	if err := opts.Server.Shutdown(srvCtx); err != nil {
		KLogger.Panicln("http server shutdown with error:", err)
	}
	srvCancel()
}

// Initialize App
func New() OptionList {
	return &optionList{}
}

func GetDefaultOptions() Options {
	return Options{
		Config: &config.Configs{},
		//Database: &db.DB{},
		Redis: nil,
	}
}

func GetOptions(opts *optionList) (Options, error) {
	var options Options
	options = GetDefaultOptions()
	for _, opt := range opts.options {
		if opt != nil {
			if err := opt(&options); err != nil {
				return Options{}, err
			}
		}
	}
	return options, nil
}
