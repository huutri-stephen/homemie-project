package posrgres

import (
	"cloud.google.com/go/cloudsqlconn"
	"context"
	"database/sql"
	"fmt"
	"homemie/config"
	"homemie/pkg/custom_error"
	"homemie/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"net"
)

var (
	db *sql.DB
)

const TRANSACTION_KEY = "TRANSACTION"

func GetDbConnection(ctx context.Context, config *config.Config) (*sql.DB, error) {
	if config.PG.Mode == "NORMAL" {
		logger.Infof(ctx, "Using normal connection")

		var err error
		db, err = sql.Open("postgres", config.PG.URL)
		if err != nil {
			logger.Errorf(ctx, "Failed to create db connection, error %s", err)
			return nil, err
		}
		return db, err
	}
	if config.PG.Mode == "GCP" {
		logger.Infof(ctx, "Using GCP connection")
		return getDdbConnectionFromGCP(config)
	}

	return nil, fmt.Errorf("unknown mode %s", config.PG.Mode)
}

func getDdbConnectionFromGCP(config *config.Config) (*sql.DB, error) {

	dbUser := config.PG.DbUser
	if dbUser == "" {
		log.Fatalf("Fatal Error in connect_connector.go: DB_USER environment variable not set.\n")
	}
	dbPwd := config.PG.DbPassword
	if dbPwd == "" {
		log.Fatalf("Fatal Error in connect_connector.go: DB_PASS environment variable not set.\n")
	}
	dbName := config.PG.DbName
	if dbName == "" {
		log.Fatalf("Fatal Error in connect_connector.go: DB_NAME environment variable not set.\n")
	}
	instanceConnectionName := config.PG.InstanceConnectionName
	if instanceConnectionName == "" {
		log.Fatalf("Fatal Error in connect_connector.go: INSTANCE_CONNECTION_NAME environment variable not set.\n")
	}
	// optional
	usePrivate := config.PG.UsePrivateIP

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
	logger.Infof(context.Background(), "dns: %s", dsn)

	parseConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	var opts []cloudsqlconn.Option
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	// Use the Cloud SQL connector to handle connecting to the instance.
	// This approach does *NOT* require the Cloud SQL proxy.
	parseConfig.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(parseConfig)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil
}

func SetupBoilerConnection(ctx context.Context, config *config.Config) {
	var err error
	db, err = GetDbConnection(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	boil.DebugMode = true
	//Set db to global variable
	boil.SetDB(db)
}

type Transaction struct {
	context.Context
	Tx *sql.Tx
}

func (t *Transaction) Commit(ctx context.Context) error {
	err := t.Tx.Commit()
	if err != nil {
		logger.Errorf(ctx, "Fail to commit transaction, error %s", err)
		return custom_error.NewInternalServerError(custom_error.INTERNAL_SERVER_ERROR, fmt.Sprintf("Fail to commit transaction,error %s", err))
	}
	return nil
}

func (t *Transaction) Rollback(ctx context.Context) {
	rollbackErr := t.Tx.Rollback()
	if rollbackErr != nil {
		logger.Errorf(ctx, "Transaction failed and then failed to rollback")
	}
}

func CreateTransaction(ctx context.Context) (*Transaction, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Errorf(ctx, "Fail to create transaction , error %s", err)
		return nil, custom_error.NewInternalServerError(custom_error.INTERNAL_SERVER_ERROR, fmt.Sprintf("Fail to create transaction: %s", err))
	}
	return &Transaction{
		Context: context.WithValue(ctx, TRANSACTION_KEY, tx),
		Tx:      tx,
	}, nil
}

func GetTransactionFromContext(ctx context.Context) boil.ContextExecutor {
	tx, ok := ctx.Value(TRANSACTION_KEY).(*sql.Tx)
	if !ok {
		return db
	}
	return tx
}
