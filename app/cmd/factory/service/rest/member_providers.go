package rest

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/factory/mrauth"
	"github.com/mondegor/go-components/mrauth/component/get"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/user"
	"github.com/mondegor/go-webcore/mrdebug"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

const (
	accessTokenTableName  = "printshop_auth.auth_tokens" //nolint:gosec
	accessTokenPrimaryKey = "token_name"
)

// initMemberProviders - создаёт объект mraccess.MemberProvider с указанными настройками.
func initMemberProviders(opts app.Options) (realm2provider map[string]mraccess.MemberProvider) {
	if len(opts.Cfg.AccessControl.Realms) == 0 {
		mrlog.Error(opts.Logger, "Auth: AccessControl.Realms is empty")

		return nil
	}

	realm2provider = make(map[string]mraccess.MemberProvider, len(opts.Cfg.AccessControl.Realms))
	realms := make([]auth.UserRealm, 0, len(opts.Cfg.AccessControl.Realms))
	domain2realms := make(map[string][]auth.UserRealm, len(opts.Cfg.AccessControl.Realms))

	for _, realm := range opts.Cfg.AccessControl.Realms {
		switch realm.AuthToken.AccessType {
		// если метод аутентификации указан JWT, то будут приниматься от клиентов JWT токены
		case "jwt":
			mrlog.Debug(opts.Logger, fmt.Sprintf("Auth.JWT: realm=%s, secret=%s", realm.Name, opts.Cfg.AccessControl.JWTSecret))

		// стандартный режим: будут приниматься от клиентов токены, хранящиеся в таблице accessTokenTableName
		default:
			mrlog.Debug(opts.Logger, fmt.Sprintf("Auth.Session: realm=%s, table=%s", realm.Name, accessTokenTableName))
		}

		domain := realm.Name
		if val, _, ok := strings.Cut(realm.Name, "/"); ok {
			domain = val
		}

		realms = append(realms, realm)
		domain2realms[domain] = append(domain2realms[domain], realm)

		realm2provider[realm.Name] = createMemberProviderByTokenType(
			opts,
			realm.AuthToken.AccessType,
			[]string{realm.Name},
		)
	}

	realm2provider["*"] = createMemberProviderGroup(opts, realms)

	for domain, realms := range domain2realms {
		realm2provider[domain+"/*"] = createMemberProviderGroup(opts, realms)
	}

	return realm2provider
}

func createMemberProviderGroup(opts app.Options, userRealms []auth.UserRealm) mraccess.MemberProvider {
	type2realms := make(map[string][]string, len(userRealms))

	for _, realm := range userRealms {
		type2realms[realm.AuthToken.AccessType] = append(type2realms[realm.AuthToken.AccessType], realm.Name)
	}

	memberProviders := make([]get.ProviderWithTokenType, 0, len(type2realms))

	for tokenType, realms := range type2realms {
		memberProviders = append(
			memberProviders,
			get.ProviderWithTokenType{
				TokenType: tokenType,
				Provider:  createMemberProviderByTokenType(opts, tokenType, realms),
			},
		)
	}

	if len(memberProviders) == 1 {
		return memberProviders[0].Provider
	}

	return mrauth.NewUserProvider(memberProviders...)
}

func createMemberProviderByTokenType(opts app.Options, tokenType string, allowedRealms []string) mraccess.MemberProvider {
	// в отладочном режиме можно указать произвольного пользователя, в указанном realm и с указанным kind
	for _, realm := range allowedRealms {
		if !opts.Cfg.Debugging.Debug || opts.Cfg.Debugging.AuthorizedUser.Realm != realm || opts.Cfg.Debugging.AuthorizedUser.ID == "" {
			continue
		}

		mrlog.Debug(
			opts.Logger,
			fmt.Sprintf(
				"Auth.Debug: userId=%s, realm=%s, kind=%s, lang=%s",
				opts.Cfg.Debugging.AuthorizedUser.ID,
				opts.Cfg.Debugging.AuthorizedUser.Realm,
				opts.Cfg.Debugging.AuthorizedUser.Kind,
				opts.Cfg.Debugging.AuthorizedUser.LangCode,
			),
		)

		return mrdebug.NewMemberProvider(
			user.New(
				uuid.MustParse(opts.Cfg.Debugging.AuthorizedUser.ID),
				opts.Cfg.Debugging.AuthorizedUser.Realm+"/"+opts.Cfg.Debugging.AuthorizedUser.Kind,
				opts.Cfg.Debugging.AuthorizedUser.LangCode,
			),
		)
	}

	// JWT режим: принимаются от клиентов JWT токены
	if tokenType == "jwt" {
		return mrauth.NewUserProviderJWT(
			opts.UseCaseErrorWrapper,
			opts.Cfg.AccessControl.JWTSecret,
			allowedRealms,
		)
	}

	// Session режим: принимаются от клиентов токены, хранящиеся в таблице accessTokenTableName
	return mrauth.NewUserProviderSession(
		opts.PostgresConnManager,
		opts.UseCaseErrorWrapper,
		opts.StorageErrorWrapper,
		mrsql.DBTableInfo{
			Name:       accessTokenTableName,
			PrimaryKey: accessTokenPrimaryKey,
		},
		allowedRealms,
	)
}
