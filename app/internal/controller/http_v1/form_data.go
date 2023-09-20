package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/view"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    formDataListURL = "/v1/forms"
    formDataItemURL = "/v1/forms/:fid"
    formDataChangeStatusURL = "/v1/forms/:fid/status"
    formDataCompileURL = "/v1/forms/:fid/compile"
)

type FormData struct {
    service usecase.FormDataService
    serviceUIFormData usecase.UIFormDataService
}

func NewFormData(service usecase.FormDataService,
                 serviceUIFormData usecase.UIFormDataService) *FormData {
    return &FormData{
        service: service,
        serviceUIFormData: serviceUIFormData,
    }
}

func (ht *FormData) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, formDataListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, formDataListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, formDataItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, formDataItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, formDataItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, formDataChangeStatusURL, ht.ChangeStatus())

    router.HttpHandlerFunc(http.MethodPatch, formDataCompileURL, ht.Compile())
}

func (ht *FormData) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *FormData) newListFilter(c mrcore.ClientData) *entity.FormDataListFilter {
    var listFilter entity.FormDataListFilter

    parseFilterDetailing(c, &listFilter.Detailing)
    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *FormData) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *FormData) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateFormDataRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
           return err
        }

        item := entity.FormData{
            ParamName: request.ParamName,
            Caption: request.Caption,
            Detailing: request.Detailing,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := mrview.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgFormDataSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *FormData) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreFormDataRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormData{
            Id:        ht.getItemId(c),
            Version:   request.Version,
            ParamName: request.ParamName,
            Caption:   request.Caption,
            Detailing: request.Detailing,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormData) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := mrcom.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormData{
            Id:      ht.getItemId(c),
            Version: request.Version,
            Status:  request.Status,
        }

        err := ht.service.ChangeStatus(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormData) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormData) Compile() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.serviceUIFormData.CompileForm(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *FormData) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("fid"))

    if id > 0 {
        return id
    }

    return 0
}
