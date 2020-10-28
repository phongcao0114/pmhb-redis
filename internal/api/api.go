package api

import (
	"net/http"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/app/handlers"
	"pmhb-redis/internal/app/services"

	redigo "github.com/gomodule/redigo/redis"
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

// CreateTransactionHandler function
func CreateTransactionHandler(conf *config.Configs, redisClient *redigo.Pool) *handlers.TransactionHandler {
	redisConn := services.NewRedisConnectorSrv(conf, redisClient)

	//repo := repositories.NewTransactionsRepo(conf, )
	srv := services.NewTransactionsService(conf, redisConn)
	return handlers.NewTransactionHandler(conf, srv)
}
