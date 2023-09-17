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
    formFieldTemplateListURL = "/v1/form-field-templates"
    formFieldTemplateItemURL = "/v1/form-field-templates/:id"
    formFieldTemplateChangeStatusURL = "/v1/form-field-templates/:id/status"
)

type FormFieldTemplate struct {
    service usecase.FormFieldTemplateService
}

func NewFormFieldTemplate(service usecase.FormFieldTemplateService) *FormFieldTemplate {
    return &FormFieldTemplate{
        service: service,
    }
}

func (ht *FormFieldTemplate) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, formFieldTemplateListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, formFieldTemplateListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, formFieldTemplateItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, formFieldTemplateItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, formFieldTemplateItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, formFieldTemplateChangeStatusURL, ht.ChangeStatus())
}

func (ht *FormFieldTemplate) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *FormFieldTemplate) newListFilter(c mrcore.ClientData) *entity.FormFieldTemplateListFilter {
    var listFilter entity.FormFieldTemplateListFilter

    parseFilterDetailing(c, &listFilter.Detailing)
    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *FormFieldTemplate) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *FormFieldTemplate) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateFormFieldTemplate{}

        if err := c.ParseAndValidate(&request); err != nil {
           return err
        }

        item := entity.FormFieldTemplate{
            ParamName: request.ParamName,
            Caption: request.Caption,
            Type: request.Type,
            Detailing: request.Detailing,
            Body: request.Body,
        }

        err := ht.service.Create(c.Context(), &item)

        if err != nil {
            return err
        }

        response := mrview.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgFormFieldTemplateSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *FormFieldTemplate) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreFormFieldTemplate{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormFieldTemplate{
            Id:        ht.getItemId(c),
            Version:   request.Version,
            ParamName: request.ParamName,
            Caption:   request.Caption,
            Type:      request.Type,
            Detailing: request.Detailing,
            Body:      request.Body,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormFieldTemplate) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := mrcom.ChangeItemStatusRequest{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormFieldTemplate{
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

func (ht *FormFieldTemplate) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormFieldTemplate) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
