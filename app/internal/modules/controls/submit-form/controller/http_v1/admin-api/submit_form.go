package http_v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	module "print-shop-back/internal/modules/controls/submit-form"
	view_shared "print-shop-back/internal/modules/controls/submit-form/controller/http_v1/shared/view"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	usecase "print-shop-back/internal/modules/controls/submit-form/usecase/admin-api"
	usecase_shared "print-shop-back/internal/modules/controls/submit-form/usecase/shared"
	"print-shop-back/pkg/shared/view"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

const (
	submitFormListURL             = "/v1/controls/submit-forms"
	submitFormItemURL             = "/v1/controls/submit-forms/:id"
	submitFormItemChangeStatusURL = "/v1/controls/submit-forms/:id/status"

	submitFormVersionItemJsonURL    = "/v1/controls/submit-forms/:id/versions/:version/json"
	submitFormItemPrepareForTestURL = "/v1/controls/submit-forms/:id/prepare-for-test"
	submitFormItemPublishURL        = "/v1/controls/submit-forms/:id/publish"
)

type (
	SubmitForm struct {
		parser         view_shared.RequestParser
		sender         mrserver.FileResponseSender
		useCase        usecase.SubmitFormUseCase
		useCaseVersion usecase.FormVersionUseCase
		listSorter     mrview.ListSorter
	}
)

func NewSubmitForm(
	parser view_shared.RequestParser,
	sender mrserver.FileResponseSender,
	useCase usecase.SubmitFormUseCase,
	useCaseVersion usecase.FormVersionUseCase,
	listSorter mrview.ListSorter,
) *SubmitForm {
	return &SubmitForm{
		parser:         parser,
		sender:         sender,
		useCase:        useCase,
		useCaseVersion: useCaseVersion,
		listSorter:     listSorter,
	}
}

func (ht *SubmitForm) Handlers() []mrserver.HttpHandler {
	return []mrserver.HttpHandler{
		{http.MethodGet, submitFormListURL, "", ht.GetList},
		{http.MethodPost, submitFormListURL, "", ht.Create},

		{http.MethodGet, submitFormItemURL, "", ht.Get},
		{http.MethodPatch, submitFormItemURL, "", ht.Store},
		{http.MethodDelete, submitFormItemURL, "", ht.Remove},

		{http.MethodPatch, submitFormItemChangeStatusURL, "", ht.ChangeStatus},

		{http.MethodGet, submitFormVersionItemJsonURL, "", ht.GetVersionJson},
		{http.MethodPatch, submitFormItemPrepareForTestURL, "", ht.PrepareForTest},
		{http.MethodPatch, submitFormItemPublishURL, "modControlsSubmitFormToPublish", ht.Publish},
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
		RewriteName: request.RewriteName,
		ParamName:   request.ParamName,
		Caption:     request.Caption,
		Detailing:   request.Detailing,
	}

	if itemID, err := ht.useCase.Create(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	} else {
		return ht.sender.Send(
			w,
			http.StatusCreated,
			view.SuccessCreatedItemResponse{
				ItemID: itemID.String(),
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
		ID:          ht.getItemID(r),
		TagVersion:  request.TagVersion,
		RewriteName: request.RewriteName,
		ParamName:   request.ParamName,
		Caption:     request.Caption,
	}

	if err := ht.useCase.Store(r.Context(), item); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *SubmitForm) ChangeStatus(w http.ResponseWriter, r *http.Request) error {
	request := view.ChangeItemStatusRequest{}

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

func (ht *SubmitForm) GetVersionJson(w http.ResponseWriter, r *http.Request) error {
	primary := entity.PrimaryKey{
		FormID:  ht.getItemID(r),
		Version: ht.getItemVersion(r),
	}
	body, err := ht.useCaseVersion.GetItemJson(r.Context(), primary, true)

	if err != nil {
		return err
	}

	return ht.sender.SendAttachmentFile(
		r.Context(),
		w,
		mrtype.File{
			FileInfo: mrtype.FileInfo{
				ContentType:  "application/json",
				OriginalName: fmt.Sprintf(module.JsonFileNamePattern, primary.FormID, primary.Version),
				Size:         int64(len(body)),
			},
			Body: io.NopCloser(bytes.NewReader(body)),
		},
	)
}

func (ht *SubmitForm) PrepareForTest(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCaseVersion.PrepareForTest(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *SubmitForm) Publish(w http.ResponseWriter, r *http.Request) error {
	if err := ht.useCaseVersion.Publish(r.Context(), ht.getItemID(r)); err != nil {
		return ht.wrapError(err, r)
	}

	return ht.sender.SendNoContent(w)
}

func (ht *SubmitForm) getItemID(r *http.Request) uuid.UUID {
	return ht.parser.PathParamUUID(r, "id")
}

func (ht *SubmitForm) getItemVersion(r *http.Request) int32 {
	return int32(ht.parser.PathParamInt64(r, "version"))
}

func (ht *SubmitForm) getRawItemID(r *http.Request) string {
	return ht.parser.PathParamString(r, "id")
}

func (ht *SubmitForm) wrapError(err error, r *http.Request) error {
	if mrcore.FactoryErrUseCaseEntityNotFound.Is(err) {
		return usecase_shared.FactoryErrSubmitFormNotFound.Wrap(err, ht.getRawItemID(r))
	}

	if mrcore.FactoryErrUseCaseEntityVersionInvalid.Is(err) {
		return mrerr.NewCustomError("tagVersion", err)
	}

	if mrcore.FactoryErrUseCaseSwitchStatusRejected.Is(err) {
		return mrerr.NewCustomError("status", err)
	}

	if usecase_shared.FactoryErrSubmitFormRewriteNameAlreadyExists.Is(err) {
		return mrerr.NewCustomError("rewriteName", err)
	}

	if usecase_shared.FactoryErrSubmitFormParamNameAlreadyExists.Is(err) {
		return mrerr.NewCustomError("paramName", err)
	}

	return err
}
