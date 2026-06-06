package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrworkflow/itemstatus"

	"print-shop-back/internal/dictionaries/materialtype/module"
	"print-shop-back/internal/dictionaries/materialtype/section/pub/entity"
)

type (
	// MaterialTypePostgres - comment struct.
	MaterialTypePostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewMaterialTypePostgres - создаёт объект MaterialTypePostgres.
func NewMaterialTypePostgres(client mrstorage.DBConnManager) *MaterialTypePostgres {
	return &MaterialTypePostgres{
		client: client,
	}
}

// Fetch - comment method.
func (re *MaterialTypePostgres) Fetch(ctx context.Context, _ entity.MaterialTypeParams) ([]entity.MaterialType, error) {
	sql := `
        SELECT
            type_id,
            type_caption
        FROM
            ` + module.DBTableNameMaterialTypes + `
        WHERE
            type_status = $1 AND deleted_at IS NULL
        ORDER BY
            type_caption ASC, type_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		itemstatus.Enabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.MaterialType, 0)

	for cursor.Next() {
		var row entity.MaterialType

		err = cursor.Scan(
			&row.ID,
			&row.Caption,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}
