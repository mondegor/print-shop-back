package repository

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-facture"
	entity "print-shop-back/internal/modules/dictionaries/paper-facture/entity/public-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	PaperFacturePostgres struct {
		client mrstorage.DBConn
	}
)

func NewPaperFacturePostgres(
	client mrstorage.DBConn,
) *PaperFacturePostgres {
	return &PaperFacturePostgres{
		client: client,
	}
}

func (re *PaperFacturePostgres) Fetch(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, error) {
	sql := `
        SELECT
            facture_id,
            facture_caption
        FROM
            ` + module.DBSchema + `.paper_factures
        WHERE
            facture_status = $1
        ORDER BY
            facture_caption ASC, facture_id ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.PaperFacture, 0)

	for cursor.Next() {
		var row entity.PaperFacture

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
