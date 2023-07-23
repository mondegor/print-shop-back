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

func (c *ItemMetaData) TableInfo() *entity.TableInfo {
    return c.tableInfo
}

func (c *ItemMetaData) PrepareSelect(query squirrel.SelectBuilder) squirrel.SelectBuilder {
    for _, cond := range c.conditions {
        query = query.Where(cond)
    }

    return query
}

func (c *ItemMetaData) PrepareUpdate(query squirrel.UpdateBuilder) squirrel.UpdateBuilder {
    for _, cond := range c.conditions {
        query = query.Where(cond)
    }

    return query
}

func (c *ItemMetaData) PrepareDelete(query squirrel.DeleteBuilder) squirrel.DeleteBuilder {
    for _, cond := range c.conditions {
        query = query.Where(cond)
    }

    return query
}
