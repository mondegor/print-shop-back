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
    catalogLaminateGetListURL = "/v1/catalog-laminates"
    catalogLaminateGetItemURL = "/v1/catalog-laminates/:id"
    catalogLaminateCreateURL = "/v1/catalog-laminates"
    catalogLaminateStoreURL = "/v1/catalog-laminates/:id"
    catalogLaminateChangeStatusURL = "/v1/catalog-laminates/:id/status"
    catalogLaminateRemove = "/v1/catalog-laminates/:id"
)

type CatalogLaminate struct {
    service usecase.CatalogLaminateService
}

func NewCatalogLaminate(service usecase.CatalogLaminateService) *CatalogLaminate {
    return &CatalogLaminate{
        service: service,
    }
}

func (f *CatalogLaminate) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogLaminateGetListURL, f.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogLaminateGetItemURL, f.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogLaminateCreateURL, f.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogLaminateStoreURL, f.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogLaminateChangeStatusURL, f.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogLaminateRemove, f.Remove())
}

func (f *CatalogLaminate) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := f.service.GetList(c.Context(), f.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (f *CatalogLaminate) newListFilter(c mrapp.ClientData) *entity.CatalogLaminateListFilter {
    var listFilter entity.CatalogLaminateListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (f *CatalogLaminate) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := f.service.GetItem(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (f *CatalogLaminate) Create() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.CreateCatalogLaminate{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminate{
            Article: request.Article,
            Caption: request.Caption,
            TypeId: request.TypeId,
            Length: request.Length,
            Weight: request.Weight,
            Thickness: request.Thickness,
        }

        err := f.service.Create(c.Context(), &item)

        if err != nil {
            if usecase.ErrCatalogLaminateArticleAlreadyExists.Is(err) {
                return mrlib.NewUserErrorListWithError("article", err)
            }

            if usecase.ErrCatalogLaminateTypeNotFound.Is(err) {
                return mrlib.NewUserErrorListWithError("typeId", err)
            }

            return err
        }

        response := dto.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: c.Locale().GetMessage(
                "msgCatalogLaminateSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (f *CatalogLaminate) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogLaminate{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminate{
            Id: f.getItemId(c),
            Version: request.Version,
            Article: request.Article,
            Caption: request.Caption,
            TypeId: request.TypeId,
            Length: request.Length,
            Weight: request.Weight,
            Thickness: request.Thickness,
        }

        err := f.service.Store(c.Context(), &item)

        if err != nil {
            if usecase.ErrCatalogLaminateTypeNotFound.Is(err) {
                return mrlib.NewUserErrorListWithError("typeId", err)
            }

            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogLaminate) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminate{
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

func (f *CatalogLaminate) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := f.service.Remove(c.Context(), f.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (f *CatalogLaminate) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
