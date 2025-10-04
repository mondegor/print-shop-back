package httpv1

import (
	"net/http"

	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
	checkLoginURL            = "/v1/check/login"
	checkPasswordStrengthURL = "/v1/check/password-strength" //nolint:gosec
)

// Check - comment struct.
type Check struct {
	parser  validate.RequestParser
	sender  mrserver.ResponseSender
	useCase mrauth.CheckUserUseCase
}

// NewCheck - создаёт объект Check.
func NewCheck(
	parser validate.RequestParser,
	sender mrserver.ResponseSender,
	useCase mrauth.CheckUserUseCase,
) *Check {
	return &Check{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера Check.
func (ht *Check) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: checkLoginURL, Func: ht.CheckLogin},
		{Method: http.MethodPost, URL: checkPasswordStrengthURL, Func: ht.CheckPasswordStrength},
	}
}

// CheckLogin - comment method.
func (ht *Check) CheckLogin(w http.ResponseWriter, r *http.Request) error {
	request := model.CheckLoginRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	if err := ht.useCase.CheckAvailability(r.Context(), request.Realm, request.UserLogin); err != nil {
		if mrauth.ErrEmailAlreadyExists.Is(err) || mrauth.ErrPhoneAlreadyExists.Is(err) {
			return mrerr.NewCustomError("userLogin", err)
		}

		return err
	}

	return ht.sender.SendNoContent(w)
}

// CheckPasswordStrength - comment method.
func (ht *Check) CheckPasswordStrength(w http.ResponseWriter, r *http.Request) error {
	request := model.CheckPasswordStrengthRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	strength, err := ht.useCase.CheckPasswordStrength(r.Context(), request.Password)
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		model.CheckPasswordStrengthResponse{
			Strength: strength.String(),
		},
	)
}
