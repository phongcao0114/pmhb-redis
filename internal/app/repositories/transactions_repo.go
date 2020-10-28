package repositories

//
//import (
//	"context"
//	"database/sql"
//	"fmt"
//	"pmhb-redis/internal/app/config"
//	"pmhb-redis/internal/app/models"
//	"pmhb-redis/internal/kerrors"
//	"pmhb-redis/internal/pkg/db"
//	"pmhb-redis/internal/pkg/klog"
//)
//
//const (
//	// TransactionsRepositoryPrefix prefix repo
//	TransactionsRepositoryPrefix = "Transactions_repository"
//)
//
//type (
//	// TransactionsRepo defines mssqldb.Client connection for each client
//	TransactionsRepo struct {
//		s       *db.DB
//		c       *config.Configs
//		errRepo kerrors.KError
//		logger  klog.Logger
//	}
//
//	//TransactionsRepository groups all function integrate with transaction collection in mssqldbdb
//	TransactionsRepository interface {
//		GetTransaction(ctx context.Context, req models.GetTransactionRepoReq) ([]models.Transactions, error)
//		InsertTransaction(ctx context.Context, req models.InsertTransactionRepoReq) (int64, error)
//	}
//)
//
//// NewTransactionsRepo opens the connection to DB from repositories package
//func NewTransactionsRepo(configs *config.Configs, s *db.DB) *TransactionsRepo {
//	// Return model
//	return &TransactionsRepo{
//		s:       s,
//		c:       configs,
//		errRepo: kerrors.WithPrefix(TransactionsRepositoryPrefix),
//		logger:  klog.WithPrefix(TransactionsRepositoryPrefix),
//	}
//}
//
//// GetTransaction function
//func (tr *TransactionsRepo) GetTransaction(ctx context.Context, req models.GetTransactionRepoReq) ([]models.Transactions, error) {
//	var dbList []models.Transactions
//	ctx = context.Background()
//
//	err := tr.s.PingContext(ctx)
//	if err != nil {
//		return dbList, tr.errRepo.Wrap(err, kerrors.DatabaseServerError, nil)
//	}
//
//	db := fmt.Sprintf("%s.%s.%s", tr.c.MSSQL.DatabaseName, "dbo", tr.c.MSSQL.Tables.Transactions)
//	tsql := fmt.Sprintf("SELECT transaction_id , transaction_name FROM %s WHERE transaction_id=%d;", db, req.TransactionID)
//
//	// Execute query
//	rows, err := tr.s.QueryContext(ctx, tsql)
//	if err != nil {
//		return dbList, tr.errRepo.Wrap(err, kerrors.CannotGetDataFromDB, nil)
//	}
//
//	defer rows.Close()
//	// Iterate through the result set.
//	for rows.Next() {
//		// Get values from row.
//		var trans models.Transactions
//		err := rows.Scan(&trans.TransactionID, &trans.TransactionName)
//		if err != nil {
//			return dbList, tr.errRepo.Wrap(err, kerrors.CannotGetDataFromDB, nil)
//		}
//		dbList = append(dbList, trans)
//	}
//
//	return dbList, nil
//}
//
//// InsertTransaction function
//func (tr *TransactionsRepo) InsertTransaction(ctx context.Context, req models.InsertTransactionRepoReq) (int64, error) {
//	var newID int64
//	ctx = context.Background()
//
//	err := tr.s.PingContext(ctx)
//	if err != nil {
//		return newID, tr.errRepo.Wrap(err, kerrors.DatabaseServerError, nil)
//	}
//
//	db := fmt.Sprintf("%s.%s.%s", tr.c.MSSQL.DatabaseName, "dbo", tr.c.MSSQL.Tables.Transactions)
//	tsql := fmt.Sprintf("INSERT INTO %s (transaction_name) VALUES (@transaction_name); select convert(bigint, SCOPE_IDENTITY());", db)
//
//	stmt, err := tr.s.Prepare(tsql)
//	if err != nil {
//		return newID, tr.errRepo.Wrap(err, kerrors.DatabaseServerError, nil)
//	}
//	defer stmt.Close()
//
//	row := stmt.QueryRowContext(
//		ctx,
//		sql.Named("transaction_name", req.TransactionName),
//	)
//
//	err = row.Scan(&newID)
//	if err != nil {
//		return newID, tr.errRepo.Wrap(err, kerrors.DatabaseScanErr, nil)
//	}
//
//	return newID, nil
//}
//
//// OracleGetTransaction function
//func (tr *TransactionsRepo) OracleGetTransaction(ctx context.Context, req models.GetTransactionRepoReq) ([]models.Transactions, error) {
//	var dbList []models.Transactions
//	ctx = context.Background()
//
//	err := tr.s.PingContext(ctx)
//	if err != nil {
//		return dbList, tr.errRepo.Wrap(err, kerrors.DatabaseServerError, nil)
//	}
//
//	table := tr.c.Oracle.Tables.Transactions
//
//	query := fmt.Sprintf("SELECT TRANSACTION_ID, TRANSACTION_NAME FROM %s WHERE TRANSACTION_ID=:id", table)
//
//	//Execute query
//	rows, err := tr.s.QueryContext(
//		ctx,
//		query,
//		sql.Named("id", req.TransactionID),
//	)
//	if err != nil {
//		return dbList, tr.errRepo.Wrap(err, kerrors.CannotGetDataFromDB, nil)
//	}
//
//	defer rows.Close()
//	// Iterate through the result set.
//	for rows.Next() {
//		// Get values from row.
//		var trans models.Transactions
//		err := rows.Scan(&trans.TransactionID, &trans.TransactionName)
//		if err != nil {
//			return dbList, tr.errRepo.Wrap(err, kerrors.CannotGetDataFromDB, nil)
//		}
//		dbList = append(dbList, trans)
//	}
//
//	return dbList, nil
//}
//
//// OracleInsertTransaction function
//func (tr *TransactionsRepo) OracleInsertTransaction(ctx context.Context, req models.InsertTransactionRepoReq) (int64, error) {
//	var newID int64
//	ctx = context.Background()
//
//	err := tr.s.PingContext(ctx)
//	if err != nil {
//		return newID, tr.errRepo.Wrap(err, kerrors.DatabaseServerError, nil)
//	}
//
//	// TODO Create Unique ID for inserting (not to use auto increment)
//	table := tr.c.Oracle.Tables.Transactions
//	tsql := fmt.Sprintf("INSERT INTO %s (TRANSACTION_NAME) VALUES (:TRANSACTION_NAME)", table)
//
//	stmt, err := tr.s.Prepare(tsql)
//	if err != nil {
//		return newID, tr.errRepo.Wrap(err, kerrors.DatabaseServerError, nil)
//	}
//	defer stmt.Close()
//
//	stmt.QueryRowContext(
//		ctx,
//		sql.Named("TRANSACTION_NAME", req.TransactionName),
//	)
//
//	return newID, nil
//}
//
//// // UpdateOneTransaction function
//// func (tr *TransactionsRepo) UpdateOneTransaction(ctx context.Context, search, update bson.M) error {
//// 	count, err := tr.s.Database(tr.c.MSSQL.DatabaseName).Collection(tr.c.MSSQL.Tables.Transactions).UpdateOne(
//// 		context.Background(),
//// 		search,
//// 		update)
//// 	if err != nil {
//// 		return tr.errRepo.Wrap(err, kerrors.DatabaseServerError, nil)
//// 	}
//// 	if count.ModifiedCount == 0 {
//// 		return tr.errRepo.Wrap(errors.new("No data has been found"), kerrors.NotFoundItemInQuery, nil)
//// 	}
//// 	return nil
//// }
