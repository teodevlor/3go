package registry

import (
	"context"
	"fmt"
	"sync"

	"go-structure/config"
	v1 "go-structure/internal/api/v1"
	otpcontroller "go-structure/internal/controller"
	controller "go-structure/internal/controller/app_user"
	"go-structure/internal/helper/database"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sarulabs/di"
)

const (
	// config
	ConfigDIName = "config_di"
	ConfigDir    = "./config"
	ConfigName   = "dev.yml"

	// api
	ApiDIName = "api_di"
)

var (
	buildOnce sync.Once
	builder   *di.Builder
	container di.Container
)

func BuildDependencyInjectContainer() {
	buildOnce.Do(func() {
		builder, _ = di.NewBuilder()
		if err := buildConfigs(); err != nil {
			panic(err)
		}
		if err := buildDatabase(); err != nil {
			panic(err)
		}
		if err := buildAdapters(); err != nil {
			panic(err)
		}
		if err := buildRepositories(); err != nil {
			panic(err)
		}
		if err := buildUsecases(); err != nil {
			panic(err)
		}
		if err := buildControllers(); err != nil {
			panic(err)
		}
		if err := buildApis(); err != nil {
			panic(err)
		}
		container = builder.Build()
	})
}

func GetDependency(dependencyName string) interface{} {
	return container.Get(dependencyName)
}

func CleanDependency() error {
	return container.Clean()
}

func buildConfigs() error {
	defs := []di.Def{}

	configDef := di.Def{
		Name:  ConfigDIName,
		Scope: di.App,
		Build: func(di di.Container) (interface{}, error) {
			if err := config.Load(ConfigDir, ConfigName); err != nil {
				return nil, err
			}
			return config.Cfg, nil
		},
		Close: func(obj interface{}) error {
			return nil
		},
	}
	defs = append(defs, configDef)
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}

func buildDatabase() error {
	// Pool Postgres for SQLC (ORM - SQLC)
	pgPoolDef := di.Def{
		Name:  PostgresPoolDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(ConfigDIName).(*config.Config)
			dbCfg := cfg.DatabasePostgre
			dsn := fmt.Sprintf(
				"postgres://%s:%s@%s:%d/%s?sslmode=disable",
				dbCfg.User,
				dbCfg.Password,
				dbCfg.Host,
				dbCfg.Port,
				dbCfg.DBName,
			)
			pool, err := pgxpool.New(context.Background(), dsn)
			if err != nil {
				return nil, err
			}
			if err := pool.Ping(context.Background()); err != nil {
				pool.Close()
				return nil, err
			}
			return pool, nil
		},
		Close: func(obj interface{}) error {
			if pool, ok := obj.(*pgxpool.Pool); ok && pool != nil {
				pool.Close()
			}
			return nil
		},
	}

	// Transaction Manager for managing transactions
	txManagerDef := di.Def{
		Name:  TransactionManagerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return database.NewTransactionManager(pool), nil
		},
	}

	return builder.Add(pgPoolDef, txManagerDef)
}

func buildApis() error {
	defs := []di.Def{}

	apiDef := di.Def{
		Name:  ApiDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			router := gin.Default()
			// router := gin.New()
			userProfileController := ctn.Get(UserProfileControllerDIName).(controller.UserProfileController)
			otpController := ctn.Get(OTPControllerDIName).(otpcontroller.OTPController)
			v1.NewApiV1(router, userProfileController, otpController)
			return router, nil
		},
	}

	defs = append(defs, apiDef)
	err := builder.Add(defs...)
	if err != nil {
		return err
	}
	return nil
}
