package pub

import (
	"time"

	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-components/mrauth/bag/contactaddress"
	"github.com/mondegor/go-components/mrauth/bag/crypt"
	"github.com/mondegor/go-components/mrauth/component/secureoperation"
	"github.com/mondegor/go-components/mrauth/component/secureoperation/action"
	"github.com/mondegor/go-components/mrauth/service"
	"github.com/mondegor/go-components/mrauth/usecase/check"
	"github.com/mondegor/go-components/mrauth/usecase/security"
	"github.com/mondegor/go-components/mrauth/usecase/security/handler"
	"github.com/mondegor/go-sysmess/mrlib/crypt/password"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1"
	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/bag"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

func createUnitSecurity(opts auth.Options) ([]mrserver.HttpController, error) {
	var list []mrserver.HttpController

	if c, err := newUnitSecurity(opts); err != nil {
		return nil, err
	} else {
		list = append(list, c)
	}

	return list, nil
}

//nolint:unparam
func newUnitSecurity(opts auth.Options) (*httpv1.Security, error) {
	useCase := security.NewChangeProperty(
		opts.DBConnManager,
		createSecureOperationPostgres(opts),
		check.NewAuthHelper(
			createCheckUserPostgres(opts),
			createUserRealmPostgres(opts),
			contactaddress.NewParser(), // ??????
			opts.UsecaseErrorWrapper,
		),
		opts.NotifierAPI,
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
		opts.UsecaseErrorWrapper,
	)

	useCaseApplyOperation := security.NewApplyOperation(
		opts.DBConnManager,
		createSecureOperationPostgres(opts),
		opts.UsecaseErrorWrapper,
		map[string]mrauth.OperationHandler{
			secureoperation.NameConfirmChangeEmail: handler.NewChangeEmail(
				opts.DBConnManager,
				createUserPostgres(opts),
				opts.NotifierAPI,
				opts.UsecaseErrorWrapper,
			),
			secureoperation.NameConfirmChangePhone: handler.NewChangePhone(
				opts.DBConnManager,
				createUserPostgres(opts),
				opts.NotifierAPI,
				opts.UsecaseErrorWrapper,
			),
			secureoperation.NameConfirmChangePassword: handler.NewChangePassword(
				createAuth2faPostgres(opts),
				opts.NotifierAPI,
				opts.UsecaseErrorWrapper,
				opts.Logger,
			),
			secureoperation.NameConfirmDisable2FA: handler.NewDisable2FA(
				opts.DBConnManager,
				createAuth2faPostgres(opts),
				opts.NotifierAPI,
				opts.UsecaseErrorWrapper,
			),
		},
	)

	useCaseApplyTOTPGenerator := security.NewApplyTOTPGenerator(
		opts.DBConnManager,
		createAuth2faPostgres(opts),
		createSecureOperationPostgres(opts),
		opts.NotifierAPI,
		opts.UsecaseErrorWrapper,
		"PrintShopApp", // TODO:
	)

	controller := httpv1.NewSecurity(
		opts.RequestParsers.Parser,
		opts.ResponseSender,
		useCase,
		useCaseApplyTOTPGenerator,
		useCaseApplyOperation,
		bag.NewOperationResponse(opts.WithDebugInfo),
	)

	return controller, nil
}
