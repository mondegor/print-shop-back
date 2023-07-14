package http_v1

import (
    "print-shop-back/internal/controller/dto"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrentity"
    "fmt"
    "net/http"
)

const (
    formDataGetListURL = "/v1/forms"
    formDataGetItemURL = "/v1/forms/:id"
    formDataCreateURL = "/v1/forms"
    formDataStoreURL = "/v1/forms/:id"
    formDataChangeStatusURL = "/v1/forms/:id/status/"
    formDataRemove = "/v1/forms/:id"
)

type FormData struct {
    service usecase.FormDataService
}

func NewFormData(service usecase.FormDataService) *FormData {
    return &FormData{
        service: service,
    }
}

func (f *FormData) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, formDataGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, formDataGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, formDataCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, formDataStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, formDataChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, formDataRemove, f.Remove())
}

func (f *FormData) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *FormData) newListFilter(c mrapp.ClientData) *entity.FormDataListFilter {
    var listFilter entity.FormDataListFilter

    parseFilterDetailing(c, &listFilter.Detailing)
    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *FormData) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *FormData) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateFormData{}

        if err := c.ParseAndValidate(&request); err != nil {
           return err
        }

        item := entity.FormData{
            Caption: request.Caption,
            Detailing: request.Detailing,
        }

        err := f.service.Create(c.Context(), &item)

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

func (f *FormData) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreFormData{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormData{
            Id: f.getItemId(c),
            Version: request.Version,
            Caption: request.Caption,
            Detailing: request.Detailing,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *FormData) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.FormData{
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

func (f *FormData) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *FormData) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
