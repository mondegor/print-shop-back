package httpv1

import (
	"net/http"
	"time"

	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-components/mrauth/entity"
	"github.com/mondegor/go-components/mrauth/enum"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlib/casttype"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrreq"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/bag"
	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
	authSignupURL  = "/v1/signup"
	authSigninURL  = "/v1/signin"
	authSessionURL = "/v1/session"
	authUserURL    = "/v1/user"
)

type (
	// Auth - comment struct.
	Auth struct {
		parser                   validate.RequestParser
		sender                   mrserver.ResponseSender
		refreshTokenCookie       *bag.RefreshTokenCookie
		useCaseCreateUser        mrauth.ConfirmCreateUserUseCase
		useCaseConfirmAuthUser   mrauth.ConfirmAuthUserUseCase
		useCaseConfirmOperation  mrauth.ConfirmOperationUseCase
		useCaseSession           mrauth.SessionUseCase
		useCaseUserInfo          mrauth.UserInfoUseCase
		factoryOperationResponse factoryOperationResponse
	}

	factoryOperationResponse interface {
		NewConfirmOperation(operation entity.SecureOperation, message string) model.WaitingConfirmOperationResponse
		NewErrorConfirmOperation(operation entity.SecureOperation, lz mrcore.Localizer, err error) model.ErrorConfirmOperationResponse
	}
)

// NewAuth - создаёт контроллер Auth.
func NewAuth(
	parser validate.RequestParser,
	sender mrserver.ResponseSender,
	useCaseCreateUser mrauth.ConfirmCreateUserUseCase,
	useCaseConfirmAuthUser mrauth.ConfirmAuthUserUseCase,
	useCaseConfirmOperation mrauth.ConfirmOperationUseCase,
	useCaseSession mrauth.SessionUseCase,
	useCaseUserInfo mrauth.UserInfoUseCase,
	factoryOperationResponse factoryOperationResponse,
) *Auth {
	return &Auth{
		parser: parser,
		sender: sender,
		refreshTokenCookie: bag.NewRefreshTokenCookie(
			"RTID",           // TODO: options !!!!!!!
			"localhost",      // TODO: options !!!!!!!
			"/",              // TODO: options !!!!!!!
			180*24*time.Hour, // TODO: options !!!!!!!
		),
		useCaseCreateUser:        useCaseCreateUser,
		useCaseConfirmAuthUser:   useCaseConfirmAuthUser,
		useCaseConfirmOperation:  useCaseConfirmOperation,
		useCaseSession:           useCaseSession,
		useCaseUserInfo:          useCaseUserInfo,
		factoryOperationResponse: factoryOperationResponse,
	}
}

// Handlers - возвращает обработчики контроллера Auth.
func (ht *Auth) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: authSignupURL, Permission: mrserver.PermissionGuestOnly, Func: ht.Signup},
		{Method: http.MethodPost, URL: authSigninURL, Permission: mrserver.PermissionGuestOnly, Func: ht.Signin},
		{Method: http.MethodPost, URL: authSessionURL, Permission: mrserver.PermissionGuestOnly, Func: ht.OpenSession},
		{Method: http.MethodPatch, URL: authSessionURL, Permission: mrserver.PermissionAnyUser, Func: ht.ContinueSession},
		{Method: http.MethodDelete, URL: authSessionURL, Func: ht.CloseSession},
		{Method: http.MethodGet, URL: authUserURL, Func: ht.UserInfo},
	}
}

// Signup - comment method.
func (ht *Auth) Signup(w http.ResponseWriter, r *http.Request) error {
	request := model.CreateUserRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	lz := ht.parser.Localizer(r)

	op, err := ht.useCaseCreateUser.Perform(r.Context(), request.Realm, lz.Language(), request.UserEmail)
	if err != nil {
		if mrauth.ErrEmailAlreadyExists.Is(err) {
			return mrerr.NewCustomError("userEmail", err)
		}

		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ht.factoryOperationResponse.NewConfirmOperation(
			op,
			lz.Translate("Confirm the creation of the user"),
		),
	)
}

// Signin - comment method.
func (ht *Auth) Signin(w http.ResponseWriter, r *http.Request) error {
	request := model.AuthorizeUserRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	lz := ht.parser.Localizer(r)

	// TODO: ограничивать частую отправку событий на авторизацию
	// TODO: писать, что код подтверждения уже был выслан, повторить попытку можно через N минут

	// TODO: проверить, что открыто не более X сессий

	op, err := ht.useCaseConfirmAuthUser.Perform(r.Context(), request.Realm, lz.Language(), request.UserLogin)
	if err != nil {
		if mrauth.ErrLoginNotExists.Is(err) {
			return mrerr.NewCustomError("userLogin", err)
		}

		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ht.factoryOperationResponse.NewConfirmOperation(
			op,
			lz.Translate("Confirm your identity to sign in"),
		),
	)
}

// OpenSession - comment method.
func (ht *Auth) OpenSession(w http.ResponseWriter, r *http.Request) error {
	request := model.LoginByTokenRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	lz := ht.parser.Localizer(r)

	// TODO: useCaseConfirmOperation и useCaseSession вложить в useCaseGroup

	// иначе операцию необходимо сначала подтвердить
	op, err := ht.useCaseConfirmOperation.Perform(r.Context(), lz.Language(), request.Token, request.Secret)
	if err != nil {
		if mrauth.ErrConfirmCodeIsIncorrect.Is(err) || mrauth.ErrNoAttemptsToConfirmOperation.Is(err) {
			return ht.sender.Send(
				w,
				http.StatusBadRequest,
				ht.factoryOperationResponse.NewErrorConfirmOperation(op, lz, mrerr.NewCustomError("secret", err)),
			)
		}

		if mr.ErrUseCaseEntityNotFound.Is(err) {
			return mrauth.ErrTokenNotFoundOrExpired.Wrap(err)
		}

		return err
	}

	// если необходимо дополнительное подтверждение (2fa)
	if op.Status == enum.OperationStatusOpened {
		return ht.sender.Send(
			w,
			http.StatusOK,
			ht.factoryOperationResponse.NewConfirmOperation(
				op,
				lz.Translate("your account has been success registered"),
			),
		)
	}

	// если операция была подтверждена
	tk, err := ht.useCaseSession.Open(r.Context(), ht.parser.DetailedIP(r), op)
	if err != nil {
		return ht.wrapError(err)
	}

	if request.UseCookie {
		// for web version
		ht.refreshTokenCookie.SetValue(w, tk.RefreshToken)
		tk.RefreshToken = ""
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		model.SuccessAccessResponse{
			AccessToken:  tk.AccessToken,
			ExpiresIn:    uint32(tk.ExpiresIn / time.Second), //nolint:gosec
			RefreshToken: tk.RefreshToken,
		},
	)
}

// ContinueSession - comment method.
func (ht *Auth) ContinueSession(w http.ResponseWriter, r *http.Request) error {
	refreshToken := ht.refreshTokenCookie.GetValue(r)
	useCookie := true

	if refreshToken == "" {
		request := model.ContinueSessionRequest{}

		if err := ht.parser.Validate(r, &request); err != nil {
			return err
		}

		refreshToken = request.RefreshToken
		useCookie = false
	}

	tk, err := ht.useCaseSession.Continue(r.Context(), ht.parser.Localizer(r).Language(), refreshToken)
	if err != nil {
		if mr.ErrUseCaseEntityNotFound.Is(err) {
			return mrauth.ErrTokenNotFoundOrExpired.Wrap(err)
		}

		return err
	}

	if useCookie {
		// for web version
		ht.refreshTokenCookie.SetValue(w, tk.RefreshToken)
		tk.RefreshToken = ""
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		model.SuccessAccessResponse{
			AccessToken:  tk.AccessToken,
			ExpiresIn:    uint32(tk.ExpiresIn / time.Second), //nolint:gosec
			RefreshToken: tk.RefreshToken,
		},
	)
}

// CloseSession - comment method.
func (ht *Auth) CloseSession(w http.ResponseWriter, r *http.Request) error {
	accessToken := mrreq.ParseAccessToken(r.Header)
	if accessToken == "" {
		return mr.ErrHttpClientUnauthorized.New()
	}

	if err := ht.useCaseSession.Close(r.Context(), accessToken); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

// UserInfo - comment method.
func (ht *Auth) UserInfo(w http.ResponseWriter, r *http.Request) error {
	info, err := ht.useCaseUserInfo.Get(r.Context(), ht.parser.UserID(r))
	if err != nil {
		return err
	}

	realms := make([]model.UserRealm, 0, len(info.Realms))
	for _, realm := range info.Realms {
		realms = append(
			realms,
			model.UserRealm{
				Name:     realm.Realm,
				UserKind: realm.Kind,
			},
		)
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		model.UserInfoResponse{
			Email:        info.User.Email,
			Phone:        casttype.UintToPhone(info.User.Phone),
			LangCode:     info.User.LangCode,
			LastLoginIP:  info.Stat.LastLoginIP.Real.String(),
			LastLoggedAt: info.Stat.LastLoggedAt.Round(1 * time.Second).Format(time.RFC3339),
			Auth2faType:  info.Auth2fa.Type,
			Realms:       realms,
			Status:       info.User.Status,
		},
	)
}

func (ht *Auth) wrapError(err error) error {
	return err
}
