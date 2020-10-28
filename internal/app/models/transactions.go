package models

type (
	// Transactions structure
	Transactions struct {
		TransactionID   int    `json:"transaction_id" bson:"transaction_id"`
		TransactionName string `json:"transaction_name" bson:"transaction_name"`
	}
	// GetTransactionRepoReq structure
	GetTransactionRepoReq struct {
		TransactionID int `json:"transaction_id"`
	}
	// GetTransactionSrvReq structure
	GetTransactionSrvReq struct {
		TransactionID int `json:"transaction_id"`
	}
	// GetTransactionReq structure
	GetTransactionReq struct {
		TransactionID int `json:"transaction_id"`
	}
	// GetTransactionRes structure
	GetTransactionRes struct {
		ListTransactions []Transactions `json:"list_transactions"`
	}
)

type (
	// InsertTransactionRepoReq structure
	InsertTransactionRepoReq struct {
		TransactionName string `json:"transaction_name"`
	}
	// InsertTransactionSrvReq structure
	InsertTransactionSrvReq struct {
		TransactionName string `json:"transaction_name"`
	}
	// InsertTransactionReq structure
	InsertTransactionReq struct {
		TransactionName string `json:"transaction_name"`
	}
	// InsertTransactionSrvRes structure
	InsertTransactionSrvRes struct {
		TransactionID   int64  `json:"transaction_id"`
		TransactionName string `json:"transaction_name"`
	}
)
