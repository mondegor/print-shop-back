package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/entity"
)

type (
	// PrintFormatPostgres - comment struct.
	PrintFormatPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewPrintFormatPostgres - создаёт объект PrintFormatPostgres.
func NewPrintFormatPostgres(client mrstorage.DBConnManager) *PrintFormatPostgres {
	return &PrintFormatPostgres{
		client: client,
	}
}

// Fetch - comment method.
func (re *PrintFormatPostgres) Fetch(ctx context.Context, _ entity.PrintFormatParams) ([]entity.PrintFormat, error) {
	sql := `
        SELECT
            format_id,
            format_caption
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
        WHERE
            format_status = $1 AND deleted_at IS NULL
        ORDER BY
            format_caption ASC, format_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		mrenum.ItemStatusEnabled,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.PrintFormat, 0)

	for cursor.Next() {
		var row entity.PrintFormat

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
