package repository

import (
	"context"
	module "print-shop-back/internal/modules/catalog/laminate"
	entity "print-shop-back/internal/modules/catalog/laminate/entity/public-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	LaminatePostgres struct {
		client mrstorage.DBConn
	}
)

func NewLaminatePostgres(
	client mrstorage.DBConn,
) *LaminatePostgres {
	return &LaminatePostgres{
		client: client,
	}
}

func (re *LaminatePostgres) Fetch(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, error) {
	sql := `
        SELECT
            laminate_id,
			laminate_article,
			laminate_caption,
			type_id,
			laminate_length,
			laminate_weight,
			laminate_thickness
        FROM
            ` + module.DBSchema + `.laminates
        WHERE
            laminate_status = $1
        ORDER BY
            laminate_caption ASC, laminate_id ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.Laminate, 0)

	for cursor.Next() {
		var row entity.Laminate

		err = cursor.Scan(
			&row.ID,
			&row.Article,
			&row.Caption,
			&row.TypeID,
			&row.Length,
			&row.Weight,
			&row.Thickness,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}
