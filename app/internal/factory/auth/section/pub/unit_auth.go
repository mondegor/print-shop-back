package pub

import (
	"time"

	"github.com/mondegor/go-components/mrauth/bag/contactaddress"
	"github.com/mondegor/go-components/mrauth/bag/crypt"
	"github.com/mondegor/go-components/mrauth/component/secureoperation"
	"github.com/mondegor/go-components/mrauth/component/secureoperation/action"
	"github.com/mondegor/go-components/mrauth/service"
	usecaseauth "github.com/mondegor/go-components/mrauth/usecase/auth"
	"github.com/mondegor/go-components/mrauth/usecase/check"
	"github.com/mondegor/go-components/mrauth/usecase/operation"
	"github.com/mondegor/go-components/mrauth/usecase/session"
	"github.com/mondegor/go-components/mrauth/usecase/session/handler"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/bag"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
	"github.com/mondegor/print-shop-back/internal/factory/auth/mapping"
)

func createUnitAuth(opts auth.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitAuth(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

//nolint:unparam
func newUnitAuth(opts auth.Options) (*httpv1.Auth, error) {
	contactAddressParser := contactaddress.NewParser()
	checkUserUseCase := check.NewAuthHelper(
		createCheckUserPostgres(opts),
		createUserRealmPostgres(opts),
		contactAddressParser,
		opts.UsecaseErrorWrapper,
	)

	useCaseCreateUser := usecaseauth.NewCreateUser(
		opts.DBConnManager,
		checkUserUseCase,
		createSecureOperationPostgres(opts),
		opts.NotifierAPI,
		opts.Locker,
		contactAddressParser,
		opts.UsecaseErrorWrapper,
		mapping.OptionUserRealmsToConfirmCreateUserRealms(opts.UserRealms),
	)

	useCaseConfirmAuthSession := usecaseauth.NewCreateSession(
		opts.DBConnManager,
		checkUserUseCase,
		createSecureOperationPostgres(opts),
		opts.NotifierAPI,
		contactAddressParser,
		service.NewFactoryConfirm2FA(
			createUserPostgres(opts),
			createAuth2faPostgres(opts),
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
			opts.UsecaseErrorWrapper,
		),
		opts.UsecaseErrorWrapper,
		mapping.OptionUserRealmsToConfirmCreateSessionRealms(opts.UserRealms),
	)

	useCaseCreateSession := session.NewSession(
		opts.DBConnManager,
		createAuthTokenPostgres(opts),
		createUserActivityStatPostgres(opts),
		createSecureOperationPostgres(opts),
		handler.NewCreateUser(
			opts.DBConnManager,
			createUserPostgres(opts),
			createUserRealmPostgres(opts),
			opts.NotifierAPI,
			opts.UsecaseErrorWrapper,
			opts.Logger,
		),
		handler.NewBeforeAuthUser(
			createUserPostgres(opts),
			createUserRealmPostgres(opts),
			opts.NotifierAPI,
			opts.UsecaseErrorWrapper,
			opts.Logger,
		),
		opts.EventEmitter,
		opts.UsecaseErrorWrapper,
		opts.Logger,
		mapping.OptionUserRealmsToCreateSessionRealms(opts.UserRealms, opts.JWT),
	)

	// ДУБЛИРУЕТСЯ!!!!
	useCaseConfirmOperation := operation.NewConfirmOperation(
		opts.DBConnManager,
		createSecureOperationPostgres(opts),
		opts.NotifierAPI,
		secureoperation.NewConfirmCode(
			crypt.NewTokenGenerator(int(opts.OperationConfirm.TokenLength)), // DEFAULT
			crypt.NewCodeGenerator(int(opts.OperationConfirm.CodeLength)),   // DEFAULT
		),
		opts.UsecaseErrorWrapper,
	)

	useCaseUserInfo := usecaseauth.NewUserInfo(
		opts.DBConnManager,
		createUserPostgres(opts),
		createAuth2faPostgres(opts),
		createUserActivityStatPostgres(opts),
		createUserRealmPostgres(opts),
		opts.UsecaseErrorWrapper,
	)

	controller := httpv1.NewAuth(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCaseCreateUser,
		useCaseConfirmAuthSession,
		useCaseConfirmOperation,
		useCaseCreateSession,
		useCaseUserInfo,
		bag.NewOperationResponse(opts.WithDebugInfo),
	)

	return controller, nil
}
