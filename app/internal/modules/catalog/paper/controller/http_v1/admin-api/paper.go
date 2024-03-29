package http_v1

import (
	"net/http"
	module "print-shop-back/internal/modules/catalog/paper"
	view_shared "print-shop-back/internal/modules/catalog/paper/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/catalog/paper/entity/admin-api"
	usecase "print-shop-back/internal/modules/catalog/paper/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/catalog/paper/usecase/shared"
	"print-shop-back/pkg/modules/dictionaries"
	"print-shop-back/pkg/shared/view"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	paperListURL             = "/v1/catalog/papers"
	paperItemURL             = "/v1/catalog/papers/:id"
	paperItemChangeStatusURL = "/v1/catalog/papers/:id/status"
)

type (
	Paper struct {
		parser     view_shared.RequestParser
		sender     mrserver.ResponseSender
		useCase    usecase.PaperUseCase
		listSorter mrview.ListSorter
	}
)

func NewPaper(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.PaperUseCase,
	listSorter mrview.ListSorter,
) *Paper {
	return &Paper{
		parser:     parser,
		sender:     sender,
		useCase:    useCase,
		listSorter: listSorter,
	}
}

func (ht *Paper) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, paperListURL, "", ht.GetList},
		{http.MethodPost, paperListURL, "", ht.Create},

		{http.MethodGet, paperItemURL, "", ht.Get},
		{http.MethodPut, paperItemURL, "", ht.Store},
		{http.MethodDelete, paperItemURL, "", ht.Remove},

		{http.MethodPatch, paperItemChangeStatusURL, "", ht.ChangeStatus},
	}
}

func (ht *Paper) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		PaperListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *Paper) listParams(r *http.Request) entity.PaperParams {
	return entity.PaperParams{
		Filter: entity.PaperListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			ColorIDs:   ht.parser.FilterKeyInt32List(r, module.ParamNameFilterCatalogPaperColorIDs),
			FactureIDs: ht.parser.FilterKeyInt32List(r, module.ParamNameFilterCatalogPaperFactureIDs),
			Length:     ht.parser.FilterRangeInt64(r, module.ParamNameFilterLengthRange),
			Width:      ht.parser.FilterRangeInt64(r, module.ParamNameFilterWidthRange),
			Density:    ht.parser.FilterRangeInt64(r, module.ParamNameFilterDensityRange),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *Paper) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *Paper) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreatePaperRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Paper{
		Article:   request.Article,
		Caption:   request.Caption,
		ColorID:   request.ColorID,
		FactureID: request.FactureID,
		Length:    request.Length,
		Width:     request.Width,
		Density:   request.Density,
		Thickness: request.Thickness,
		Sides:     request.Sides,
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

func (ht *Paper) Store(w http.ResponseWriter, r *http.Request) error {
	request := StorePaperRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Paper{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Article:    request.Article,
		Caption:    request.Caption,
		ColorID:    request.ColorID,
		FactureID:  request.FactureID,
		Length:     request.Length,
		Width:      request.Width,
		Density:    request.Density,
		Thickness:  request.Thickness,
		Sides:      request.Sides,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Paper) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.Paper{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Paper) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *Paper) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *Paper) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *Paper) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrPaperNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if usecase_shared.FactoryErrPaperArticleAlreadyExists.Is(err) {
		return mrerr.NewCustomError("article", err)
	}

	if dictionaries.FactoryErrPaperColorNotFound.Is(err) {
		return mrerr.NewCustomError("colorId", err)
	}

	if dictionaries.FactoryErrPaperFactureNotFound.Is(err) {
		return mrerr.NewCustomError("factureId", err)
	}

	return err
}
