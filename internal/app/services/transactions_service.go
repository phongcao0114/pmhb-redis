package services

import (
	"context"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/app/models"
	"pmhb-redis/internal/kerrors"
	"pmhb-redis/internal/pkg/klog"
)

const (
	// TransactionsServicePrefix prefix logger
	TransactionsServicePrefix = "Transactions_service"
)

type (
	// TransactionsSrv groups all transactions service together
	TransactionsSrv struct {
		conf   *config.Configs
		errSrv kerrors.KError
		logger klog.Logger

		//transactionsRepo repositories.TransactionsRepository
		redisConn RedisConnectorService
	}

	//TransactionsService interface
	TransactionsService interface {
		GetTransactions(ctx context.Context, req *models.GetTransactionSrvReq) ([]models.Transactions, error)
		InsertTransaction(ctx context.Context, req *models.InsertTransactionSrvReq) (models.InsertTransactionSrvRes, error)
	}
)

//NewTransactionsService init a new transactions service
func NewTransactionsService(conf *config.Configs, redisConn RedisConnectorService) *TransactionsSrv {
	return &TransactionsSrv{
		conf:   conf,
		errSrv: kerrors.WithPrefix(TransactionsServicePrefix),
		logger: klog.WithPrefix(TransactionsServicePrefix),

		//transactionsRepo: repo,
		redisConn: redisConn,
	}
}

// GetTransactions function service
func (tr *TransactionsSrv) GetTransactions(ctx context.Context, req *models.GetTransactionSrvReq) ([]models.Transactions, error) {
	return []models.Transactions{}, nil
}

// InsertTransaction function service
func (tr *TransactionsSrv) InsertTransaction(ctx context.Context, req *models.InsertTransactionSrvReq) (models.InsertTransactionSrvRes, error) {
	return models.InsertTransactionSrvRes{}, nil
}
