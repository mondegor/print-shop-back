package http_v1

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/mrapp"
    "calc-user-data-back-adm/pkg/mrcontext"
)

func parseFilterDetailing(c mrapp.ClientData, detailing *[]entity.ItemDetailing) {
    items, err := mrcontext.EnumListFromRequest(c.Request(), "detailing")

    if err == nil {
        var itemDetailing entity.ItemDetailing

        for _, item := range items {
            err = itemDetailing.ParseAndSet(item)

            if err != nil {
                c.Logger().Warn(err.Error())
                continue
            }

            *detailing = append(*detailing, itemDetailing)
        }
    } else {
        c.Logger().Warn(err.Error())
    }

    if len(*detailing) == 0 {
        *detailing = append(*detailing, entity.ItemDetailingNormal)
    }
}

func parseFilterStatuses(c mrapp.ClientData, statuses *[]entity.ItemStatus) {
    items, err := mrcontext.EnumListFromRequest(c.Request(), "statuses")

    if err == nil {
        var itemStatus entity.ItemStatus

        for _, item := range items {
            if item == entity.ItemStatusRemoved.String() {
                continue
            }

            err = itemStatus.ParseAndSet(item)

            if err != nil {
                c.Logger().Warn(err.Error())
                continue
            }

            *statuses = append(*statuses, itemStatus)
        }
    } else {
        c.Logger().Warn(err.Error())
    }

    if len(*statuses) == 0 {
        *statuses = append(*statuses, entity.ItemStatusEnabled)
    }
}
