package http_v1

import (
	"net/http"
	view_shared "print-shop-back/internal/modules/controls/submit-form/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/submit-form/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/controls/submit-form/usecase/shared"
	"print-shop-back/pkg/modules/controls"
	"print-shop-back/pkg/shared/view"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	formElementListURL     = "/v1/controls/submit-form-elements"
	formElementItemURL     = "/v1/controls/submit-form-elements/:id"
	formElementItemMoveURL = "/v1/controls/submit-form-elements/:id/move"
)

type (
	FormElement struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.FormElementUseCase
	}
)

func NewFormElement(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.FormElementUseCase,
) *FormElement {
	return &FormElement{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

func (ht *FormElement) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodPost, formElementListURL, "", ht.Create},

		{http.MethodGet, formElementItemURL, "", ht.Get},
		{http.MethodPatch, formElementItemURL, "", ht.Store},
		{http.MethodDelete, formElementItemURL, "", ht.Remove},

		{http.MethodPatch, formElementItemMoveURL, "", ht.Move},
	}
}

func (ht *FormElement) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *FormElement) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateFormElementRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.FormElement{
		FormID:     request.FormID,
		TemplateID: request.TemplateID,
		ParamName:  request.ParamName,
		Caption:    request.Caption,
		Required:   &request.Required,
	}

	if itemID, err := ht.useCase.Create(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	} else {
		return ht.sender.Send(
			w,
			http.StatusCreated,
			view.SuccessCreatedItemInt32Response{
				ItemID: itemID,
			},
		)
	}
}

func (ht *FormElement) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreFormElementRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.FormElement{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		ParamName:  request.ParamName,
		Caption:    request.Caption,
		Required:   request.Required,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormElement) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormElement) Move(w http.ResponseWriter, r *http.Request) error {
	request := view.MoveItemRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	if err := ht.useCase.MoveAfterID(r.Context(), ht.getItemID(r), request.AfterNodeID); err != nil {
		return ht.wrapErrorNode(r, err)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormElement) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *FormElement) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *FormElement) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrFormElementNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if usecase_shared.FactoryErrSubmitFormNotFound.Is(err) ||
		usecase_shared.FactoryErrFormElementDetailingNotAllowed.Is(err) {
		return mrerr.NewCustomError("formId", err)
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if usecase_shared.FactoryErrFormElementParamNameAlreadyExists.Is(err) {
		return mrerr.NewCustomError("paramName", err)
	}

	if controls.FactoryErrElementTemplateNotFound.Is(err) ||
		controls.FactoryErrElementTemplateIsDisabled.Is(err) {
		return mrerr.NewCustomError("templateId", err)
	}

	return err
}

func (ht *FormElement) wrapErrorNode(r *http.Request, err error) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrFormElementNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrorderer.FactoryErrAfterNodeNotFound.Is(err) {
		return mrerr.NewCustomError("afterNodeId", err)
	}

	return err
}
