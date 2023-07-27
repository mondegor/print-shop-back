package repository

import (
    "print-shop-back/internal/entity"

    "github.com/Masterminds/squirrel"
)

type (
    Condition any

	ItemMetaData struct {
        tableInfo  *entity.TableInfo
        conditions []Condition
    }
)

func NewItemMetaData(tableName string, primaryKeyName string, conds []Condition) *ItemMetaData {
    return &ItemMetaData{
        tableInfo: &entity.TableInfo{
            Name:       tableName,
            PrimaryKey: primaryKeyName,
        },
        conditions: conds,
    }
}

func (mt *ItemMetaData) TableInfo() *entity.TableInfo {
    return mt.tableInfo
}

func (mt *ItemMetaData) PrepareSelect(query squirrel.SelectBuilder) squirrel.SelectBuilder {
    for _, cond := range mt.conditions {
        query = query.Where(cond)
    }

    return query
}

func (mt *ItemMetaData) PrepareUpdate(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
    for _, cond := range mt.conditions {
        query = query.Where(cond)
    }

    return query
}

func (mt *ItemMetaData) PrepareDelete(query squirrel.DeleteBuilder) squirrel.DeleteBuilder {
    for _, cond := range mt.conditions {
        query = query.Where(cond)
    }

    return query
}
