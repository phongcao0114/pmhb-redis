package handlers

import (
	"net/http"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/app/models"
	"pmhb-redis/internal/app/response"
	"pmhb-redis/internal/app/services"
	"pmhb-redis/internal/app/utils"
	"pmhb-redis/internal/kerrors"
	"pmhb-redis/internal/pkg/klog"
)

const (
	// TransactionHandlerPrefix prefix logger
	TransactionHandlerPrefix = "Transaction_handler"
)

// TransactionHandler struct defines the variables for specifying interface.
type TransactionHandler struct {
	conf       *config.Configs
	errHandler kerrors.KError
	logger     klog.Logger

	srv services.TransactionsService
}

// NewTransactionHandler connects to the service from handler.
func NewTransactionHandler(conf *config.Configs, s services.TransactionsService) *TransactionHandler {
	return &TransactionHandler{
		conf:       conf,
		errHandler: kerrors.WithPrefix(TransactionHandlerPrefix),
		logger:     klog.WithPrefix(TransactionHandlerPrefix),

		srv: s,
	}
}

// GetTransaction handler handles the upcoming request.
func (th *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {

	var req models.RequestInfo
	var body models.GetTransactionReq
	err := utils.DecodeToBody(&th.errHandler, &req, &body, r)
	if err != nil {
		response.WriteJSON(w)(response.HandleError(r, req.Header, err))
		return
	}

	list, err := th.srv.GetTransactions(r.Context(), &models.GetTransactionSrvReq{
		TransactionID: body.TransactionID,
	})
	if err != nil {
		response.WriteJSON(w)(response.HandleError(r, req.Header, err))
		return
	}

	commitModels := models.GetTransactionRes{
		ListTransactions: list,
	}

	response.WriteJSON(w)(response.HandleSuccess(r, req.Header, commitModels))
	return

}

// InsertTransaction handler handles the upcoming request.
func (th *TransactionHandler) InsertTransaction(w http.ResponseWriter, r *http.Request) {

	var req models.RequestInfo
	body := models.InsertTransactionReq{}
	err := utils.DecodeToBody(&th.errHandler, &req, &body, r)
	if err != nil {
		response.WriteJSON(w)(response.HandleError(r, req.Header, err))
		return
	}

	commitModels, err := th.srv.InsertTransaction(r.Context(), &models.InsertTransactionSrvReq{
		TransactionName: body.TransactionName,
	})
	if err != nil {
		response.WriteJSON(w)(response.HandleError(r, req.Header, err))
		return
	}

	response.WriteJSON(w)(response.HandleSuccess(r, req.Header, commitModels))
	return

}
