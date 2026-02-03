package registry

import (
	"errors"
	"fmt"
	"sync"

	"go-structure/config"
	"go-structure/internal/helper/database"
)

const (
	DriverPostgres = "postgres"
	DriverMySQL    = "mysql"

	// DI name for each driver. Decide to use postgres or mysql in di_repository / di_usecase.
	DatabasePostgresDIName = "database_postgres_di"
	DatabaseMySQLDIName    = "database_mysql_di"

	// PostgresPoolDIName is the DI name for *pgxpool.Pool (used by sqlc).
	PostgresPoolDIName = "postgres_pool_di"
)

var (
	ErrDriverNotRegistered = errors.New("database: driver chưa được đăng ký")
	ErrDriverEmpty         = errors.New("database: driver không được để trống")
)

type DBBuilder func(cfg *config.Config) (database.DBHelper, error)

var (
	dbRegistry   = make(map[string]DBBuilder)
	dbRegistryMu sync.RWMutex
)

func init() {
	RegisterDB(DriverPostgres, func(cfg *config.Config) (database.DBHelper, error) {
		return database.NewPostgresDB(cfg.DatabasePostgre)
	})
	RegisterDB(DriverMySQL, func(cfg *config.Config) (database.DBHelper, error) {
		return database.NewMySQLDB(cfg.DatabaseMySQL)
	})
}

func RegisterDB(driver string, builder DBBuilder) {
	dbRegistryMu.Lock()
	defer dbRegistryMu.Unlock()
	if driver == "" {
		panic("database: driver name cannot be empty")
	}
	dbRegistry[driver] = builder
}

func RegisteredDBDrivers() []string {
	dbRegistryMu.RLock()
	defer dbRegistryMu.RUnlock()
	names := make([]string, 0, len(dbRegistry))
	for name := range dbRegistry {
		names = append(names, name)
	}
	return names
}

func NewDB(driver string, cfg *config.Config) (database.DBHelper, error) {
	if driver == "" {
		return nil, ErrDriverEmpty
	}
	dbRegistryMu.RLock()
	builder, ok := dbRegistry[driver]
	dbRegistryMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrDriverNotRegistered, driver)
	}
	return builder(cfg)
}
