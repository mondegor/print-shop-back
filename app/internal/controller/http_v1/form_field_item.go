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
    formFieldItemListURL = "/v1/forms/:fid/fields"
    formFieldItemItemURL = "/v1/forms/:fid/fields/:id"
    formFieldItemMoveURL = "/v1/forms/:fid/fields/:id/move"
)

type FormFieldItem struct {
    service usecase.FormFieldItemService
    serviceFormData usecase.FormDataService
    serviceFormFieldTemplate usecase.FormFieldTemplateService
}

func NewFormFieldItem(service usecase.FormFieldItemService,
                      serviceFormData usecase.FormDataService,
                      serviceFormFieldTemplate usecase.FormFieldTemplateService) *FormFieldItem {
    return &FormFieldItem{
        service: service,
        serviceFormData: serviceFormData,
        serviceFormFieldTemplate: serviceFormFieldTemplate,
    }
}

func (ht *FormFieldItem) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, formFieldItemListURL, ht.FormDataMiddleware(ht.GetList()))
    router.HttpHandlerFunc(http.MethodPost, formFieldItemListURL, ht.FormDataMiddleware(ht.Create()))

    router.HttpHandlerFunc(http.MethodGet, formFieldItemItemURL, ht.FormDataMiddleware(ht.Get()))
    router.HttpHandlerFunc(http.MethodPut, formFieldItemItemURL, ht.FormDataMiddleware(ht.Store()))
    router.HttpHandlerFunc(http.MethodDelete, formFieldItemItemURL, ht.FormDataMiddleware(ht.Remove()))

    router.HttpHandlerFunc(http.MethodPatch, formFieldItemMoveURL, ht.FormDataMiddleware(ht.Move()))
}

func (ht *FormFieldItem) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *FormFieldItem) newListFilter(c mrapp.ClientData) *entity.FormFieldItemListFilter {
    var listFilter entity.FormFieldItemListFilter

    listFilter.FormId = ht.getFormId(c)
    parseFilterDetailing(c, &listFilter.Detailing)

    return &listFilter
}

func (ht *FormFieldItem) Get() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c), ht.getFormId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *FormFieldItem) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateFormFieldItem{}

        if err := c.ParseAndValidate(&request); err != nil {
           return err
        }

        item := entity.FormFieldItem{
            FormId: ht.getFormId(c),
            TemplateId: request.TemplateId,
            ParamName: request.ParamName,
            Caption: request.Caption,
            Required: request.Required,
        }

        err := ht.service.Create(c.Context(), &item)

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

func (ht *FormFieldItem) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreFormFieldItem{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormFieldItem{
            Id:        ht.getItemId(c),
            FormId:    ht.getFormId(c),
            Version:   request.Version,
            ParamName: request.ParamName,
            Caption:   request.Caption,
            Required:  request.Required,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            if usecase.ErrFormFieldItemParamNameAlreadyExists.Is(err) {
                return mrlib.NewUserErrorListWithError("paramName", err)
            }

            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormFieldItem) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c), ht.getFormId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormFieldItem) Move() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.MoveFormFieldItem{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        err := ht.service.MoveAfterId(
            c.Context(),
            ht.getItemId(c),
            request.AfterNodeId,
            ht.getFormId(c),
        )

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *FormFieldItem) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
