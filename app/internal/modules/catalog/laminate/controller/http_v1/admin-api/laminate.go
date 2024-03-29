package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/catalog/laminate"
	view_shared "print-shop-back/internal/modules/catalog/laminate/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/catalog/laminate/entity/admin-api"
	usecase "print-shop-back/internal/modules/catalog/laminate/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/laminate/usecase/shared"
	"print-shop-back/pkg/modules/dictionaries"
	"print-shop-back/pkg/shared/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	laminateListURL             = "/v1/catalog/laminates"
	laminateItemURL             = "/v1/catalog/laminates/:id"
	laminateItemChangeStatusURL = "/v1/catalog/laminates/:id/status"
)

type (
	Laminate struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.LaminateUseCase
		listSorter mrview.ListSorter
	}
)

func NewLaminate(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.LaminateUseCase,
	listSorter mrview.ListSorter,
) *Laminate {
	return &Laminate{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

func (ht *Laminate) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, laminateListURL, "", ht.GetList},
		{http.MethodPost, laminateListURL, "", ht.Create},

		{http.MethodGet, laminateItemURL, "", ht.Get},
		{http.MethodPut, laminateItemURL, "", ht.Store},
		{http.MethodDelete, laminateItemURL, "", ht.Remove},

		{http.MethodPatch, laminateItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *Laminate) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		LaminateListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Laminate) listParams(r *http.Request) entity.LaminateParams {
	return entity.LaminateParams{
		Filter: entity.LaminateListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			TypeIDs:    ht.parser.FilterKeyInt32List(r, module.ParamNameFilterCatalogLaminateTypeIDs),
			Length:     ht.parser.FilterRangeInt64(r, module.ParamNameFilterLengthRange),
			Weight:     ht.parser.FilterRangeInt64(r, module.ParamNameFilterWeightRange),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *Laminate) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *Laminate) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateLaminateRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Laminate{
		Article:   request.Article,
		Caption:   request.Caption,
		TypeID:    request.TypeID,
		Length:    request.Length,
		Weight:    request.Weight,
		Thickness: request.Thickness,
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

func (ht *Laminate) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreLaminateRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Laminate{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Article:    request.Article,
		Caption:    request.Caption,
		TypeID:     request.TypeID,
		Length:     request.Length,
		Weight:     request.Weight,
		Thickness:  request.Thickness,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Laminate) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Laminate{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Laminate) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Laminate) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *Laminate) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Laminate) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrLaminateNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if usecase_shared.FactoryErrLaminateArticleAlreadyExists.Is(err) {
		return mrerr.NewCustomError("article", err)
	}

	if dictionaries.FactoryErrLaminateTypeNotFound.Is(err) {
		return mrerr.NewCustomError("typeId", err)
	}

	return err
}
