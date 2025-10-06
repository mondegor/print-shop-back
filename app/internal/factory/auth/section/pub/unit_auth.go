package pub

import (
	"time"

	"github.com/mondegor/go-components/mrauth/bag/contactaddress"
	"github.com/mondegor/go-components/mrauth/component/secureoperation/action"
	"github.com/mondegor/go-components/mrauth/repository"
	"github.com/mondegor/go-components/mrauth/service"
	usecaseauth "github.com/mondegor/go-components/mrauth/usecase/auth"
	"github.com/mondegor/go-components/mrauth/usecase/check"
	"github.com/mondegor/go-components/mrauth/usecase/operation"
	"github.com/mondegor/go-components/mrauth/usecase/session"
	"github.com/mondegor/go-components/mrauth/usecase/session/handler"
	"github.com/mondegor/go-components/mrnotifier"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlock"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/bag"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
	"github.com/mondegor/print-shop-back/internal/factory/auth/mapping"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initUnitAuthController(
	logger mrlog.Logger,
	eventEmitter mrevent.Emitter,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	storageUser *repository.UserPostgres,
	storageCheckUser *repository.CheckUserPostgres,
	storageUserRealm *repository.UserRealmPostgres,
	storageAuth2fa *repository.Auth2faPostgres,
	storageUserActivityStat *repository.UserActivityStatPostgres,
	storageAuthToken *repository.AuthTokenPostgres,
	storageSecureOperation *repository.SecureOperationPostgres,
	useCaseConfirmOperation *operation.ConfirmOperation,
	locker mrlock.Locker,
	requestParser *validate.Parser,
	responseSender mrserver.ResponseSender,
	notifierAPI mrnotifier.NoticeProducer,
	withDebugInfo bool,
	userRealms []auth.UserRealm,
	jwtConfig auth.JWTConfig,
) (mrserver.HttpController, error) {
	contactAddressParser := contactaddress.NewParser()

	checkUserUseCase := check.NewAuthHelper(
		storageCheckUser,
		storageUserRealm,
		contactAddressParser,
		useCaseErrorWrapper,
	)

	useCaseCreateUser := usecaseauth.NewCreateUser(
		dbConnManager,
		checkUserUseCase,
		storageSecureOperation,
		notifierAPI,
		locker,
		contactAddressParser,
		useCaseErrorWrapper,
		mapping.OptionUserRealmsToConfirmCreateUserRealms(userRealms),
	)

	useCaseConfirmAuthSession := usecaseauth.NewCreateSession(
		dbConnManager,
		checkUserUseCase,
		storageSecureOperation,
		notifierAPI,
		contactAddressParser,
		service.NewFactoryConfirm2FA(
			storageUser,
			storageAuth2fa,
			action.NewConfirmBy2fa(
				[]action.Option{
					action.WithMaxAttempts(5), // TODO: в настройки
					action.WithExpiry(30 * time.Minute),
				},
				[]action.Option{
					action.WithMaxAttempts(5), // TODO: в настройки
					action.WithExpiry(30 * time.Minute),
				},
			),
			useCaseErrorWrapper,
		),
		useCaseErrorWrapper,
		mapping.OptionUserRealmsToConfirmCreateSessionRealms(userRealms),
	)

	useCaseCreateSession := session.NewSession(
		dbConnManager,
		storageAuthToken,
		storageUserActivityStat,
		storageSecureOperation,
		handler.NewCreateUser(
			dbConnManager,
			storageUser,
			storageUserRealm,
			notifierAPI,
			useCaseErrorWrapper,
			logger,
		),
		handler.NewBeforeAuthUser(
			storageUser,
			storageUserRealm,
			notifierAPI,
			useCaseErrorWrapper,
			logger,
		),
		eventEmitter,
		useCaseErrorWrapper,
		logger,
		mapping.OptionUserRealmsToCreateSessionRealms(userRealms, jwtConfig),
	)

	useCaseUserInfo := usecaseauth.NewUserInfo(
		dbConnManager,
		storageUser,
		storageAuth2fa,
		storageUserActivityStat,
		storageUserRealm,
		useCaseErrorWrapper,
	)

	controller := httpv1.NewAuth(
		requestParser,
		responseSender,
		useCaseCreateUser,
		useCaseConfirmAuthSession,
		useCaseConfirmOperation,
		useCaseCreateSession,
		useCaseUserInfo,
		bag.NewOperationResponse(withDebugInfo),
	)

	return controller, nil
}
