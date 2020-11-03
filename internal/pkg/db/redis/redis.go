package redis

import (
	"pmhb-redis/internal/app/config"

	"github.com/go-redis/redis"
)

func CreateRedisClusterClient(sconf *config.Configs) (*redis.ClusterClient, error) {

	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    sconf.Redis.Addresses,
		Password: sconf.Redis.Password,
	})

	if err := c.Ping().Err(); err != nil {
		return nil, err
	}
	return c, nil

}
