package pub

import (
	"github.com/mondegor/go-components/mrauth/repository"
	"github.com/mondegor/go-storage/mrsql"

	"github.com/mondegor/print-shop-back/internal/auth/module"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

func createUserPostgres(opts auth.Options) *repository.UserPostgres {
	return repository.NewUserPostgres(
		opts.DBConnManager,
		opts.StorageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".users",
			PrimaryKey: "user_id",
		},
	)
}

func createCheckUserPostgres(opts auth.Options) *repository.CheckUserPostgres {
	return repository.NewCheckUserPostgres(
		opts.DBConnManager,
		opts.StorageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".users",
			PrimaryKey: "user_id",
		},
	)
}

func createUserRealmPostgres(opts auth.Options) *repository.UserRealmPostgres {
	return repository.NewUserRealmPostgres(
		opts.DBConnManager,
		opts.StorageErrorWrapper,
		module.DBSchema+".users_realms",
	)
}

func createAuth2faPostgres(opts auth.Options) *repository.Auth2faPostgres {
	return repository.NewAuth2faPostgres(
		opts.DBConnManager,
		opts.StorageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".users_auth_2fa",
			PrimaryKey: "user_id",
		},
	)
}

func createUserActivityStatPostgres(opts auth.Options) *repository.UserActivityStatPostgres {
	return repository.NewUserActivityStatPostgres(
		opts.DBConnManager,
		opts.StorageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".users_activity_stat",
			PrimaryKey: "user_id",
		},
	)
}

// func createUserActivityLogPostgres(opts auth.Options) *repository.UserActivityLogPostgres {
// 	return repository.NewUserActivityLogPostgres(
// 		opts.DBConnManager,
//      opts.StorageErrorWrapper,
// 		module.DBSchema+".secure_operations_log",
// 	)
// }

// 111111111111111111
func createAuthTokenPostgres(opts auth.Options) *repository.AuthTokenPostgres {
	return repository.NewAuthTokenPostgres(
		opts.DBConnManager,
		opts.StorageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".auth_tokens",
			PrimaryKey: "refresh_token",
		},
	)
}

func createSecureOperationPostgres(opts auth.Options) *repository.SecureOperationPostgres {
	return repository.NewSecureOperationPostgres(
		opts.DBConnManager,
		opts.StorageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".secure_operations",
			PrimaryKey: "operation_token",
		},
	)
}

// func createSecureOperationLogPostgres(opts auth.Options) *repository.SecureOperationLogPostgres {
// 	return repository.NewSecureOperationLogPostgres(
// 		opts.DBConnManager,
//      opts.StorageErrorWrapper,
// 		module.DBSchema+".secure_operations_log",
// 	)
// }
