package rest

import (
	"net/http"

	authvalidate "github.com/mondegor/go-components/mrauth/validate"
	authcfg "github.com/mondegor/go-components/wire/mrauth/config"
	auth "github.com/mondegor/go-components/wire/mrauth/infra/pub"
	"github.com/mondegor/go-sysmess/mraccess"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"print-shop-back/internal/app"
)

// TODO: дублирование название таблиц.
const (
	serviceAuthTokensTableName      = "printshop_auth.auth_tokens" //nolint:gosec
	serviceSecureOperationTableName = "printshop_auth.secure_operations"
	// serviceSecureOperationLogTableName = "printshop_auth.secure_operations_log".
	serviceSessionsTableName = "printshop_auth.sessions"
	serviceUsersTableName    = "printshop_auth.users"
	// serviceUsersActivityLogTableName = "printshop_auth.users_activity_log".
	serviceUsersActivityStatTableName = "printshop_auth.users_activity_stat"
	serviceUsersAuth2faTableName      = "printshop_auth.users_auth_2fa"
	serviceUsersRealmsTableName       = "printshop_auth.users_realms"
)

// RegisterRestRouterAuthHandlers - регистрирует в указанном роутере обработчики секции AuthAPI.
func RegisterRestRouterAuthHandlers(
	router mrserver.HttpRouter,
	opts app.Options,
	actionGroup mraccess.ActionGroup,
	userProvider mraccess.UserProvider,
) error {
	router.HandlerFunc(http.MethodGet, actionGroup.BasePath, mrresp.HandlerGetStatusOkAsJSON(opts.Logger))

	controllers, err := initing.CreateHttpControllers(
		opts.Logger,
		getAuthAPIControllers(opts),
		initing.WithCheckAccessMiddleware(opts.Logger, actionGroup, userProvider, opts.PermsProvider),
	)
	if err != nil {
		return err
	}

	router.Register(controllers...)

	return nil
}

func getAuthAPIControllers(opts app.Options) []initing.HttpModule {
	return []initing.HttpModule{
		auth.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.Locker,
			// opts.RequestParsers.Parser,
			authvalidate.NewParser( // TODO: объединить со стандартным Parser или сделать свой? Может там нужно меньше парсеров
				opts.RequestParsers.Int64,
				opts.RequestParsers.Uint64,
				opts.RequestParsers.String,
				opts.RequestParsers.UUID,
				opts.RequestParsers.Validator,
				opts.RequestParsers.ClientIP,
				opts.RequestParsers.User,
				opts.RequestParsers.Locale,
			),
			opts.ResponseSenders.Sender,
			opts.ResponseSenders.FileSender,
			opts.NotifierAPI,
			opts.Cfg.AccessControl.Realms,
			opts.Cfg.AccessControl.DefaultOperationConfirm,
			authcfg.JWT{
				Method: opts.Cfg.AccessControl.JWTMethod,
				Secret: []byte(opts.Cfg.AccessControl.JWTSecret),
			},
			nil, // appResolver
			nil, // locationResolver
			serviceAuthTokensTableName,
			serviceSecureOperationTableName,
			// serviceSecureOperationLogTableName,
			serviceSessionsTableName,
			serviceUsersTableName,
			// serviceUsersActivityLogTableName,
			serviceUsersActivityStatTableName,
			serviceUsersAuth2faTableName,
			serviceUsersRealmsTableName,
			opts.DebugFunc,
		),
	}
}
