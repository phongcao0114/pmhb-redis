package services

import (
	"context"
	"encoding/json"
	"errors"
	"pmhb-redis/internal/app/config"
	"pmhb-redis/internal/pkg/klog"
	"testing"

	repoMock "pmhb-redis/internal/app/repositories/mock_repo"
	srvMock "pmhb-redis/internal/app/services/mock_services"
	kafkaMock "pmhb-redis/internal/kafka/mock_kafka"

	"pmhb-redis/internal/app/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactions(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	logger := klog.WithPrefix("test-GetTransactions")
	redisSrv := srvMock.NewMockRedisConnectorService(ctl)
	kafkaProducer := kafkaMock.NewMockProducerKafka(ctl)
	transRepo := repoMock.NewMockTransactionsRepository(ctl)

	conf := &config.Configs{}

	transSrv := NewTransactionsService(conf, kafkaProducer, transRepo, redisSrv)

	scenarios := GetTransactionsCases()

	for _, scena := range scenarios {
		switch scena.Name {
		case 0:
			transRepo.EXPECT().GetTransaction(
				gomock.Any(),
				gomock.Any(),
			).Return(
				scena.ExpectGetTransactions.Model,
				scena.ExpectGetTransactions.Err,
			)

			outputModel, outputErr := transSrv.GetTransactions(scena.Input.Context, &scena.Input.Body)
			expectModel, expectErr := scena.ExpectGetTransactions.Model, scena.ExpectGetTransactions.Err

			if o, e := outputErr, expectErr; o != e {
				logger.KErrorf(scena.Input.Context, "Failed Error checking: ouput-%v  expect-%v", o, e)
				t.Fail()
			}
			buff1, errJSON1 := json.Marshal(outputModel)
			buff2, errJSON2 := json.Marshal(expectModel)
			if errJSON1 != nil || errJSON2 != nil {
				t.Fail()
			}

			assert.JSONEq(t, string(buff2), string(buff1), "two data are not same")

		case 1:
			transRepo.EXPECT().GetTransaction(
				gomock.Any(),
				gomock.Any(),
			).Return(
				scena.ExpectGetTransactions.Model,
				scena.ExpectGetTransactions.Err,
			)

			outputModel, outputErr := transSrv.GetTransactions(scena.Input.Context, &scena.Input.Body)
			expectModel, expectErr := scena.ExpectGetTransactions.Model, scena.ExpectGetTransactions.Err

			if o, e := outputErr, expectErr; o != e {
				logger.KErrorf(scena.Input.Context, "Failed Error checking: ouput-%v  expect-%v", o, e)
				t.Fail()
			}
			buff1, errJSON1 := json.Marshal(outputModel)
			buff2, errJSON2 := json.Marshal(expectModel)
			if errJSON1 != nil || errJSON2 != nil {
				t.Fail()
			}

			assert.JSONEq(t, string(buff2), string(buff1), "two data are not same")
		}
	}
}

type (
	InputGetTransactions struct {
		Context context.Context
		Body    models.GetTransactionSrvReq
	}
	ExpectGetTransactions struct {
		Model []models.Transactions
		Err   error
	}
	ExpectTransGetTransactions struct {
		Model []models.Transactions
		Err   error
	}
	TestGetTransactionsCase struct {
		Name        int
		Description string
		Input       InputGetTransactions
		ExpectTransGetTransactions
		ExpectGetTransactions
	}
)

func GetTransactionsCases() []TestGetTransactionsCase {
	var output []TestGetTransactionsCase

	//===================== 1. Happy case ===================== //
	case1 := TestGetTransactionsCase{
		Name:        0,
		Description: "Successfully calling to get transaction",
		Input: InputGetTransactions{
			Body: models.GetTransactionSrvReq{
				TransactionID: 1,
			},
		},
		ExpectTransGetTransactions: ExpectTransGetTransactions{
			Model: []models.Transactions{},
			Err:   nil,
		},
		ExpectGetTransactions: ExpectGetTransactions{
			Model: []models.Transactions{},
			Err:   nil,
		},
	}
	case1.ExpectTransGetTransactions.Model = append(case1.ExpectTransGetTransactions.Model, models.Transactions{
		TransactionID:   1,
		TransactionName: "transaction_name_1",
	})
	case1.ExpectGetTransactions.Model = append(case1.ExpectGetTransactions.Model, models.Transactions{
		TransactionID:   1,
		TransactionName: "transaction_name_1",
	})
	output = append(output, case1)

	//===================== 2. Database Server Error Case ===================== //
	case2 := TestGetTransactionsCase{
		Name:        1,
		Description: "Err calling to get transaction",
		Input: InputGetTransactions{
			Body: models.GetTransactionSrvReq{
				TransactionID: 2,
			},
		},
		ExpectTransGetTransactions: ExpectTransGetTransactions{
			Model: []models.Transactions{},
			Err:   errors.New("Database Server Error"),
		},
		ExpectGetTransactions: ExpectGetTransactions{
			Model: []models.Transactions{},
			Err:   errors.New("Database Server Error"),
		},
	}
	output = append(output, case2)

	return output
}

func TestInsertTransaction(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	logger := klog.WithPrefix("test-InsertTransactions")
	redisSrv := srvMock.NewMockRedisConnectorService(ctl)
	kafkaProducer := kafkaMock.NewMockProducerKafka(ctl)
	transRepo := repoMock.NewMockTransactionsRepository(ctl)

	conf := &config.Configs{
		AppID: "789",
	}

	trasSrv := NewTransactionsService(conf, kafkaProducer, transRepo, redisSrv)

	scenarios := InsertTransactionsCases()

	for _, scena := range scenarios {
		switch scena.Name {
		case 0:
			transRepo.EXPECT().InsertTransaction(
				gomock.Any(),
				gomock.Any(),
			).Return(
				scena.ExpectTransInsertTransactions.TransactionID,
				scena.ExpectTransInsertTransactions.Err,
			)

			kafkaProducer.EXPECT().Send(
				gomock.Any(),
				gomock.Any(),
			).AnyTimes()

			outputModel, outputErr := trasSrv.InsertTransaction(scena.Input.Context, &scena.Input.Body)
			expectModel, expectErr := scena.ExpectInsertTransactions.Model, scena.ExpectInsertTransactions.Err

			if o, e := outputErr, expectErr; o != e {
				logger.KErrorf(scena.Input.Context, "Failed Error checking: ouput-%v  expect-%v", o, e)
				t.Fail()
			}
			buff1, errJSON1 := json.Marshal(outputModel)
			buff2, errJSON2 := json.Marshal(expectModel)
			if errJSON1 != nil || errJSON2 != nil {
				t.Fail()
			}

			assert.JSONEq(t, string(buff2), string(buff1), "two data are not same")

		case 1:
			transRepo.EXPECT().InsertTransaction(
				gomock.Any(),
				gomock.Any(),
			).Return(
				scena.ExpectTransInsertTransactions.TransactionID,
				scena.ExpectTransInsertTransactions.Err,
			)

			outputModel, outputErr := trasSrv.InsertTransaction(scena.Input.Context, &scena.Input.Body)
			expectModel, expectErr := scena.ExpectInsertTransactions.Model, scena.ExpectInsertTransactions.Err

			if o, e := outputErr, expectErr; o.Error() != e.Error() {
				logger.KErrorf(scena.Input.Context, "Failed Error checking: ouput-%v  expect-%v", o, e)
				t.Fail()
			}
			buff1, errJSON1 := json.Marshal(outputModel)
			buff2, errJSON2 := json.Marshal(expectModel)
			if errJSON1 != nil || errJSON2 != nil {
				t.Fail()
			}

			assert.JSONEq(t, string(buff2), string(buff1), "two data are not same")
		}
	}
}

type (
	InputInsertTransactions struct {
		Context context.Context
		Body    models.InsertTransactionSrvReq
	}
	ExpectInsertTransactions struct {
		Model models.InsertTransactionSrvRes
		Err   error
	}
	ExpectTransInsertTransactions struct {
		TransactionID int64
		Err           error
	}
	TestInsertTransactionsCase struct {
		Name        int
		Description string
		Input       InputInsertTransactions
		ExpectTransInsertTransactions
		ExpectInsertTransactions
	}
)

func InsertTransactionsCases() []TestInsertTransactionsCase {
	var output []TestInsertTransactionsCase

	//===================== 1. Happy case ===================== //
	case1 := TestInsertTransactionsCase{
		Name:        0,
		Description: "Successfully calling to get transaction",
		Input: InputInsertTransactions{
			Body: models.InsertTransactionSrvReq{
				TransactionName: "transaction_name_1",
			},
		},
		ExpectTransInsertTransactions: ExpectTransInsertTransactions{
			TransactionID: 1,
			Err:           nil,
		},
		ExpectInsertTransactions: ExpectInsertTransactions{
			Model: models.InsertTransactionSrvRes{
				TransactionID:   1,
				TransactionName: "transaction_name_1",
			},
			Err: nil,
		},
	}

	output = append(output, case1)

	//===================== 2. Database Server Error Case ===================== //
	case2 := TestInsertTransactionsCase{
		Name:        1,
		Description: "Err calling to get transaction",
		Input: InputInsertTransactions{
			Body: models.InsertTransactionSrvReq{
				TransactionName: "transaction_name_2",
			},
		},
		ExpectTransInsertTransactions: ExpectTransInsertTransactions{
			Err: errors.New("Database Server Error"),
		},
		ExpectInsertTransactions: ExpectInsertTransactions{
			Model: models.InsertTransactionSrvRes{},
			Err:   errors.New("Database Server Error"),
		},
	}
	output = append(output, case2)

	return output
}
