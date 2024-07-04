package httpv1

import (
	"net/http"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
	"github.com/mondegor/print-shop-back/pkg/view"

	"github.com/mondegor/go-components/mrsort"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	formElementListURL     = "/v1/controls/submit-form-elements"
	formElementItemURL     = "/v1/controls/submit-form-elements/{id}"
	formElementItemMoveURL = "/v1/controls/submit-form-elements/{id}/move"
)

type (
	// FormElement - comment struct.
	FormElement struct {
		parser  validate.RequestSubmitFormParser
		sender  mrserver.ResponseSender
		useCase usecase.FormElementUseCase
	}
)

// NewFormElement - создаёт контроллер FormElement.
func NewFormElement(parser validate.RequestSubmitFormParser, sender mrserver.ResponseSender, useCase usecase.FormElementUseCase) *FormElement {
	return &FormElement{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
	}
}

// Handlers - возвращает обработчики контроллера FormElement.
func (ht *FormElement) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodPost, URL: formElementListURL, Func: ht.Create},

		{Method: http.MethodGet, URL: formElementItemURL, Func: ht.Get},
		{Method: http.MethodPatch, URL: formElementItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: formElementItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: formElementItemMoveURL, Func: ht.Move},
	}
}

// Get - comment method.
func (ht *FormElement) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// Create - comment method.
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

	itemID, err := ht.useCase.Create(r.Context(), item)
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(
		w,
		http.StatusCreated,
		view.SuccessCreatedItemInt32Response{
			ItemID: itemID,
		},
	)
}

// Store - comment method.
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

// Remove - comment method.
func (ht *FormElement) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// Move - comment method.
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
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrFormElementNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if module.ErrSubmitFormNotFound.Is(err) ||
		module.ErrFormElementDetailingNotAllowed.Is(err) {
		return mrerr.NewCustomError("formId", err)
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if module.ErrFormElementParamNameAlreadyExists.Is(err) {
		return mrerr.NewCustomError("paramName", err)
	}

	if api.ErrElementTemplateNotFound.Is(err) ||
		api.ErrElementTemplateIsDisabled.Is(err) {
		return mrerr.NewCustomError("templateId", err)
	}

	return err
}

func (ht *FormElement) wrapErrorNode(r *http.Request, err error) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return module.ErrFormElementNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrsort.ErrAfterNodeNotFound.Is(err) {
		return mrerr.NewCustomError("afterNodeId", err)
	}

	return err
}
