package http_v1

import (
    "fmt"
    "net/http"
    "print-shop-back/internal/controller/view"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-sysmess/mrerr"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrview"
)

const (
    catalogLaminateListURL = "/v1/catalog-laminates"
    catalogLaminateItemURL = "/v1/catalog-laminates/:id"
    catalogLaminateChangeStatusURL = "/v1/catalog-laminates/:id/status"
)

type CatalogLaminate struct {
    service usecase.CatalogLaminateService
}

func NewCatalogLaminate(service usecase.CatalogLaminateService) *CatalogLaminate {
    return &CatalogLaminate{
        service: service,
    }
}

func (ht *CatalogLaminate) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, catalogLaminateListURL, ht.GetList())
    router.HttpHandlerFunc(http.MethodPost, catalogLaminateListURL, ht.Create())

    router.HttpHandlerFunc(http.MethodGet, catalogLaminateItemURL, ht.Get())
    router.HttpHandlerFunc(http.MethodPut, catalogLaminateItemURL, ht.Store())
    router.HttpHandlerFunc(http.MethodDelete, catalogLaminateItemURL, ht.Remove())

    router.HttpHandlerFunc(http.MethodPut, catalogLaminateChangeStatusURL, ht.ChangeStatus())
}

func (ht *CatalogLaminate) GetList() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        items, err := ht.service.GetList(c.Context(), ht.newListFilter(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, items)
    }
}

func (ht *CatalogLaminate) newListFilter(c mrcore.ClientData) *entity.CatalogLaminateListFilter {
    var listFilter entity.CatalogLaminateListFilter

    parseFilterStatuses(c, &listFilter.Statuses)

    return &listFilter
}

func (ht *CatalogLaminate) Get() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.service.GetItem(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponse(http.StatusOK, item)
    }
}

func (ht *CatalogLaminate) Create() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.CreateCatalogLaminateRequest{}

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
                return mrerr.NewListWith("article", err)
            }

            if usecase.ErrCatalogLaminateTypeNotFound.Is(err) {
                return mrerr.NewListWith("typeId", err)
            }

            return err
        }

        response := mrview.CreateItemResponse{
            ItemId: fmt.Sprintf("%d", item.Id),
            Message: mrctx.Locale(c.Context()).TranslateMessage(
                "msgCatalogLaminateSuccessCreated",
                "entity has been success created",
            ),
        }

        return c.SendResponse(http.StatusCreated, response)
    }
}

func (ht *CatalogLaminate) Store() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := view.StoreCatalogLaminateRequest{}

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
                return mrerr.NewListWith("typeId", err)
            }

            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogLaminate) ChangeStatus() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        request := mrcom.ChangeItemStatusRequest{}

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

func (ht *CatalogLaminate) Remove() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        err := ht.service.Remove(c.Context(), ht.getItemId(c))

        if err != nil {
            return err
        }

        return c.SendResponseNoContent()
    }
}

func (ht *CatalogLaminate) getItemId(c mrcore.ClientData) mrentity.KeyInt32 {
    id := mrentity.KeyInt32(c.RequestPath().GetInt("id"))

    if id > 0 {
        return id
    }

    return 0
}
