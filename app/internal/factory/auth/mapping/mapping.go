package mapping

import (
	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-components/mrauth/bag/crypt"
	"github.com/mondegor/go-components/mrauth/bag/jwt"
	bagsession "github.com/mondegor/go-components/mrauth/bag/session"
	"github.com/mondegor/go-components/mrauth/component/secureoperation"
	"github.com/mondegor/go-components/mrauth/component/secureoperation/action"
	usecaseauth "github.com/mondegor/go-components/mrauth/usecase/auth"
	"github.com/mondegor/go-components/mrauth/usecase/session"

	"github.com/mondegor/print-shop-back/internal/factory/auth"
)

// OptionUserRealmsToConfirmCreateUserRealms - comment func.
func OptionUserRealmsToConfirmCreateUserRealms(realms []auth.UserRealm) []usecaseauth.CreateUserRealm {
	mappedRealms := make([]usecaseauth.CreateUserRealm, 0, len(realms))

	for _, realm := range realms {
		// добавляются только области пользователей, с поддержкой регистрации
		if realm.RegisterUserKind == "none" {
			continue
		}

		mappedRealms = append(
			mappedRealms,
			usecaseauth.CreateUserRealm{
				Name: realm.Name,
				Operation: secureoperation.NewCreateUser(
					realm.Name,
					realm.RegisterUserKind,
					crypt.NewTokenGenerator(int(realm.AuthToken.Length)),
					crypt.NewCodeGenerator(int(realm.OperationConfirm.CodeLength)),
					action.WithMaxAttempts(realm.OperationConfirm.SendByEmail.MaxAttempts),
					action.WithMaxResends(realm.OperationConfirm.SendByEmail.MaxResends),
					action.WithMinResendTime(realm.OperationConfirm.SendByEmail.MinResendTime),
					action.WithExpiry(realm.OperationConfirm.SessionExpiry),
				),
			},
		)
	}

	return mappedRealms
}

// OptionUserRealmsToConfirmCreateSessionRealms - comment func.
func OptionUserRealmsToConfirmCreateSessionRealms(realms []auth.UserRealm) []usecaseauth.CreateSessionRealm {
	mappedRealms := make([]usecaseauth.CreateSessionRealm, 0, len(realms))

	for _, realm := range realms {
		mappedRealms = append(
			mappedRealms,
			usecaseauth.CreateSessionRealm{
				Name: realm.Name,
				Operation: secureoperation.NewAuthorizeUser(
					crypt.NewTokenGenerator(int(realm.AuthToken.Length)),
					crypt.NewCodeGenerator(int(realm.OperationConfirm.CodeLength)),
					secureoperation.WithAuthorizeUserConfirmByEmailOpts(
						action.WithMaxAttempts(realm.OperationConfirm.SendByEmail.MaxAttempts),
						action.WithMaxResends(realm.OperationConfirm.SendByEmail.MaxResends),
						action.WithMinResendTime(realm.OperationConfirm.SendByEmail.MinResendTime),
						action.WithExpiry(realm.OperationConfirm.SessionExpiry),
					),
					secureoperation.WithAuthorizeUserConfirmByPhoneOpts(
						action.WithMaxAttempts(realm.OperationConfirm.SendByPhone.MaxAttempts),
						action.WithMaxResends(realm.OperationConfirm.SendByPhone.MaxResends),
						action.WithMinResendTime(realm.OperationConfirm.SendByPhone.MinResendTime),
						action.WithExpiry(realm.OperationConfirm.SessionExpiry),
					),
					secureoperation.WithAuthorizeUserConfirmPhoneByEmail(true),
				),
			},
		)
	}

	return mappedRealms
}

// OptionUserRealmsToCreateSessionRealms - comment func.
func OptionUserRealmsToCreateSessionRealms(realms []auth.UserRealm, jwt2 auth.JWT) []session.CreateSessionRealm {
	mappedRealms := make([]session.CreateSessionRealm, 0, len(realms))

	for _, realm := range realms {
		var tokenIssuer mrauth.TokenIssuer

		switch realm.AuthToken.AccessType {
		case "jwt":
			tokenIssuer = jwt.NewTokenIssuer(
				crypt.NewTokenGenerator(int(realm.AuthToken.Length)),
				realm.AuthToken.AccessExpiry,
				realm.AuthToken.RefreshExpiry,
				jwt2.Method,
				jwt2.Secret,
			)
		default:
			tokenIssuer = bagsession.NewTokenIssuer(
				crypt.NewTokenGenerator(int(realm.AuthToken.Length)),
				realm.AuthToken.AccessExpiry,
				realm.AuthToken.RefreshExpiry,
			)
		}

		mappedRealms = append(
			mappedRealms,
			session.CreateSessionRealm{
				Name:        realm.Name,
				TokenIssuer: tokenIssuer,
			},
		)
	}

	return mappedRealms
}
