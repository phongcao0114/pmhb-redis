// Package services contains the services of the system.
package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/kerrors"
	"pmhb-redis/internal/pkg/klog"
	"sync"

	"github.com/gomodule/redigo/redis"
)

const (
	// RedisConnectorServicePrefix constant
	RedisConnectorServicePrefix = "Redis_connector_service"
)

type (
	// RedisConnectorService interface
	RedisConnectorService interface {
		Get(ctx context.Context, key string) ([]byte, error)
		Set(ctx context.Context, key string, data interface{}) error
		MSet(ctx context.Context, keys []interface{}, data []interface{}) error
		MSetWithExpired(ctx context.Context, keys []interface{}, data []interface{}, seconds int) error
		SetWithExpired(ctx context.Context, key string, data interface{}, seconds int) error
		Delete(ctx context.Context, key string) error
	}

	// RedisConnectorSrv structure
	RedisConnectorSrv struct {
		errSvr kerrors.KError
		conf   *config.Configs
		redis  *redis.Pool
		logger klog.Logger
		mutex  sync.Mutex
	}
)

// NewRedisConnectorSrv declare new connection
func NewRedisConnectorSrv(conf *config.Configs, redis *redis.Pool) *RedisConnectorSrv {
	return &RedisConnectorSrv{
		errSvr: kerrors.WithPrefix(RedisConnectorServicePrefix),
		conf:   conf,
		redis:  redis,
		logger: klog.WithPrefix(RedisConnectorServicePrefix),
		mutex:  sync.Mutex{},
	}
}

// Get function
func (r *RedisConnectorSrv) Get(ctx context.Context, key string) ([]byte, error) {
	conn := r.redis.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		r.logger.KError(ctx, err.Error())
		return nil, err
	}
	r.logger.KDebug(ctx, string(data))
	return data, nil
}

// Set function
func (r *RedisConnectorSrv) Set(ctx context.Context, key string, data interface{}) error {
	// serialize data object to JSON
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	conn := r.redis.Get()
	defer conn.Close()
	reply, err := conn.Do("SET", key, json)
	if err != nil {
		r.logger.KError(ctx, err.Error())
		return err
	}
	_, err = conn.Do("EXPIRE", key, r.conf.Redis.TTL)
	if err != nil {
		return fmt.Errorf("error setting ttl %v", err)
	}
	r.logger.KDebug(ctx, reply)
	return nil
}

// SetWithExpired function
func (r *RedisConnectorSrv) SetWithExpired(ctx context.Context, key string, data interface{}, seconds int) error {
	// serialize data object to JSON
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	conn := r.redis.Get()
	defer conn.Close()
	reply, err := conn.Do("SET", key, json)
	if err != nil {
		r.logger.KError(ctx, err.Error())
		return err
	}
	_, err = conn.Do("EXPIRE", key, seconds)
	if err != nil {
		return fmt.Errorf("error setting ttl %v", err)
	}
	r.logger.KDebug(ctx, reply)
	return nil
}

// Delete function
func (r *RedisConnectorSrv) Delete(ctx context.Context, key string) error {
	conn := r.redis.Get()
	defer conn.Close()
	reply, err := conn.Do("DEL", key)
	if err != nil {
		r.logger.KError(ctx, err.Error())
		return err
	}
	r.logger.KDebug(ctx, reply)
	return nil
}

// MSet function
func (r *RedisConnectorSrv) MSet(ctx context.Context, keys []interface{}, data []interface{}) error {
	if len(keys) != len(data) {
		return errors.New("Length of keys and values are not equals")
	}

	// Open connection
	conn := r.redis.Get()
	defer conn.Close()

	var pairs []interface{}
	for idx, i := range data {
		json, err := json.Marshal(i)
		if err != nil {
			return err
		}
		pairs = append(pairs, keys[idx], json)
	}
	reply, err := conn.Do("MSET", pairs...)
	if err != nil {
		r.logger.KError(ctx, err.Error())
		return err
	}
	_, err = conn.Do("EXPIRE", pairs, r.conf.Redis.TTL)
	if err != nil {
		return fmt.Errorf("error setting ttl %v", err)
	}
	r.logger.KDebug(ctx, reply)
	return nil
}

// MSetWithExpired function
func (r *RedisConnectorSrv) MSetWithExpired(ctx context.Context, keys []interface{}, data []interface{}, seconds int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(keys) != len(data) {
		return errors.New("Length of keys and values are not equals")
	}

	// Open connection
	conn := r.redis.Get()
	defer conn.Close()

	var pairs []interface{}
	for idx, i := range data {
		json, err := json.Marshal(i)
		if err != nil {
			return err
		}
		pairs = append(pairs, keys[idx], json)
	}
	reply, err := conn.Do("MSET", pairs...)
	if err != nil {
		r.logger.KError(ctx, err.Error())
		return err
	}

	for _, i := range keys {
		_, err = conn.Do("EXPIRE", i, seconds)
		if err != nil {
			return fmt.Errorf("error setting ttl %v", err)
		}
	}
	r.logger.KDebug(ctx, reply)
	return nil
}
