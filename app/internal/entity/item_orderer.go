package entity

import "print-shop-back/pkg/mrentity"

const ModelNameItemOrderer = "ItemOrderer"

type (
    ItemOrdererNode struct {
        Id mrentity.KeyInt32
        PrevId mrentity.ZeronullInt32
        NextId mrentity.ZeronullInt32
        OrderField mrentity.ZeronullInt64
    }
)
