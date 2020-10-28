package app

import (
	"fmt"
	"net/http"
	"pmhb-redis/internal/api"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/pkg/db/redis"
)

func UseRedis() Option {
	return func(o *Options) error {
		o.Redis = redis.CreateRedisPool(o.Config)
		return nil
	}
}

func UseRouter() Option {
	return func(o *Options) error {
		router, err := api.NewRouter(o.Config, &api.InfraConn{
			//DBconn:      o.Database,
			RedisClient: o.Redis,
		})
		if err != nil {
			return err
		}

		srv := &http.Server{
			Addr:              fmt.Sprint(":", o.Config.HTTPServerPort),
			Handler:           router,
			ReadTimeout:       o.Config.HTTPServer.ReadTimeout,
			WriteTimeout:      o.Config.HTTPServer.WriteTimeout,
			ReadHeaderTimeout: o.Config.HTTPServer.ReadHeaderTimeout,
		}

		o.Server = srv
		return nil
	}
}

func SetConfig(conf *config.Configs) Option {
	return func(o *Options) error {
		o.Config = conf
		return nil
	}
}
