package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/dto"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrentity"
)

const (
    formDataGetListURL = "/v1/forms"
    formDataGetItemURL = "/v1/forms/:fid"
    formDataCreateURL = "/v1/forms"
    formDataStoreURL = "/v1/forms/:fid"
    formDataChangeStatusURL = "/v1/forms/:fid/status/"
    formDataRemoveURL = "/v1/forms/:fid"
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

func (ht *FormData) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, formDataGetListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodGet, formDataGetItemURL, ht.GetItem())
    router.HttpHandlerFunc(http.MethodPost, formDataCreateURL, ht.Create())
    router.HttpHandlerFunc(http.MethodPut, formDataStoreURL, ht.Store())
    router.HttpHandlerFunc(http.MethodPut, formDataChangeStatusURL, ht.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, formDataRemoveURL, ht.Remove())
    router.HttpHandlerFunc(http.MethodPatch, formDataCompileURL, ht.Compile())
}

func (ht *FormData) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *FormData) newListFilter(c mrapp.ClientData) *entity.FormDataListFilter {
    var listFilter entity.FormDataListFilter

    parseFilterDetailing(c, &listFilter.Detailing)
    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *FormData) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *FormData) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateFormData{}

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

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgFormDataSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *FormData) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreFormData{}

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

func (ht *FormData) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

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

func (ht *FormData) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormData) Compile() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.serviceUIFormData.CompileForm(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *FormData) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("fid"))

    if id > 0 {
        return id
    }

    return 0
}
