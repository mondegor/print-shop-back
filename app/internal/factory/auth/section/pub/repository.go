package pub

import (
	"github.com/mondegor/go-components/mrauth/repository"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/print-shop-back/internal/auth/module"
)

func initUserPostgres(
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
) *repository.UserPostgres {
	return repository.NewUserPostgres(
		dbConnManager,
		storageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".users",
			PrimaryKey: "user_id",
		},
	)
}

func initCheckUserPostgres(
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
) *repository.CheckUserPostgres {
	return repository.NewCheckUserPostgres(
		dbConnManager,
		storageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".users",
			PrimaryKey: "user_id",
		},
	)
}

func initUserRealmPostgres(
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
) *repository.UserRealmPostgres {
	return repository.NewUserRealmPostgres(
		dbConnManager,
		storageErrorWrapper,
		module.DBSchema+".users_realms",
	)
}

func initAuth2faPostgres(
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
) *repository.Auth2faPostgres {
	return repository.NewAuth2faPostgres(
		dbConnManager,
		storageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".users_auth_2fa",
			PrimaryKey: "user_id",
		},
	)
}

func initUserActivityStatPostgres(
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
) *repository.UserActivityStatPostgres {
	return repository.NewUserActivityStatPostgres(
		dbConnManager,
		storageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".users_activity_stat",
			PrimaryKey: "user_id",
		},
	)
}

// func initUserActivityLogPostgres(
//	 storageErrorWrapper mrerr.ErrorWrapper,
//	 dbConnManager mrstorage.DBConnManager,
// ) *repository.UserActivityLogPostgres {
//	 return repository.NewUserActivityLogPostgres(
//		 dbConnManager,
//		 storageErrorWrapper,
//		 module.DBSchema+".secure_operations_log",
//	 )
// }

func initAuthTokenPostgres(
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
) *repository.AuthTokenPostgres {
	return repository.NewAuthTokenPostgres(
		dbConnManager,
		storageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".auth_tokens",
			PrimaryKey: "refresh_token",
		},
	)
}

func initSecureOperationPostgres(
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
) *repository.SecureOperationPostgres {
	return repository.NewSecureOperationPostgres(
		dbConnManager,
		storageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       module.DBSchema + ".secure_operations",
			PrimaryKey: "operation_token",
		},
	)
}

// func initSecureOperationLogPostgres(
//	 storageErrorWrapper mrerr.ErrorWrapper,
//	 dbConnManager mrstorage.DBConnManager,
// ) *repository.SecureOperationLogPostgres {
//	 return repository.NewSecureOperationLogPostgres(
//		 dbConnManager,
//		 storageErrorWrapper,
//		 module.DBSchema+".secure_operations_log",
//	 )
// }
