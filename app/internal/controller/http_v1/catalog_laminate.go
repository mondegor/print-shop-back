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
    catalogLaminateRemoveURL = "/v1/catalog-laminates/:id"
)

type CatalogLaminate struct {
    service usecase.CatalogLaminateService
}

func NewCatalogLaminate(service usecase.CatalogLaminateService) *CatalogLaminate {
    return &CatalogLaminate{
        service: service,
    }
}

func (ht *CatalogLaminate) AddHandlers(router mrapp.Router) {
    router.HttpHandlerFunc(http.MethodGet, catalogLaminateGetListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodGet, catalogLaminateGetItemURL, ht.GetItem())
    router.HttpHandlerFunc(http.MethodPost, catalogLaminateCreateURL, ht.Create())
    router.HttpHandlerFunc(http.MethodPut, catalogLaminateStoreURL, ht.Store())
    router.HttpHandlerFunc(http.MethodPut, catalogLaminateChangeStatusURL, ht.ChangeStatus())
    router.HttpHandlerFunc(http.MethodDelete, catalogLaminateRemoveURL, ht.Remove())
}

func (ht *CatalogLaminate) GetList() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogLaminate) newListFilter(c mrapp.ClientData) *entity.CatalogLaminateListFilter {
    var listFilter entity.CatalogLaminateListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogLaminate) GetItem() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogLaminate) Create() mrapp.HttpHandlerFunc {
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

        err := ht.service.Create(c.Context(), &item)

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

func (ht *CatalogLaminate) Store() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.StoreCatalogLaminate{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminate{
            Id:        ht.getItemId(c),
            Version:   request.Version,
            Article:   request.Article,
            Caption:   request.Caption,
            TypeId:    request.TypeId,
            Length:    request.Length,
            Weight:    request.Weight,
            Thickness: request.Thickness,
        }

        err := ht.service.Store(c.Context(), &item)

        if err != nil {
            if usecase.ErrCatalogLaminateTypeNotFound.Is(err) {
                return mrlib.NewUserErrorListWithError("typeId", err)
            }

            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogLaminate) ChangeStatus() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        request := dto.ChangeItemStatus{}

        if err := c.ParseAndValidate(&request); err != nil {
            return err
        }

        item := entity.CatalogLaminate{
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

func (ht *CatalogLaminate) Remove() mrapp.HttpHandlerFunc {
    return func(c mrapp.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogLaminate) getItemId(c mrapp.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
