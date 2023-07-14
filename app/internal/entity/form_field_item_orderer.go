package entity

import "calc-user-data-back-adm/pkg/mrentity"

type Node struct {
    Id mrentity.KeyInt32
    PrevId mrentity.ZeronullInt32
    NextId mrentity.ZeronullInt32
    OrderField mrentity.ZeronullInt64
}
