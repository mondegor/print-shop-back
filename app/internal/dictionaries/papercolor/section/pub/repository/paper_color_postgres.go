package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/entity"
)

type (
	// PaperColorPostgres - comment struct.
	PaperColorPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewPaperColorPostgres - создаёт объект PaperColorPostgres.
func NewPaperColorPostgres(client mrstorage.DBConnManager) *PaperColorPostgres {
	return &PaperColorPostgres{
		client: client,
	}
}

// Fetch - comment method.
func (re *PaperColorPostgres) Fetch(ctx context.Context, _ entity.PaperColorParams) ([]entity.PaperColor, error) {
	sql := `
        SELECT
            color_id,
            color_caption
        FROM
            ` + module.DBTableNamePaperColors + `
        WHERE
            color_status = $1 AND deleted_at IS NULL
        ORDER BY
            color_caption ASC, color_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.PaperColor, 0)

	for cursor.Next() {
		var row entity.PaperColor

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
