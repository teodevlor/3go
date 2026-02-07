package registry

const (
	// PostgresPoolDIName is the DI name for *pgxpool.Pool (used by sqlc).
	PostgresPoolDIName = "postgres_pool_di"

	// TransactionManagerDIName is the DI name for TransactionManager (pgx).
	TransactionManagerDIName = "transaction_manager_di"
)
