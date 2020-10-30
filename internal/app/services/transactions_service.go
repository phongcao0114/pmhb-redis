package services

import (
	"context"
	"encoding/json"
	"errors"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/app/models"
	"pmhb-redis/internal/kerrors"
	"pmhb-redis/internal/pkg/klog"
	"time"

	"github.com/go-redis/redis"
)

const (
	// EmployeesServicePrefix prefix logger
	EmployeesServicePrefix = "Employees_service"
)

type (
	// EmployeesSrv groups all employees service together
	EmployeesSrv struct {
		conf   *config.Configs
		errSrv kerrors.KError
		logger klog.Logger

		//employeesRepo repositories.EmployeesRepository
		redisConn RedisConnectorService
	}

	//EmployeesService interface
	EmployeesService interface {
		GetEmployee(ctx context.Context, key string) (models.Employee, error)
		SetEmployee(ctx context.Context, key string, employee models.Employee, expiryTime int) (models.Employee, error)
	}
)

//NewEmployeesService init a new employees service
func NewEmployeesService(conf *config.Configs, redisConn RedisConnectorService) *EmployeesSrv {
	return &EmployeesSrv{
		conf:   conf,
		errSrv: kerrors.WithPrefix(EmployeesServicePrefix),
		logger: klog.WithPrefix(EmployeesServicePrefix),

		//employeesRepo: repo,
		redisConn: redisConn,
	}
}

// GetEmployees function service
func (tr *EmployeesSrv) GetEmployee(ctx context.Context, key string) (models.Employee, error) {
	employee := models.Employee{}
	if key == "" {
		return employee, tr.errSrv.Wrap(errors.New(kerrors.KeyMissingInRedis.String()), kerrors.KeyMissingInRedis, nil)
	}

	val, err := tr.redisConn.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return employee, tr.errSrv.Wrap(err, kerrors.NotFoundInRedis, nil)
		}
		return employee, tr.errSrv.Wrap(err, kerrors.CannotGetRedis, nil)
	}
	err = json.Unmarshal([]byte(val), &employee)
	if err != nil {
		return employee, tr.errSrv.Wrap(err, kerrors.UnMarshalFail, nil)
	}
	return employee, nil
}

// SetEmployee function service
func (tr *EmployeesSrv) SetEmployee(ctx context.Context, key string, employee models.Employee, expiryTime int) (models.Employee, error) {
	if key == "" {
		return employee, tr.errSrv.Wrap(errors.New(kerrors.KeyMissingInRedis.String()), kerrors.KeyMissingInRedis, nil)
	}
	if employee.Name == "" {
		return employee, tr.errSrv.Wrap(errors.New(kerrors.EmployeeNameMissing.String()), kerrors.EmployeeNameMissing, nil)
	}
	if employee.Position == "" {
		return employee, tr.errSrv.Wrap(errors.New(kerrors.EmployeePositionMissing.String()), kerrors.EmployeePositionMissing, nil)
	}
	err := tr.redisConn.Set(ctx, key, employee, time.Duration(expiryTime)*time.Second)
	if err != nil {
		return employee, tr.errSrv.Wrap(err, kerrors.CannotSetRedis, nil)
	}
	return employee, nil
}
