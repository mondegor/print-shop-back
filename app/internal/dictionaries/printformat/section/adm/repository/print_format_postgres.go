package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/shared/repository"
)

type (
	// PrintFormatPostgres - comment struct.
	PrintFormatPostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
	}
)

// NewPrintFormatPostgres - создаёт объект PrintFormatPostgres.
func NewPrintFormatPostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect) *PrintFormatPostgres {
	return &PrintFormatPostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

// NewSelectParams - comment method.
func (re *PrintFormatPostgres) NewSelectParams(params entity.PrintFormatParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLike("UPPER(format_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterRangeFloat64("format_width", mrtype.RangeFloat64(params.Filter.Width), 0, mrlib.EqualityThresholdE9),
				w.FilterRangeFloat64("format_height", mrtype.RangeFloat64(params.Filter.Height), 0, mrlib.EqualityThresholdE9),
				w.FilterAnyOf("format_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("format_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *PrintFormatPostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.PrintFormat, error) {
	whereStr, whereArgs := params.Where.ToSQL()

	sql := `
        SELECT
            format_id,
            tag_version,
            format_caption as caption,
            format_width as width,
            format_height as height,
            format_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
        WHERE
            ` + whereStr + `
        ORDER BY
            ` + params.OrderBy.String() + params.Limit.String() + `;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		whereArgs...,
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
			&row.TagVersion,
			&row.Caption,
			&row.Width,
			&row.Height,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

// FetchTotal - comment method.
func (re *PrintFormatPostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
        WHERE
            ` + whereStr + `;`

	var totalRow int64

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		whereArgs...,
	).Scan(
		&totalRow,
	)

	return totalRow, err
}

// FetchOne - comment method.
func (re *PrintFormatPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PrintFormat, error) {
	sql := `
        SELECT
            tag_version,
            format_caption,
            format_width,
            format_height,
            format_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
        WHERE
            format_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.PrintFormat{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.Width,
		&row.Height,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *PrintFormatPostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository.PrintFormatFetchStatusPostgres(ctx, re.client, rowID)
}

// Insert - comment method.
func (re *PrintFormatPostgres) Insert(ctx context.Context, row entity.PrintFormat) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
            (
                format_caption,
                format_width,
                format_height,
                format_status
            )
        VALUES
            ($1, $2, $3, $4)
        RETURNING
            format_id;`

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Width,
		row.Height,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *PrintFormatPostgres) Update(ctx context.Context, row entity.PrintFormat) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
			format_caption = $3,
			format_width = $4,
			format_height = $5
        WHERE
            format_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		row.Caption,
		row.Width,
		row.Height,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// UpdateStatus - comment method.
func (re *PrintFormatPostgres) UpdateStatus(ctx context.Context, row entity.PrintFormat) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            format_status = $3
        WHERE
            format_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		row.Status,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// Delete - comment method.
func (re *PrintFormatPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNamePrintFormats + `
        SET
            tag_version = tag_version + 1,
			deleted_at = NOW()
        WHERE
            format_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}
