package http_v1

import (
    "print-shop-back/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrctx"
    "github.com/mondegor/go-webcore/mrreq"
)

func parseFilterDetailing(c mrcore.ClientData, detailing *[]entity.ItemDetailing) {
    items, err := mrreq.EnumList(c.Request(), "detailing")

    if err == nil {
        var itemDetailing entity.ItemDetailing

        for _, item := range items {
            err = itemDetailing.ParseAndSet(item)

            if err != nil {
                mrctx.Logger(c.Context()).Warn(err.Error())
                continue
            }

            *detailing = append(*detailing, itemDetailing)
        }
    } else {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }

    if len(*detailing) == 0 {
        *detailing = append(*detailing, entity.ItemDetailingNormal)
    }
}

func parseFilterStatuses(c mrcore.ClientData, statuses *[]mrcom.ItemStatus) {
    items, err := mrreq.EnumList(c.Request(), "statuses")

    if err == nil {
        var itemStatus mrcom.ItemStatus

        for _, item := range items {
            if item == mrcom.ItemStatusRemoved.String() {
                continue
            }

            err = itemStatus.ParseAndSet(item)

            if err != nil {
                mrctx.Logger(c.Context()).Warn(err.Error())
                continue
            }

            *statuses = append(*statuses, itemStatus)
        }
    } else {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }

    if len(*statuses) == 0 {
        *statuses = append(*statuses, mrcom.ItemStatusEnabled)
    }
}

func parseFilterResourceStatuses(c mrcore.ClientData, statuses *[]entity.ResourceStatus) {
    items, err := mrreq.EnumList(c.Request(), "statuses")

    if err == nil {
        var resourceStatus entity.ResourceStatus

        for _, item := range items {
            if item == mrcom.ItemStatusRemoved.String() {
                continue
            }

            err = resourceStatus.ParseAndSet(item)

            if err != nil {
                mrctx.Logger(c.Context()).Warn(err.Error())
                continue
            }

            *statuses = append(*statuses, resourceStatus)
        }
    } else {
        mrctx.Logger(c.Context()).Warn(err.Error())
    }
}
