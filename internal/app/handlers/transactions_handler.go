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

	"fmt"
)

const (
	// EmployeeHandlerPrefix prefix logger
	EmployeeHandlerPrefix = "Employee_handler"
)

// EmployeeHandler struct defines the variables for specifying interface.
type EmployeeHandler struct {
	conf       *config.Configs
	errHandler kerrors.KError
	logger     klog.Logger

	srv services.EmployeesService
}

// NewEmployeeHandler connects to the service from handler.
func NewEmployeeHandler(conf *config.Configs, s services.EmployeesService) *EmployeeHandler {
	return &EmployeeHandler{
		conf:       conf,
		errHandler: kerrors.WithPrefix(EmployeeHandlerPrefix),
		logger:     klog.WithPrefix(EmployeeHandlerPrefix),

		srv: s,
	}
}

// GetEmployee handler handles the upcoming request.
func (th *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	var req models.RequestInfo
	var body models.GetEmployeeReq
	err := utils.DecodeToBody(&th.errHandler, &req, &body, r)
	if err != nil {
		response.WriteJSON(w)(response.HandleError(r, req.Header, err))
		return
	}
	fmt.Printf("req: %+v\n", req)
	employee, err := th.srv.GetEmployee(r.Context(), body.Key)
	if err != nil {
		response.WriteJSON(w)(response.HandleError(r, req.Header, err))
		return
	}
	response.WriteJSON(w)(response.HandleSuccess(r, req.Header, employee))
	return
}

// InsertEmployee handler handles the upcoming request.
func (th *EmployeeHandler) InsertEmployee(w http.ResponseWriter, r *http.Request) {

	var req models.RequestInfo
	body := models.SetEmployeeReq{}
	err := utils.DecodeToBody(&th.errHandler, &req, &body, r)
	if err != nil {
		response.WriteJSON(w)(response.HandleError(r, req.Header, err))
		return
	}

	commitModels, err := th.srv.SetEmployee(r.Context(), body.Key, body.Employee, body.ExpiryTime)
	if err != nil {
		response.WriteJSON(w)(response.HandleError(r, req.Header, err))
		return
	}

	response.WriteJSON(w)(response.HandleSuccess(r, req.Header, commitModels))
	return

}
