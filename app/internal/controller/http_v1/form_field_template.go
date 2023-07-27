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
    formFieldTemplateGetListURL = "/v1/form-field-templates"
    formFieldTemplateGetItemURL = "/v1/form-field-templates/:id"
    formFieldTemplateCreateURL = "/v1/form-field-templates"
    formFieldTemplateStoreURL = "/v1/form-field-templates/:id"
    formFieldTemplateChangeStatusURL = "/v1/form-field-templates/:id/status"
    formFieldTemplateRemoveURL = "/v1/form-field-templates/:id"
)

type FormFieldTemplate struct {
    service usecase.FormFieldTemplateService
}

func NewFormFieldTemplate(service usecase.FormFieldTemplateService) *FormFieldTemplate {
    return &FormFieldTemplate{
        service: service,
    }
}

func (ht *FormFieldTemplate) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, formFieldTemplateGetListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodGet, formFieldTemplateGetItemURL, ht.GetItem())
    router.HttpHandlerFunc(http.MethodPost, formFieldTemplateCreateURL, ht.Create())
    router.HttpHandlerFunc(http.MethodPut, formFieldTemplateStoreURL, ht.Store())
    router.HttpHandlerFunc(http.MethodPut, formFieldTemplateChangeStatusURL, ht.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, formFieldTemplateRemoveURL, ht.Remove())
}

func (ht *FormFieldTemplate) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *FormFieldTemplate) newListFilter(c mrapp.ClientData) *entity.FormFieldTemplateListFilter {
    var listFilter entity.FormFieldTemplateListFilter

    parseFilterDetailing(c, &listFilter.Detailing)
    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *FormFieldTemplate) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *FormFieldTemplate) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateFormFieldTemplate{}

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

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgFormFieldTemplateSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *FormFieldTemplate) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreFormFieldTemplate{}

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

func (ht *FormFieldTemplate) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

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

func (ht *FormFieldTemplate) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormFieldTemplate) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
