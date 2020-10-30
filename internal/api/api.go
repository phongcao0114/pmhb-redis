package api

import (
	"net/http"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/app/handlers"
	"pmhb-redis/internal/app/services"

	"github.com/go-redis/redis"
)

type (
	middleware = func(http.Handler) http.Handler
	route      struct {
		desc        string
		path        string
		method      string
		handler     http.HandlerFunc
		middlewares []middleware
	}
)

// CreateEmployeeHandler function
func CreateEmployeeHandler(conf *config.Configs, redisClient *redis.ClusterClient) *handlers.EmployeeHandler {
	redisConn := services.NewRedisConnectorSrv(conf, redisClient)

	//repo := repositories.NewEmployeesRepo(conf, )
	srv := services.NewEmployeesService(conf, redisConn)
	return handlers.NewEmployeeHandler(conf, srv)
}
