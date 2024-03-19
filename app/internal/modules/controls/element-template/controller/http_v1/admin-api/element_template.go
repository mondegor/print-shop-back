package http_v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	module "print-shop-back/internal/modules/controls/element-template"
	view_shared "print-shop-back/internal/modules/controls/element-template/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/controls/element-template/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/element-template/usecase/admin-api"
	"print-shop-back/pkg/modules/controls"
	"print-shop-back/pkg/shared/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrparser"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	elementListTemplateURL             = "/v1/controls/element-templates"
	elementTemplateItemURL             = "/v1/controls/element-templates/:id"
	elementTemplateItemChangeStatusURL = "/v1/controls/element-templates/:id/status"
	elementTemplateItemJsonURL         = "/v1/controls/element-templates/:id/json"
)

type (
	ElementTemplate struct {
		parser     view_shared.RequestParser
		sender     mrserver.FileResponseSender
		useCase    usecase.ElementTemplateUseCase
		listSorter mrview.ListSorter
	}
)

func NewElementTemplate(
	parser view_shared.RequestParser,
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

func (ht *ElementTemplate) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, elementListTemplateURL, "", ht.GetList},
		{http.MethodPost, elementListTemplateURL, "", ht.Create},

		{http.MethodGet, elementTemplateItemURL, "", ht.Get},
		{http.MethodPatch, elementTemplateItemURL, "", ht.Store},
		{http.MethodDelete, elementTemplateItemURL, "", ht.Remove},

		{http.MethodPatch, elementTemplateItemChangeStatusURL, "", ht.ChangeStatus},

		{http.MethodGet, elementTemplateItemJsonURL, "", ht.GetJson},
	}
}

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

func (ht *ElementTemplate) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

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

func (ht *ElementTemplate) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreElementTemplateRequest{}
	rawElementTemplate := []byte(r.FormValue(module.ParamNameElementTemplateObject))

	if err := ht.parser.ValidateContent(r.Context(), rawElementTemplate, &request); err != nil {
		return err
	}

	file, err := ht.parser.FormFileContent(r, module.ParamNameElementTemplateAttachment)

	if err != nil {
		// указывать файл необязательно
		if !mrcore.FactoryErrHttpFileUpload.Is(err) {
			return mrparser.WrapFileError(err, module.ParamNameElementTemplateAttachment)
		}

		file.Body = []byte{}
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
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return controls.FactoryErrElementTemplateNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
