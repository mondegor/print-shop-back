package http_v1

import (
    "calc-user-data-back-adm/internal/controller/dto"
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/internal/usecase"
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrentity"
    "fmt"
    "net/http"
)

const (
    formFieldTemplateGetListURL = "/v1/form-field-templates"
    formFieldTemplateGetItemURL = "/v1/form-field-templates/:id"
    formFieldTemplateCreateURL = "/v1/form-field-templates"
    formFieldTemplateStoreURL = "/v1/form-field-templates/:id"
    formFieldTemplateChangeStatusURL = "/v1/form-field-templates/:id/status"
    formFieldTemplateRemove = "/v1/form-field-templates/:id"
)

type FormFieldTemplate struct {
    service usecase.FormFieldTemplateService
}

func NewFormFieldTemplate(service usecase.FormFieldTemplateService) *FormFieldTemplate {
    return &FormFieldTemplate{
        service: service,
    }
}

func (f *FormFieldTemplate) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, formFieldTemplateGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, formFieldTemplateGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, formFieldTemplateCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, formFieldTemplateStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, formFieldTemplateChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, formFieldTemplateRemove, f.Remove())
}

func (f *FormFieldTemplate) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *FormFieldTemplate) newListFilter(c mrapp.ClientData) *entity.FormFieldTemplateListFilter {
    var listFilter entity.FormFieldTemplateListFilter

    parseFilterDetailing(c, &listFilter.Detailing)
    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *FormFieldTemplate) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *FormFieldTemplate) Create() mrapp.HttpHandlerFunc {
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

        err := f.service.Create(c.Context(), &item)

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

func (f *FormFieldTemplate) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreFormFieldTemplate{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormFieldTemplate{
            Id: f.getItemId(c),
            Version: request.Version,
            ParamName: request.ParamName,
            Caption: request.Caption,
            Type: request.Type,
            Detailing: request.Detailing,
            Body: request.Body,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *FormFieldTemplate) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormFieldTemplate{
            Id: f.getItemId(c),
            Version: request.Version,
            Status: request.Status,
        }

        err := f.service.ChangeStatus(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *FormFieldTemplate) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *FormFieldTemplate) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
