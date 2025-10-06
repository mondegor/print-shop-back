package pub

import (
	"github.com/mondegor/go-components/mrnotifier"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/module"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
	"github.com/mondegor/print-shop-back/internal/initing"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

// InitHttpModule - создаются все компоненты модуля и возвращаются к нему контролеры.
func InitHttpModule(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	storageErrorWrapper mrerr.ErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	locker mrlock.Locker,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
	responseFileSender mrserver.FileResponseSender,
	notifierAPI mrnotifier.NoticeProducer,
	withDebugInfo bool,
	userRealms []auth.UserRealm,
	operationConfirm auth.OperationConfirm,
	jwtConfig auth.JWTConfig,
) initing.HttpModule {
	storageUser := initUserPostgres(storageErrorWrapper, dbConnManager)
	storageCheckUser := initCheckUserPostgres(storageErrorWrapper, dbConnManager)
	storageUserRealm := initUserRealmPostgres(storageErrorWrapper, dbConnManager)
	storageAuth2fa := initAuth2faPostgres(storageErrorWrapper, dbConnManager)
	storageUserActivityStat := initUserActivityStatPostgres(storageErrorWrapper, dbConnManager)
	// storageUserActivityLog := initUserActivityLogPostgres(storageErrorWrapper, dbConnManager)
	storageAuthToken := initAuthTokenPostgres(storageErrorWrapper, dbConnManager)
	storageSecureOperation := initSecureOperationPostgres(storageErrorWrapper, dbConnManager)
	// storageSecureOperationLog := initSecureOperationLogPostgres(storageErrorWrapper, dbConnManager)
	useCaseConfirmOperation := initConfirmOperationUseCase(
		useCaseErrorWrapper,
		dbConnManager,
		storageSecureOperation,
		notifierAPI,
		operationConfirm,
	)

	return initing.HttpModule{
		Name:       module.Name,
		Permission: module.Permission,
		Controllers: []initing.HttpController{
			{
				Create: func() (mrserver.HttpController, error) {
					return initUnitAuthController(
						logger,
						eventEmitter,
						useCaseErrorWrapper,
						dbConnManager,
						storageUser,
						storageCheckUser,
						storageUserRealm,
						storageAuth2fa,
						storageUserActivityStat,
						storageAuthToken,
						storageSecureOperation,
						useCaseConfirmOperation,
						locker,
						requestParser,
						responseSender,
						notifierAPI,
						withDebugInfo,
						userRealms,
						jwtConfig,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initCheckController(
						useCaseErrorWrapper,
						storageCheckUser,
						storageUserRealm,
						requestParser,
						responseSender,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initOperationController(
						useCaseErrorWrapper,
						dbConnManager,
						storageSecureOperation,
						useCaseConfirmOperation,
						requestParser,
						responseSender,
						notifierAPI,
						withDebugInfo,
						operationConfirm,
					)
				},
			},
			{
				Create: func() (mrserver.HttpController, error) {
					return initSecurityController(
						logger,
						useCaseErrorWrapper,
						dbConnManager,
						storageUser,
						storageCheckUser,
						storageUserRealm,
						storageAuth2fa,
						storageSecureOperation,
						requestParser,
						responseFileSender,
						notifierAPI,
						withDebugInfo,
					)
				},
			},
		},
	}
}
