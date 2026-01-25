package httpv1

import (
	"net/http"

	"github.com/mondegor/go-components/mrordering"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-webcore/mrserver"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/shared/validate"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
	"github.com/mondegor/print-shop-back/pkg/view"
)

const (
	formElementListURL     = "/v1/controls/submit-form-elements"
	formElementItemURL     = "/v1/controls/submit-form-elements/{id}"
	formElementItemMoveURL = "/v1/controls/submit-form-elements/{id}/move"
)

type (
	// FormElement - comment struct.
	FormElement struct {
		parser       validate.RequestSubmitFormParser
		sender       mrserver.ResponseSender
		useCase      adm.FormElementUseCase
		errorWrapper errors.CustomWrapper
	}
)

// NewFormElement - создаёт контроллер FormElement.
func NewFormElement(parser validate.RequestSubmitFormParser, sender mrserver.ResponseSender, useCase adm.FormElementUseCase) *FormElement {
	return &FormElement{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
		errorWrapper: errors.NewCustomWrapper(
			module.ErrSubmitFormNotFound.Code(), "formId",
			module.ErrFormElementDetailingNotAllowed.Code(), "formId",
			errors.ErrUseCaseEntityVersionConflict.Code(), "tagVersion",
			module.ErrFormElementParamNameAlreadyExists.Code(), "paramName",
			api.ErrElementTemplateNotFound.Code(), "templateId",
			api.ErrElementTemplateIsDisabled.Code(), "templateId",
			mrordering.ErrAfterNodeNotFound.Code(), "afterNodeId",
		),
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
	req := CreateFormElementRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := entity.FormElement{
		FormID:     req.FormID,
		TemplateID: req.TemplateID,
		ParamName:  req.ParamName,
		Caption:    req.Caption,
		Required:   &req.Required,
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
	req := StoreFormElementRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	item := entity.FormElement{
		ID:         ht.getItemID(r),
		TagVersion: req.TagVersion,
		ParamName:  req.ParamName,
		Caption:    req.Caption,
		Required:   req.Required,
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
	req := view.MoveItemRequest{}

	if err := ht.parser.Validate(r, &req); err != nil {
		return err
	}

	if err := ht.useCase.MoveAfterID(r.Context(), ht.getItemID(r), req.AfterNodeID); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *FormElement) getItemID(r *http.Request) uint64 {
	return ht.parser.PathParamUint64(r, "id")
}

func (ht *FormElement) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *FormElement) wrapError(err error, r *http.Request) error {
	if errors.Is(err, errors.ErrUseCaseEntityNotFound) {
		return module.ErrFormElementNotFound.Wrap(err, ht.getRawItemID(r))
	}

	return ht.errorWrapper.Wrap(err)
}
