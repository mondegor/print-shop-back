package pub

import (
	"time"

	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-components/mrauth/bag/contactaddress"
	"github.com/mondegor/go-components/mrauth/bag/crypt"
	"github.com/mondegor/go-components/mrauth/component/secureoperation"
	"github.com/mondegor/go-components/mrauth/component/secureoperation/action"
	"github.com/mondegor/go-components/mrauth/repository"
	"github.com/mondegor/go-components/mrauth/service"
	"github.com/mondegor/go-components/mrauth/usecase/check"
	"github.com/mondegor/go-components/mrauth/usecase/security"
	"github.com/mondegor/go-components/mrauth/usecase/security/handler"
	"github.com/mondegor/go-components/mrnotifier"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlib/crypt/password"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/bag"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

func initSecurityController(
	logger mrlog.Logger,
	useCaseErrorWrapper mrerr.UseCaseErrorWrapper,
	dbConnManager mrstorage.DBConnManager,
	storageUser *repository.UserPostgres,
	storageCheckUser *repository.CheckUserPostgres,
	storageUserRealm *repository.UserRealmPostgres,
	storageAuth2fa *repository.Auth2faPostgres,
	storageSecureOperation *repository.SecureOperationPostgres,
	requestParser *validate.Parser,
	responseFileSender mrserver.FileResponseSender,
	notifierAPI mrnotifier.NoticeProducer,
	withDebugInfo bool,
) (mrserver.HttpController, error) {
	useCase := security.NewChangeProperty(
		dbConnManager,
		storageSecureOperation,
		check.NewAuthHelper(
			storageCheckUser,
			storageUserRealm,
			contactaddress.NewParser(), // ??????
			useCaseErrorWrapper,
		),
		notifierAPI,
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
		secureoperation.NewChangeEmail(
			crypt.NewTokenGenerator(64),
			crypt.NewCodeGenerator(6),
			action.WithMaxAttempts(5), // TODO: в настройки
			action.WithExpiry(30*time.Minute),
		),
		secureoperation.NewChangePhone(
			crypt.NewTokenGenerator(64),
			crypt.NewCodeGenerator(6),
			action.WithMaxAttempts(5), // TODO: в настройки
			action.WithExpiry(30*time.Minute),
		),
		secureoperation.NewChangePassword(
			crypt.NewTokenGenerator(64),
			crypt.NewCodeGenerator(6),
			action.WithMaxAttempts(5), // TODO: в настройки
			action.WithExpiry(30*time.Minute),
		),
		secureoperation.NewChangeTOTP(
			crypt.NewTokenGenerator(64),
			crypt.NewCodeGenerator(6),
			action.WithMaxAttempts(5), // TODO: в настройки
			action.WithExpiry(30*time.Minute),
		),
		secureoperation.NewDisable2FA(
			crypt.NewTokenGenerator(64),
			crypt.NewCodeGenerator(6),
			action.WithMaxAttempts(5), // TODO: в настройки
			action.WithExpiry(30*time.Minute),
		),
		func() string {
			return password.NewGenerator().Generate(16, password.CharAll) // TODO: в настройки
		},
		useCaseErrorWrapper,
	)

	useCaseApplyOperation := security.NewApplyOperation(
		dbConnManager,
		storageSecureOperation,
		useCaseErrorWrapper,
		map[string]mrauth.OperationHandler{
			secureoperation.NameConfirmChangeEmail: handler.NewChangeEmail(
				dbConnManager,
				storageUser,
				notifierAPI,
				useCaseErrorWrapper,
			),
			secureoperation.NameConfirmChangePhone: handler.NewChangePhone(
				dbConnManager,
				storageUser,
				notifierAPI,
				useCaseErrorWrapper,
			),
			secureoperation.NameConfirmChangePassword: handler.NewChangePassword(
				storageAuth2fa,
				notifierAPI,
				useCaseErrorWrapper,
				logger,
			),
			secureoperation.NameConfirmDisable2FA: handler.NewDisable2FA(
				dbConnManager,
				storageAuth2fa,
				notifierAPI,
				useCaseErrorWrapper,
			),
		},
	)

	useCaseApplyTOTPGenerator := security.NewApplyTOTPGenerator(
		dbConnManager,
		storageAuth2fa,
		storageSecureOperation,
		notifierAPI,
		useCaseErrorWrapper,
		"PrintShopApp", // TODO:
	)

	controller := httpv1.NewSecurity(
		requestParser,
		responseFileSender,
		useCase,
		useCaseApplyTOTPGenerator,
		useCaseApplyOperation,
		bag.NewOperationResponse(withDebugInfo),
	)

	return controller, nil
}
