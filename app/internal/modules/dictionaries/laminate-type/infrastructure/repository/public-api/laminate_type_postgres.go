package repository

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/laminate-type"
	entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/public-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	LaminateTypePostgres struct {
		client mrstorage.DBConn
	}
)

func NewLaminateTypePostgres(
	client mrstorage.DBConn,
) *LaminateTypePostgres {
	return &LaminateTypePostgres{
		client: client,
	}
}

func (re *LaminateTypePostgres) Fetch(ctx context.Context, params entity.LaminateTypeParams) ([]entity.LaminateType, error) {
	sql := `
        SELECT
            type_id,
            type_caption
        FROM
            ` + module.DBSchema + `.laminate_types
        WHERE
            type_status = $1
        ORDER BY
            type_caption ASC, type_id ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.LaminateType, 0)

	for cursor.Next() {
		var row entity.LaminateType

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
