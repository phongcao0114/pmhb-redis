package redis

import (
	"errors"
	"pmhb-redis/internal/app/config"
	"time"

	"github.com/FZambia/sentinel"
	redigo "github.com/gomodule/redigo/redis"
)

func CreateRedisPool(sconf *config.Configs) *redigo.Pool {
	// Prepare Sentinel
	sntnl := &sentinel.Sentinel{
		Addrs:      sconf.Redis.Addresses,
		MasterName: sconf.Redis.MasterName,
		Dial: func(addr string) (redigo.Conn, error) {
			timeout := sconf.Redis.DialTimeout
			c, err := redigo.Dial("tcp", addr, redigo.DialConnectTimeout(timeout), redigo.DialReadTimeout(timeout), redigo.DialWriteTimeout(timeout))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	// Return redis Pool
	return &redigo.Pool{
		MaxIdle:     sconf.Redis.MaxIdle,
		MaxActive:   sconf.Redis.MaxActive,
		Wait:        true,
		IdleTimeout: sconf.Redis.IdleTimeout,
		Dial: func() (redigo.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			if sconf.Redis.Password != "" {
				c, err := redigo.Dial("tcp", masterAddr, redigo.DialPassword(sconf.Redis.Password))
				if err != nil {
					return nil, err
				}
				return c, nil
			}

			c, err := redigo.Dial("tcp", masterAddr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("role check failed")
			}
			if time.Since(t) < 10*time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
