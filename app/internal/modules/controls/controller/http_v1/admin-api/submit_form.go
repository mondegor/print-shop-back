package http_v1

import (
	"fmt"
	"net/http"
	module "print-shop-back/internal/modules/controls"
	view_shared "print-shop-back/internal/modules/controls/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/controls/usecase/shared"
	"strconv"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	submitFormListURL             = "/v1/controls/submit-forms"
	submitFormItemURL             = "/v1/controls/submit-forms/:id"
	submitFormItemChangeStatusURL = "/v1/controls/submit-forms/:id/status"
	submitFormItemCompileURL      = "/v1/controls/submit-forms/:id/compile"
)

type (
	SubmitForm struct {
		parser  view_shared.RequestParser
		sender  mrserver.ResponseSender
		useCase usecase.SubmitFormUseCase
		// useCaseUISubmitForm usecase.UISubmitFormUseCase
		listSorter mrview.ListSorter
	}
)

func NewSubmitForm(
	parser view_shared.RequestParser,
	sender mrserver.ResponseSender,
	useCase usecase.SubmitFormUseCase,
	// useCaseUISubmitForm usecase.UISubmitFormUseCase,
	listSorter mrview.ListSorter,
) *SubmitForm {
	return &SubmitForm{
		parser:  parser,
		sender:  sender,
		useCase: useCase,
		// useCaseUISubmitForm: useCaseUISubmitForm,
		listSorter: listSorter,
	}
}

func (ht *SubmitForm) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, submitFormListURL, "", ht.GetList},
		{http.MethodPost, submitFormListURL, "", ht.Create},

		{http.MethodGet, submitFormItemURL, "", ht.Get},
		{http.MethodPut, submitFormItemURL, "", ht.Store},
		{http.MethodDelete, submitFormItemURL, "", ht.Remove},

		{http.MethodPut, submitFormItemChangeStatusURL, "", ht.ChangeStatus},

		{http.MethodPatch, submitFormItemCompileURL, "", ht.Compile},
	}
}

func (ht *SubmitForm) GetList(w http.ResponseWriter, r *http.Request) error {
	items, totalItems, err := ht.useCase.GetList(r.Context(), ht.listParams(r))

	if err != nil {
		return err
	}

	return ht.sender.Send(
		w,
		http.StatusOK,
		SubmitFormListResponse{
			Items: items,
			Total: totalItems,
		},
	)
}

func (ht *SubmitForm) listParams(r *http.Request) entity.SubmitFormParams {
	return entity.SubmitFormParams{
		Filter: entity.SubmitFormListFilter{
			SearchText: ht.parser.FilterString(r, module.ParamNameFilterSearchText),
			Detailing:  ht.parser.FilterElementDetailingList(r, module.ParamNameFilterElementDetailing),
			Statuses:   ht.parser.FilterStatusList(r, module.ParamNameFilterStatuses),
		},
		Sorter: ht.parser.SortParams(r, ht.listSorter),
		Pager:  ht.parser.PageParams(r),
	}
}

func (ht *SubmitForm) Get(w http.ResponseWriter, r *http.Request) error {
	item, err := ht.useCase.GetItem(r.Context(), ht.getItemID(r))

	if err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *SubmitForm) Create(w http.ResponseWriter, r *http.Request) error {
	request := CreateSubmitFormRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.SubmitForm{
		ParamName: request.ParamName,
		Caption:   request.Caption,
		Detailing: request.Detailing,
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
					"msgControlsSubmitFormSuccessCreated",
					"entity has been success created",
				),
			},
		)
	}
}

func (ht *SubmitForm) Store(w http.ResponseWriter, r *http.Request) error {
	request := StoreSubmitFormRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.SubmitForm{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		ParamName:  request.ParamName,
		Caption:    request.Caption,
		Detailing:  request.Detailing,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *SubmitForm) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := ChangeItemStatusRequest{}

	if err := ht.parser.Validate(r, &request); err != nil {
		return err
	}

	item := entity.SubmitForm{
		ID:         ht.getItemID(r),
		TagVersion: request.TagVersion,
		Status:     request.Status,
	}

	if err := ht.useCase.ChangeStatus(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *SubmitForm) Remove(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCase.Remove(r.Context(), ht.getItemID(r)); err != nil {
		return err
	}

	return ht.sender.SendNoContent(w)
}

func (ht *SubmitForm) Compile(w http.ResponseWriter, r *http.Request) error {
	// item, err := ht.useCaseUISubmitForm.CompileForm(r.Context(), ht.getItemID(r))
	err := fmt.Errorf("of")
	item := ""

	if err != nil {
		return err
	}

	return ht.sender.Send(w, http.StatusOK, item)
}

func (ht *SubmitForm) getItemID(r *http.Request) mrtype.KeyInt32 {
	return ht.parser.PathKeyInt32(r, "id")
}

func (ht *SubmitForm) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *SubmitForm) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrSubmitFormNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("version", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	return err
}
