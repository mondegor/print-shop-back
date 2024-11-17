package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/entity"
)

type (
	// PaperFacturePostgres - comment struct.
	PaperFacturePostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewPaperFacturePostgres - создаёт объект PaperFacturePostgres.
func NewPaperFacturePostgres(client mrstorage.DBConnManager) *PaperFacturePostgres {
	return &PaperFacturePostgres{
		client: client,
	}
}

// Fetch - comment method.
func (re *PaperFacturePostgres) Fetch(ctx context.Context, _ entity.PaperFactureParams) ([]entity.PaperFacture, error) {
	sql := `
        SELECT
            facture_id,
            facture_caption
        FROM
            ` + module.DBTableNamePaperFactures + `
        WHERE
            facture_status = $1 AND deleted_at IS NULL
        ORDER BY
            facture_caption ASC, facture_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
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
