package repository

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/box/module"
	"github.com/mondegor/print-shop-back/internal/calculations/box/section/pub/entity"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// BoxPostgres - comment struct.
	BoxPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewBoxPostgres - создаёт объект BoxPostgres.
func NewBoxPostgres(client mrstorage.DBConnManager) *BoxPostgres {
	return &BoxPostgres{
		client: client,
	}
}

// FetchOne - comment method.
func (re *BoxPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.CalcResult, error) {
	sql := `
		SELECT
			tag_version,
			result_article,
			result_caption,
			result_length,
			result_width,
			result_height,
			result_status,
			created_at,
			updated_at
		FROM
			` + module.DBSchema + `.` + module.DBTableNameBox + `
		WHERE
			result_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.CalcResult{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// Insert - comment method.
func (re *BoxPostgres) Insert(ctx context.Context, row entity.CalcResult) (mrtype.KeyInt32, error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameBox + `
			(
				result_caption
			)
		VALUES
			($1)
		RETURNING
			result_id;`

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Caption,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}
