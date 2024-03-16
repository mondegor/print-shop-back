package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/dictionaries/print-format"
	view_shared "print-shop-back/internal/modules/dictionaries/print-format/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/dictionaries/print-format/entity/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/print-format/usecase/admin-api"
	"print-shop-back/pkg/modules/dictionaries"
	"strconv"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	printFormatListURL             = "/v1/dictionaries/print-formats"
	printFormatItemURL             = "/v1/dictionaries/print-formats/:id"
	printFormatItemChangeStatusURL = "/v1/dictionaries/print-formats/:id/status"
)

type (
	PrintFormat struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.PrintFormatUseCase
		listSorter mrview.ListSorter
	}
)

func NewPrintFormat(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.PrintFormatUseCase,
	listSorter mrview.ListSorter,
) *PrintFormat {
	return &PrintFormat{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

func (ht *PrintFormat) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, printFormatListURL, "", ht.GetList},
		{http.MethodPost, printFormatListURL, "", ht.Create},

		{http.MethodGet, printFormatItemURL, "", ht.Get},
		{http.MethodPut, printFormatItemURL, "", ht.Store},
		{http.MethodDelete, printFormatItemURL, "", ht.Remove},

		{http.MethodPut, printFormatItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *PrintFormat) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		PrintFormatListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *PrintFormat) listParams(r *http.Request) entity.PrintFormatParams {
	return entity.PrintFormatParams{
		Filter: entity.PrintFormatListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Length:     ht.parser.FilterRangeInt64(r, module.ParamNameFilterLengthRange),
			Width:      ht.parser.FilterRangeInt64(r, module.ParamNameFilterWidthRange),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *PrintFormat) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *PrintFormat) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreatePrintFormatRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PrintFormat{
		Caption: request.Caption,
		Length:  request.Length,
		Width:   request.Width,
	}

	if itemID, err := ht.useCase.Create(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	} else {
		return ht.sender.Send(
			w,
			http.StatusCreated,
			SuccessCreatedItemResponse{
				ItemID: strconv.Itoa(int(itemID)),
				Message: mrlang.Ctx(r.Context()).TranslateMessage(
					"msgDictionariesPrintFormatSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *PrintFormat) Store(w http.ResponseWriter, r *http.Request) error {
	request := StorePrintFormatRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PrintFormat{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Caption:    request.Caption,
		Length:     request.Length,
		Width:      request.Width,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PrintFormat) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PrintFormat{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PrintFormat) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PrintFormat) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *PrintFormat) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *PrintFormat) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return dictionaries.FactoryErrPrintFormatNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
