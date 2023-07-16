package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/dto"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrlib"
)

const (
    formFieldItemGetListURL = "/v1/forms/:id/fields"
    formFieldItemGetItemURL = "/v1/forms/:id/fields/:fid"
    formFieldItemCreateURL = "/v1/forms/:id/fields"
    formFieldItemStoreURL = "/v1/forms/:id/fields/:fid"
    formFieldItemMoveURL = "/v1/forms/:id/fields/:fid/move"
    formFieldItemRemove = "/v1/forms/:id/fields/:fid"
)

type FormFieldItem struct {
    service usecase.FormFieldItemService
    serviceFormFieldItemOrderer usecase.FormFieldItemOrdererService
    serviceFormData usecase.FormDataService
    serviceFormFieldTemplate usecase.FormFieldTemplateService
}

func NewFormFieldItem(service usecase.FormFieldItemService,
                      serviceFormFieldItemOrderer usecase.FormFieldItemOrdererService,
                      serviceFormData usecase.FormDataService,
                      serviceFormFieldTemplate usecase.FormFieldTemplateService) *FormFieldItem {
    return &FormFieldItem{
        service: service,
        serviceFormFieldItemOrderer: serviceFormFieldItemOrderer,
        serviceFormData: serviceFormData,
        serviceFormFieldTemplate: serviceFormFieldTemplate,
    }
}

func (f *FormFieldItem) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, formFieldItemGetListURL, f.FormDataMiddleware(f.GetList()))
    router.HttpHandlerFunc(http.MethodGet, formFieldItemGetItemURL, f.FormDataMiddleware(f.GetItem()))
    router.HttpHandlerFunc(http.MethodPost, formFieldItemCreateURL, f.FormDataMiddleware(f.Create()))
    router.HttpHandlerFunc(http.MethodPut, formFieldItemStoreURL, f.FormDataMiddleware(f.Store()))
    router.HttpHandlerFunc(http.MethodPatch, formFieldItemMoveURL, f.FormDataMiddleware(f.Move()))
    router.HttpHandlerFunc(http.MethodDelete, formFieldItemRemove, f.FormDataMiddleware(f.Remove()))
}

func (f *FormFieldItem) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *FormFieldItem) newListFilter(c mrapp.ClientData) *entity.FormFieldItemListFilter {
    var listFilter entity.FormFieldItemListFilter

    listFilter.FormId = f.getFormId(c)
    parseFilterDetailing(c, &listFilter.Detailing)

    return &listFilter
}

func (f *FormFieldItem) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c), f.getFormId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *FormFieldItem) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateFormFieldItem{}

        if err := c.ParseAndValidate(&request); err != nil {
           return err
        }

        item := entity.FormFieldItem{
            FormId: f.getFormId(c),
            TemplateId: request.TemplateId,
            ParamName: request.ParamName,
            Caption: request.Caption,
            Required: request.Required,
        }

        err := f.service.Create(c.Context(), &item)

        if err != nil {
            if usecase.ErrFormFieldItemTemplateNotFound.Is(err) {
                return mrlib.NewUserErrorListWithError("templateId", err)
            }

            if usecase.ErrFormFieldItemParamNameAlreadyExists.Is(err) {
                return mrlib.NewUserErrorListWithError("paramName", err)
            }

            return err
        }

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgFormFieldItemSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (f *FormFieldItem) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreFormFieldItem{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormFieldItem{
            Id: f.getItemId(c),
            FormId: f.getFormId(c),
            Version: request.Version,
            ParamName: request.ParamName,
            Caption: request.Caption,
            Required: request.Required,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            if usecase.ErrFormFieldItemParamNameAlreadyExists.Is(err) {
                return mrlib.NewUserErrorListWithError("paramName", err)
            }

            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *FormFieldItem) Move() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.MoveFormFieldItem{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        err := f.serviceFormFieldItemOrderer.MoveAfterId(
            c.Context(),
            f.getItemId(c),
            request.AfterNodeId,
        )

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *FormFieldItem) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c), f.getFormId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *FormFieldItem) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("fid"))

    if id > 0 {
        return id
    }

    return 0
}
