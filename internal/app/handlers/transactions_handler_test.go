package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pmhb-redis/internal/app/models"
	"pmhb-redis/internal/app/response"
	"pmhb-redis/internal/kerrors"

	"pmhb-redis/internal/app/config"
	transSrvMock "pmhb-redis/internal/app/services/mock_services"

	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactions(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	transSrv := transSrvMock.NewMockTransactionsService(ctl)
	conf := &config.Configs{}

	transHandler := NewTransactionHandler(conf, transSrv)

	scenarios := GetTransactionsCases()

	for _, scena := range scenarios {
		inputJSON, errInput := json.Marshal(scena.Input)
		if errInput != nil {
			t.Fail()
		}

		req, errRequest := http.NewRequest("POST", "/kph/api/get", bytes.NewBuffer(inputJSON))
		if errRequest != nil {
			t.Fail()
		}

		switch scena.Name {
		case 0:
			transSrv.EXPECT().GetTransactions(
				gomock.Any(),
				gomock.Any(),
			).Return(
				scena.ExpectSrvGetTransactions.Model,
				scena.ExpectSrvGetTransactions.Err,
			)

			responseDataOut, err := DecodeTransactionsSuccess(transHandler, req, transHandler.GetTransaction)
			if err != nil {
				t.Fail()
			}
			var outputObj models.GetTransactionRes
			var resp []byte

			if resp, err = json.Marshal(responseDataOut.ServiceResponseBody); err != nil {
				t.Fail()
			}
			if err := json.Unmarshal(resp, &outputObj); err != nil {
				t.Fail()
			}

			expectObj := scena.ExpectGetTransactions.Model

			buff1, errJSON1 := json.Marshal(outputObj)
			buff2, errJSON2 := json.Marshal(expectObj)
			if errJSON1 != nil || errJSON2 != nil {
				t.Fail()
			}

			assert.JSONEq(t, string(buff2), string(buff1), "two data aren't same")
		case 1:
			transSrv.EXPECT().GetTransactions(
				gomock.Any(),
				gomock.Any(),
			).Return(
				scena.ExpectSrvGetTransactions.Model,
				scena.ExpectSrvGetTransactions.Err,
			)

			responseDataOut, err := DecodeTransactionsFail(transHandler, req, transHandler.GetTransaction)
			if err != nil {
				t.Fail()
			}
			var outputObj models.KbankResponseHeader
			var resp []byte

			if resp, err = json.Marshal(responseDataOut.KbankResponseHeader); err != nil {
				t.Fail()
			}
			if err := json.Unmarshal(resp, &outputObj); err != nil {
				t.Fail()
			}

			expectErr := scena.ExpectGetTransactions.ErrCode
			if expectErr != outputObj.Errors.ErrorCode {
				t.Fail()
			}
		}
	}
}

type (
	ExpectGetTransactions struct {
		Model     models.GetTransactionRes
		ErrCode   string
		ErrorDesc string
	}
	ExpectSrvGetTransactions struct {
		Model []models.Transactions
		Err   error
	}
	TestGetTransactionsCase struct {
		Name        int
		Description string
		Input       models.RequestInfo
		ExpectSrvGetTransactions
		ExpectGetTransactions
	}
)

func GetTransactionsCases() []TestGetTransactionsCase {
	var output []TestGetTransactionsCase

	//===================== 1. Happy case ===================== //
	case1 := TestGetTransactionsCase{
		Name:        0,
		Description: "Successfully calling to get transaction",
		Input: models.RequestInfo{
			Header: models.KbankRequestHeader{
				RequestAppID: "781",
				RequestDate:  time.Now().Format(time.RFC3339),
				RequestUID:   "request_id",
			},
			Body: models.GetTransactionReq{
				TransactionID: 1,
			},
		},
		ExpectSrvGetTransactions: ExpectSrvGetTransactions{
			Model: []models.Transactions{},
			Err:   nil,
		},
		ExpectGetTransactions: ExpectGetTransactions{
			Model: models.GetTransactionRes{
				ListTransactions: []models.Transactions{},
			},
			ErrCode:   "00",
			ErrorDesc: "Success",
		},
	}
	case1.ExpectSrvGetTransactions.Model = append(case1.ExpectSrvGetTransactions.Model, models.Transactions{
		TransactionID:   1,
		TransactionName: "transaction_1",
	})
	case1.ExpectGetTransactions.Model.ListTransactions = append(case1.ExpectGetTransactions.Model.ListTransactions, models.Transactions{
		TransactionID:   1,
		TransactionName: "transaction_1",
	})
	output = append(output, case1)

	//===================== 2. Database Server Error Case ===================== //
	case2 := TestGetTransactionsCase{
		Name:        1,
		Description: "Successfully calling to get transaction",
		Input: models.RequestInfo{
			Header: models.KbankRequestHeader{
				RequestAppID: "781",
				RequestDate:  time.Now().Format(time.RFC3339),
				RequestUID:   "request_id",
			},
			Body: models.GetTransactionReq{
				TransactionID: 2,
			},
		},
		ExpectSrvGetTransactions: ExpectSrvGetTransactions{
			Model: []models.Transactions{},
			Err: kerrors.WithPrefix("Transaction_handler").Wrap(
				errors.New("Database Server Error"),
				kerrors.DatabaseServerError,
				nil,
			),
		},
		ExpectGetTransactions: ExpectGetTransactions{
			Model: models.GetTransactionRes{
				ListTransactions: []models.Transactions{},
			},
			ErrCode:   "2000",
			ErrorDesc: "Database Server Error",
		},
	}
	output = append(output, case2)

	return output
}

func DecodeTransactionsFail(handler *TransactionHandler, req *http.Request, handerFunc http.HandlerFunc) (response.FailureResponseFormat, error) {
	rr := httptest.NewRecorder()
	http.HandlerFunc(handerFunc).ServeHTTP(rr, req)

	var respFail response.FailureResponseFormat
	if err := json.NewDecoder(rr.Body).Decode(&respFail); err != nil {
		return respFail, err
	}
	return respFail, nil
}

func DecodeTransactionsSuccess(handler *TransactionHandler, req *http.Request, handerFunc http.HandlerFunc) (response.SuccessResponseFormat, error) {
	rr := httptest.NewRecorder()
	http.HandlerFunc(handerFunc).ServeHTTP(rr, req)

	var respSuccess response.SuccessResponseFormat
	if err := json.NewDecoder(rr.Body).Decode(&respSuccess); err != nil {
		return respSuccess, err
	}
	return respSuccess, nil
}

func TestInsertTransactions(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	transSrv := transSrvMock.NewMockTransactionsService(ctl)
	conf := &config.Configs{}

	transHandler := NewTransactionHandler(conf, transSrv)

	scenarios := InsertTransactionsCases()

	for _, scena := range scenarios {
		inputJSON, errInput := json.Marshal(scena.Input)
		if errInput != nil {
			t.Fail()
		}

		req, errRequest := http.NewRequest("POST", "/kph/api/set", bytes.NewBuffer(inputJSON))
		if errRequest != nil {
			t.Fail()
		}

		switch scena.Name {
		case 0:
			transSrv.EXPECT().InsertTransaction(
				gomock.Any(),
				gomock.Any(),
			).Return(
				scena.ExpectSrvInsertTransactions.Model,
				scena.ExpectSrvInsertTransactions.Err,
			)

			responseDataOut, err := DecodeTransactionsSuccess(transHandler, req, transHandler.InsertTransaction)
			if err != nil {
				t.Fail()
			}
			var outputObj models.InsertTransactionSrvRes
			var resp []byte

			if resp, err = json.Marshal(responseDataOut.ServiceResponseBody); err != nil {
				t.Fail()
			}
			if err := json.Unmarshal(resp, &outputObj); err != nil {
				t.Fail()
			}

			expectObj := scena.ExpectInsertTransactions.Model

			buff1, errJSON1 := json.Marshal(outputObj)
			buff2, errJSON2 := json.Marshal(expectObj)
			if errJSON1 != nil || errJSON2 != nil {
				t.Fail()
			}

			assert.JSONEq(t, string(buff2), string(buff1), "two data aren't same")
		case 1:
			transSrv.EXPECT().InsertTransaction(
				gomock.Any(),
				gomock.Any(),
			).Return(
				scena.ExpectSrvInsertTransactions.Model,
				scena.ExpectSrvInsertTransactions.Err,
			)

			responseDataOut, err := DecodeTransactionsFail(transHandler, req, transHandler.InsertTransaction)
			if err != nil {
				t.Fail()
			}
			var outputObj models.KbankResponseHeader
			var resp []byte

			if resp, err = json.Marshal(responseDataOut.KbankResponseHeader); err != nil {
				t.Fail()
			}
			if err := json.Unmarshal(resp, &outputObj); err != nil {
				t.Fail()
			}

			expectErr := scena.ExpectInsertTransactions.ErrCode
			if expectErr != outputObj.Errors.ErrorCode {
				t.Fail()
			}
		}
	}
}

type (
	ExpectInsertTransactions struct {
		Model     models.InsertTransactionSrvRes
		ErrCode   string
		ErrorDesc string
	}
	ExpectSrvInsertTransactions struct {
		Model models.InsertTransactionSrvRes
		Err   error
	}
	TestInsertTransactionsCase struct {
		Name        int
		Description string
		Input       models.RequestInfo
		ExpectSrvInsertTransactions
		ExpectInsertTransactions
	}
)

func InsertTransactionsCases() []TestInsertTransactionsCase {
	var output []TestInsertTransactionsCase

	//===================== 1. Happy case ===================== //
	case1 := TestInsertTransactionsCase{
		Name:        0,
		Description: "Successfully calling to insert transaction",
		Input: models.RequestInfo{
			Header: models.KbankRequestHeader{
				RequestAppID: "781",
				RequestDate:  time.Now().Format(time.RFC3339),
				RequestUID:   "request_id",
			},
			Body: models.InsertTransactionReq{
				TransactionName: "transaction_name_1",
			},
		},
		ExpectSrvInsertTransactions: ExpectSrvInsertTransactions{
			Model: models.InsertTransactionSrvRes{
				TransactionID:   1,
				TransactionName: "transaction_name_1",
			},
			Err: nil,
		},
		ExpectInsertTransactions: ExpectInsertTransactions{
			Model: models.InsertTransactionSrvRes{
				TransactionID:   1,
				TransactionName: "transaction_name_1",
			},
			ErrCode:   "00",
			ErrorDesc: "Success",
		},
	}

	output = append(output, case1)

	//===================== 2. Database Server Error Case ===================== //
	case2 := TestInsertTransactionsCase{
		Name:        1,
		Description: "Failed calling to insert transaction",
		Input: models.RequestInfo{
			Header: models.KbankRequestHeader{
				RequestAppID: "781",
				RequestDate:  time.Now().Format(time.RFC3339),
				RequestUID:   "request_id",
			},
			Body: models.InsertTransactionReq{
				TransactionName: "transaction_name_1",
			},
		},
		ExpectSrvInsertTransactions: ExpectSrvInsertTransactions{
			Model: models.InsertTransactionSrvRes{},
			Err: kerrors.WithPrefix("Transaction_handler").Wrap(
				errors.New("Database Server Error"),
				kerrors.DatabaseServerError,
				nil,
			),
		},
		ExpectInsertTransactions: ExpectInsertTransactions{
			Model:     models.InsertTransactionSrvRes{},
			ErrCode:   "2000",
			ErrorDesc: "Database Server Error",
		},
	}
	output = append(output, case2)

	return output
}
