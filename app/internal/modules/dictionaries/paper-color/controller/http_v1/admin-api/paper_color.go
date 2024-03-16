package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/dictionaries/paper-color"
	view_shared "print-shop-back/internal/modules/dictionaries/paper-color/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/paper-color/usecase/admin-api"
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
	paperColorListURL             = "/v1/dictionaries/paper-colors"
	paperColorItemURL             = "/v1/dictionaries/paper-colors/:id"
	paperColorItemChangeStatusURL = "/v1/dictionaries/paper-colors/:id/status"
)

type (
	PaperColor struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.PaperColorUseCase
		listSorter mrview.ListSorter
	}
)

func NewPaperColor(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.PaperColorUseCase,
	listSorter mrview.ListSorter,
) *PaperColor {
	return &PaperColor{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

func (ht *PaperColor) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, paperColorListURL, "", ht.GetList},
		{http.MethodPost, paperColorListURL, "", ht.Create},

		{http.MethodGet, paperColorItemURL, "", ht.Get},
		{http.MethodPut, paperColorItemURL, "", ht.Store},
		{http.MethodDelete, paperColorItemURL, "", ht.Remove},

		{http.MethodPut, paperColorItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *PaperColor) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		PaperColorListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *PaperColor) listParams(r *http.Request) entity.PaperColorParams {
	return entity.PaperColorParams{
		Filter: entity.PaperColorListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *PaperColor) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *PaperColor) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreatePaperColorRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperColor{
		Caption: request.Caption,
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
					"msgDictionariesPaperColorSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *PaperColor) Store(w http.ResponseWriter, r *http.Request) error {
	request := StorePaperColorRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperColor{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Caption:    request.Caption,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PaperColor) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.PaperColor{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PaperColor) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *PaperColor) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *PaperColor) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *PaperColor) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return dictionaries.FactoryErrPaperColorNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
