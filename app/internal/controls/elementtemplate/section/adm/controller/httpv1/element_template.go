package httpv1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/module"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/usecase"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/shared/validate"
	"github.com/mondegor/print-shop-back/pkg/controls/api"
	"github.com/mondegor/print-shop-back/pkg/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	elementListTemplateURL             = "/v1/controls/element-templates"
	elementTemplateItemURL             = "/v1/controls/element-templates/{id}"
	elementTemplateItemChangeStatusURL = "/v1/controls/element-templates/{id}/status"
	elementTemplateItemJsonURL         = "/v1/controls/element-templates/{id}/json"
)

type (
	// ElementTemplate - comment struct.
	ElementTemplate struct {
		parser     validate.RequestElementTemplateParser
		sender     mrserver.FileResponseSender
		useCase    usecase.ElementTemplateUseCase
		listSorter mrview.ListSorter
	}
)

// NewElementTemplate - создаёт контроллер ElementTemplate.
func NewElementTemplate(
	parser validate.RequestElementTemplateParser,
	sender mrserver.FileResponseSender,
	useCase usecase.ElementTemplateUseCase,
	listSorter mrview.ListSorter,
) *ElementTemplate {
	return &ElementTemplate{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

// Handlers - возвращает обработчики контроллера ElementTemplate.
func (ht *ElementTemplate) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{Method: http.MethodGet, URL: elementListTemplateURL, Func: ht.GetList},
		{Method: http.MethodPost, URL: elementListTemplateURL, Func: ht.Create},

		{Method: http.MethodGet, URL: elementTemplateItemURL, Func: ht.Get},
		{Method: http.MethodPatch, URL: elementTemplateItemURL, Func: ht.Store},
		{Method: http.MethodDelete, URL: elementTemplateItemURL, Func: ht.Remove},

		{Method: http.MethodPatch, URL: elementTemplateItemChangeStatusURL, Func: ht.ChangeStatus},

		{Method: http.MethodGet, URL: elementTemplateItemJsonURL, Func: ht.GetJson},
	}
}

// GetList - comment method.
func (ht *ElementTemplate) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))
	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		ElementTemplateListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *ElementTemplate) listParams(r *http.Request) entity.ElementTemplateParams {
	return entity.ElementTemplateParams{
		Filter: entity.ElementTemplateListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Detailing:  ht.parser.FilterElementDetailingList(r, module.ParamNameFilterElementDetailing),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

// Get - comment method.
func (ht *ElementTemplate) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))
	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

// GetJson - comment method.
func (ht *ElementTemplate) GetJson(w http.ResponseWriter, r *http.Request) error {
	itemID := ht.getItemID(r)

	body, err := ht.useCase.GetItemJson(r.Context(), itemID, true)
	if err != nil {
		return err
	}

	return ht.sender.SendAttachmentFile(
		r.Context(),
		w,
		mrtype.File{
			FileInfo: mrtype.FileInfo{
				ContentType:  "application/json",
				OriginalName: fmt.Sprintf(module.JsonFileNamePattern, itemID),
				Size:         int64(len(body)),
			},
			Body: io.NopCloser(bytes.NewReader(body)),
		},
	)
}

// Create - comment method.
func (ht *ElementTemplate) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateElementTemplateRequest{}
	rawElementTemplate := []byte(r.FormValue(module.ParamNameElementTemplateObject))

	if err := ht.parser.ValidateContent(r.Context(), rawElementTemplate, &request); err != nil {
		return err
	}

	file, err := ht.parser.FormFileContent(r, module.ParamNameElementTemplateAttachment)
	if err != nil {
		return mrparser.WrapFileError(err, module.ParamNameElementTemplateAttachment)
	}

	item := entity.ElementTemplate{
		ParamName: request.ParamName,
		Caption:   request.Caption,
		Type:      request.Type,
		Detailing: request.Detailing,
		Body:      file.Body,
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
func (ht *ElementTemplate) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreElementTemplateRequest{}
	rawElementTemplate := []byte(r.FormValue(module.ParamNameElementTemplateObject))

	if err := ht.parser.ValidateContent(r.Context(), rawElementTemplate, &request); err != nil {
		return err
	}

	file, err := ht.parser.FormFileContent(r, module.ParamNameElementTemplateAttachment)
	if err != nil {
		// указывать файл необязательно
		if !mrcore.ErrHttpFileUpload.Is(err) {
			return mrparser.WrapFileError(err, module.ParamNameElementTemplateAttachment)
		}

		file.Body = nil
	}

	item := entity.ElementTemplate{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		ParamName:  request.ParamName,
		Caption:    request.Caption,
		Body:       file.Body,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// ChangeStatus - comment method.
func (ht *ElementTemplate) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.ElementTemplate{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

// Remove - comment method.
func (ht *ElementTemplate) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *ElementTemplate) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *ElementTemplate) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *ElementTemplate) wrapError(err error, r *http.Request) error {
	if mrcore.ErrUseCaseEntityNotFound.Is(err) {
		return api.ErrElementTemplateNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.ErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.ErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
