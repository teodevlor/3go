package registry

import (
	pgdb "go-structure/internal/orm/db/postgres"
	account_repo "go-structure/internal/repository"
	user_profile_repo "go-structure/internal/repository/app_user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sarulabs/di"
)

const (
	UserProfileRepoDIName = "user_profile_repo_di"
	AccountRepoDIName     = "account_repo_di"
)

func buildRepositories() error {
	userProfileDef := di.Def{
		Name:  UserProfileRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			q := pgdb.New(pool)
			return user_profile_repo.NewUserProfileRepository(q), nil
		},
	}

	accountDef := di.Def{
		Name:  AccountRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			q := pgdb.New(pool)
			return account_repo.NewAccountRepository(q), nil
		},
	}
	return builder.Add(userProfileDef, accountDef)
}
