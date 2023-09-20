package http_v1

import (
    "net/http"
    "print-shop-back/internal/usecase"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    commonImagesURL = "/images/*path"
)

type Common struct {
    serviceFile usecase.FileService
}

func NewCommon(serviceFile usecase.FileService) *Common {
    return &Common{
        serviceFile: serviceFile,
    }
}

func (ht *Common) AddHandlers(router mrcore.HttpRouter) {
    router.HttpHandlerFunc(http.MethodGet, commonImagesURL, ht.GetImage())
}

func (ht *Common) GetImage() mrcore.HttpHandlerFunc {
    return func(c mrcore.ClientData) error {
        item, err := ht.serviceFile.Get(c.Context(), c.RequestPath().Get("path"))

        if err != nil {
            return err
        }

        defer item.Body.Close()

        return c.SendFile(item.ContentType, item.Body)
    }
}
