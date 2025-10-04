package httpv1

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
	securityEmailURL            = "/v1/security/email"
	securityPhoneURL            = "/v1/security/phone"
	securityPasswordURL         = "/v1/security/password"
	securityPasswordGenerateURL = "/v1/security/password/generate"
	securityTOTPGeneratorURL    = "/v1/security/totp"
	securityDisable2FAURL       = "/v1/security/disable2fa"
	securityApplyOperation      = "/v1/security/apply-operation"
)

type (
	// Security - comment struct.
	Security struct {
		parser                    validate.RequestParser
		sender                    mrserver.FileResponseSender
		useCase                   mrauth.ChangeUseCase
		useCaseApplyOperationTOTP useCaseApplyOperationTOTP
		useCaseApplyOperation     useCaseApplyOperation
		factoryOperationResponse  factoryOperationResponse
	}

	useCaseApplyOperationTOTP interface {
		Execute(ctx context.Context, userID uuid.UUID, operationToken string) (totpURL mrtype.Image, err error)
	}

	useCaseApplyOperation interface {
		Execute(ctx context.Context, userID uuid.UUID, operationToken string) error
	}
)

// NewSecurity - создаёт объект Security.
func NewSecurity(
	parser validate.RequestParser,
	sender mrserver.FileResponseSender,
	useCase mrauth.ChangeUseCase,
	useCaseApplyOperationTOTP useCaseApplyOperationTOTP,
	useCaseApplyOperation useCaseApplyOperation,
	factoryOperationResponse factoryOperationResponse,
) *Security {
	return &Security{
		parser:                    parser,
		sender:                    sender,
		useCase:                   useCase,
		useCaseApplyOperationTOTP: useCaseApplyOperationTOTP,
		useCaseApplyOperation:     useCaseApplyOperation,
		factoryOperationResponse:  factoryOperationResponse,
	}
}

// Handlers - возвращает обработчики контроллера Security.
func (ht *Security) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: securityEmailURL, Func: ht.ChangeEmail},
		{Method: http.MethodPost, URL: securityPhoneURL, Func: ht.ChangePhone},
		{Method: http.MethodPost, URL: securityPasswordURL, Func: ht.ChangePassword},
		{Method: http.MethodPost, URL: securityPasswordGenerateURL, Func: ht.GeneratePassword},
		{Method: http.MethodPost, URL: securityTOTPGeneratorURL, Func: ht.ChangeTOTPGenerator},
		{Method: http.MethodPatch, URL: securityTOTPGeneratorURL, Func: ht.ApplyTOTPGenerator},
		{Method: http.MethodPost, URL: securityDisable2FAURL, Func: ht.Disable2FA},
		{Method: http.MethodPatch, URL: securityApplyOperation, Func: ht.ApplyOperation},
	}
}

// ChangeEmail - comment method.
func (ht *Security) ChangeEmail(w http.ResponseWriter, r *http.Request) error {
	request := model.ChangeEmailRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	op, err := ht.useCase.ChangeEmail(r.Context(), ht.parser.UserID(r), request.NewEmail)
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ht.factoryOperationResponse.NewConfirmOperation(
			op,
			ht.parser.Localizer(r).Translate("your account has been success registered"),
		),
	)
}

// ChangePhone - comment method.
func (ht *Security) ChangePhone(w http.ResponseWriter, r *http.Request) error {
	request := model.ChangePhoneRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	op, err := ht.useCase.ChangePhone(r.Context(), ht.parser.UserID(r), request.NewPhone)
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ht.factoryOperationResponse.NewConfirmOperation(
			op,
			ht.parser.Localizer(r).Translate("msgChangePhoneRequestCreatedSuccessfully"),
		),
	)
}

// ChangePassword - comment method.
func (ht *Security) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	request := model.ChangePasswordRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	op, err := ht.useCase.ChangePassword(r.Context(), ht.parser.UserID(r), request.NewPassword)
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ht.factoryOperationResponse.NewConfirmOperation(
			op,
			ht.parser.Localizer(r).Translate("msgChangePasswordRequestCreatedSuccessfully"),
		),
	)
}

// GeneratePassword - comment method.
func (ht *Security) GeneratePassword(w http.ResponseWriter, r *http.Request) error {
	return ht.sender.Send(
		w,
		http.StatusOK,
		model.GeneratedPasswordResponse{
			Password: ht.useCase.GeneratePassword(r.Context()),
		},
	)
}

// ChangeTOTPGenerator - comment method.
func (ht *Security) ChangeTOTPGenerator(w http.ResponseWriter, r *http.Request) error {
	op, err := ht.useCase.ChangeTOTPGenerator(r.Context(), ht.parser.UserID(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ht.factoryOperationResponse.NewConfirmOperation(
			op,
			ht.parser.Localizer(r).Translate("msgDisable2FARequestCreatedSuccessfully"),
		),
	)
}

// ApplyTOTPGenerator - comment method.
func (ht *Security) ApplyTOTPGenerator(w http.ResponseWriter, r *http.Request) error {
	request := model.ApplyOperationRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	totpImage, err := ht.useCaseApplyOperationTOTP.Execute(r.Context(), ht.parser.UserID(r), request.Token)
	if err != nil {
		return err
	}

	return ht.sender.SendFile(
		r.Context(),
		w,
		totpImage.ToFile(),
	)
}

// Disable2FA - comment method.
func (ht *Security) Disable2FA(w http.ResponseWriter, r *http.Request) error {
	op, err := ht.useCase.Disable2FA(r.Context(), ht.parser.UserID(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ht.factoryOperationResponse.NewConfirmOperation(
			op,
			ht.parser.Localizer(r).Translate("msgDisable2FARequestCreatedSuccessfully"),
		),
	)
}

// ApplyOperation - comment method.
func (ht *Security) ApplyOperation(w http.ResponseWriter, r *http.Request) error {
	request := model.ApplyOperationRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	if err := ht.useCaseApplyOperation.Execute(r.Context(), ht.parser.UserID(r), request.Token); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}
