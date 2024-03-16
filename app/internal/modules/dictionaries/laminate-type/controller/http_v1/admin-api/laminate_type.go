package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/dictionaries/laminate-type"
	view_shared "print-shop-back/internal/modules/dictionaries/laminate-type/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/admin-api"
	usecase "print-shop-back/internal/modules/dictionaries/laminate-type/usecase/admin-api"
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
	laminateTypeListURL             = "/v1/dictionaries/laminate-types"
	laminateTypeItemURL             = "/v1/dictionaries/laminate-types/:id"
	laminateTypeItemChangeStatusURL = "/v1/dictionaries/laminate-types/:id/status"
)

type (
	LaminateType struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.LaminateTypeUseCase
		listSorter mrview.ListSorter
	}
)

func NewLaminateType(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.LaminateTypeUseCase,
	listSorter mrview.ListSorter,
) *LaminateType {
	return &LaminateType{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

func (ht *LaminateType) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, laminateTypeListURL, "", ht.GetList},
		{http.MethodPost, laminateTypeListURL, "", ht.Create},

		{http.MethodGet, laminateTypeItemURL, "", ht.Get},
		{http.MethodPut, laminateTypeItemURL, "", ht.Store},
		{http.MethodDelete, laminateTypeItemURL, "", ht.Remove},

		{http.MethodPut, laminateTypeItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *LaminateType) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		LaminateTypeListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *LaminateType) listParams(r *http.Request) entity.LaminateTypeParams {
	return entity.LaminateTypeParams{
		Filter: entity.LaminateTypeListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *LaminateType) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *LaminateType) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateLaminateTypeRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.LaminateType{
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
					"msgDictionariesLaminateTypeSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *LaminateType) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreLaminateTypeRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.LaminateType{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Caption:    request.Caption,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *LaminateType) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.LaminateType{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *LaminateType) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *LaminateType) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *LaminateType) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *LaminateType) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return dictionaries.FactoryErrLaminateTypeNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
