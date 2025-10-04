package httpv1

import (
	"net/http"

	"github.com/mondegor/go-components/mrauth"
	"github.com/mondegor/go-components/mrauth/enum"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/auth/section/pub/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/pkg/validate"
)

const (
	operationConfirmURL = "/v1/operation/confirm"
	operationResendURL  = "/v1/operation/resend"
	// operationRevokeURL  = "/v1/operation/revoke".
)

// Operation - comment struct.
type Operation struct {
	parser                   validate.RequestParser
	sender                   mrserver.ResponseSender
	useCaseConfirmOperation  mrauth.ConfirmOperationUseCase
	useCaseResendConfirmCode mrauth.ResendConfirmCodeUseCase
	factoryOperationResponse factoryOperationResponse
}

// NewOperation - создаёт объект Operation.
func NewOperation(
	parser validate.RequestParser,
	sender mrserver.ResponseSender,
	useCaseConfirmOperation mrauth.ConfirmOperationUseCase,
	useCaseResendConfirmCode mrauth.ResendConfirmCodeUseCase,
	factoryOperationResponse factoryOperationResponse,
) *Operation {
	return &Operation{
		parser:                   parser,
		sender:                   sender,
		useCaseConfirmOperation:  useCaseConfirmOperation,
		useCaseResendConfirmCode: useCaseResendConfirmCode,
		factoryOperationResponse: factoryOperationResponse,
	}
}

// Handlers - возвращает обработчики контроллера Operation.
func (ht *Operation) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPatch, URL: operationConfirmURL, Permission: mrserver.PermissionAnyUser, Func: ht.Confirm},
		{Method: http.MethodPatch, URL: operationResendURL, Permission: mrserver.PermissionAnyUser, Func: ht.Resend},
		// {Method: http.MethodPatch, URL: operationRevokeURL, Func: ht.Revoke},
	}
}

// Confirm - comment method.
func (ht *Operation) Confirm(w http.ResponseWriter, r *http.Request) error {
	request := model.ConfirmOperationRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	lz := ht.parser.Localizer(r)

	op, err := ht.useCaseConfirmOperation.Perform(r.Context(), lz.Language(), request.Token, request.Secret)
	if err != nil {
		if !mrauth.ErrConfirmCodeIsIncorrect.Is(err) && !mrauth.ErrNoAttemptsToConfirmOperation.Is(err) {
			return ht.wrapError(err)
		}

		return ht.sender.Send(
			w,
			http.StatusBadRequest,
			ht.factoryOperationResponse.NewErrorConfirmOperation(op, lz, err),
		)
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
	return ht.sender.SendNoContent(w)
}

// Resend - comment method.
func (ht *Operation) Resend(w http.ResponseWriter, r *http.Request) error {
	request := model.OperationTokenRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	lz := ht.parser.Localizer(r)

	op, err := ht.useCaseResendConfirmCode.Perform(r.Context(), lz.Language(), request.Token)
	if err != nil {
		if mr.ErrUseCaseEntityNotFound.Is(err) {
			return mrauth.ErrTokenNotFoundOrExpired.Wrap(err)
		}

		if mrauth.ErrSendingNewMessagesIsTemporarilyRestricted.Is(err) {
			return ht.sender.Send(
				w,
				http.StatusBadRequest,
				ht.factoryOperationResponse.NewErrorConfirmOperation(op, lz, mrerr.NewCustomError("token", err)),
			)
		}

		return ht.wrapError(err)
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ht.factoryOperationResponse.NewConfirmOperation(
			op,
			lz.Translate("The confirmation code has been sent successfully"),
		),
	)
}

// func (ht *Operation) Revoke(w http.ResponseWriter, r *http.Request) error {
// 	request := OperationRequest{}
//
// 	if err := ht.parser.Validate(r, &request); err != nil {
// 		return err
// 	}
//
// 	if err := ht.useCase.Revoke(r.Context(), request.AuthToken); err != nil {
// 		return ht.wrapError(err)
// 	}
//
// 	return ht.sender.SendNoContent(w)
// }

func (ht *Operation) wrapError(err error) error {
	// ConfirmCode is not correct
	// operation already confirmed | operation is not opened
	return err
}
