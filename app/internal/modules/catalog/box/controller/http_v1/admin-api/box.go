package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/catalog/box"
	view_shared "print-shop-back/internal/modules/catalog/box/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/catalog/box/entity/admin-api"
	usecase "print-shop-back/internal/modules/catalog/box/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/box/usecase/shared"
	"print-shop-back/pkg/shared/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	boxListURL             = "/v1/catalog/boxes"
	boxItemURL             = "/v1/catalog/boxes/:id"
	boxItemChangeStatusURL = "/v1/catalog/boxes/:id/status"
)

type (
	Box struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.BoxUseCase
		listSorter mrview.ListSorter
	}
)

func NewBox(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.BoxUseCase,
	listSorter mrview.ListSorter,
) *Box {
	return &Box{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

func (ht *Box) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, boxListURL, "", ht.GetList},
		{http.MethodPost, boxListURL, "", ht.Create},

		{http.MethodGet, boxItemURL, "", ht.Get},
		{http.MethodPut, boxItemURL, "", ht.Store},
		{http.MethodDelete, boxItemURL, "", ht.Remove},

		{http.MethodPatch, boxItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *Box) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		BoxListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Box) listParams(r *http.Request) entity.BoxParams {
	return entity.BoxParams{
		Filter: entity.BoxListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Length:     ht.parser.FilterRangeInt64(r, module.ParamNameFilterLengthRange),
			Width:      ht.parser.FilterRangeInt64(r, module.ParamNameFilterWidthRange),
			Depth:      ht.parser.FilterRangeInt64(r, module.ParamNameFilterDepthRange),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *Box) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *Box) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateBoxRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Box{
		Article: request.Article,
		Caption: request.Caption,
		Length:  request.Length,
		Width:   request.Width,
		Depth:   request.Depth,
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

func (ht *Box) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreBoxRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Box{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Article:    request.Article,
		Caption:    request.Caption,
		Length:     request.Length,
		Width:      request.Width,
		Depth:      request.Depth,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Box) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Box{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Box) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Box) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *Box) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Box) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrBoxNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if usecase_shared.FactoryErrBoxArticleAlreadyExists.Is(err) {
		return mrerr.NewCustomError("article", err)
	}

	return err
}
