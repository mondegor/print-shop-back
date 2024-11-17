package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
)

type (
	// PrintFormatPostgres - comment struct.
	PrintFormatPostgres struct {
		client          mrstorage.DBConnManager
		sqlBuilder      mrstorage.SQLBuilder
		repoStatus      db.FieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus]
		repoSoftDeleter db.RowSoftDeleter[uint64]
		repoTotalRows   db.TotalRowsFetcher[uint64]
	}
)

// NewPrintFormatPostgres - создаёт объект PrintFormatPostgres.
func NewPrintFormatPostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *PrintFormatPostgres {
	return &PrintFormatPostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoStatus: db.NewFieldWithVersionUpdater[uint64, uint32, mrenum.ItemStatus](
			client,
			module.DBTableNamePrintFormats,
			"format_id",
			module.DBFieldTagVersion,
			"format_status",
			module.DBFieldDeletedAt,
		),
		repoSoftDeleter: db.NewRowSoftDeleter[uint64](
			client,
			module.DBTableNamePrintFormats,
			"format_id",
			module.DBFieldTagVersion,
			module.DBFieldDeletedAt,
		),
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNamePrintFormats,
		),
	}
}

// FetchWithTotal - comment method.
func (re *PrintFormatPostgres) FetchWithTotal(ctx context.Context, params entity.PrintFormatParams) (rows []entity.PrintFormat, countRows uint64, err error) {
	condition := re.sqlBuilder.Condition().Build(re.fetchCondition(params.Filter))

	total, err := re.repoTotalRows.Fetch(ctx, condition)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	if params.Pager.Size > total {
		params.Pager.Size = total
	}

	orderBy := re.sqlBuilder.OrderBy().Build(re.fetchOrderBy(params.Sorter))
	limit := re.sqlBuilder.Limit().Build(params.Pager.Index, params.Pager.Size)

	rows, err = re.fetch(ctx, condition, orderBy, limit, params.Pager.Size)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// Fetch - comment method.
func (re *PrintFormatPostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows uint64,
) ([]entity.PrintFormat, error) {
	whereStr, whereArgs := condition.ToSQL()

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
            ` + module.DBTableNamePrintFormats + `
        WHERE
            ` + whereStr + `
        ORDER BY
            ` + orderBy.String() + limit.String() + `;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		whereArgs...,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.PrintFormat, 0, maxRows)

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

// FetchOne - comment method.
func (re *PrintFormatPostgres) FetchOne(ctx context.Context, rowID uint64) (entity.PrintFormat, error) {
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
            ` + module.DBTableNamePrintFormats + `
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
func (re *PrintFormatPostgres) FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error) {
	return re.repoStatus.Fetch(ctx, rowID)
}

// Insert - comment method.
func (re *PrintFormatPostgres) Insert(ctx context.Context, row entity.PrintFormat) (rowID uint64, err error) {
	sql := `
        INSERT INTO ` + module.DBTableNamePrintFormats + `
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

	err = re.client.Conn(ctx).QueryRow(
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
func (re *PrintFormatPostgres) Update(ctx context.Context, row entity.PrintFormat) (tagVersion uint32, err error) {
	sql := `
        UPDATE
            ` + module.DBTableNamePrintFormats + `
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

	err = re.client.Conn(ctx).QueryRow(
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
func (re *PrintFormatPostgres) UpdateStatus(ctx context.Context, row entity.PrintFormat) (tagVersion uint32, err error) {
	return re.repoStatus.Update(ctx, row.ID, row.TagVersion, row.Status)
}

// Delete - comment method.
func (re *PrintFormatPostgres) Delete(ctx context.Context, rowID uint64) error {
	return re.repoSoftDeleter.Delete(ctx, rowID)
}

func (re *PrintFormatPostgres) fetchCondition(filter entity.PrintFormatListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.Expr("deleted_at IS NULL"),
				c.FilterLike("UPPER(format_caption)", strings.ToUpper(filter.SearchText)),
				c.FilterRangeFloat64("format_width", mrtype.RangeFloat64(filter.Width), 0, mrlib.EqualityThresholdE9),
				c.FilterRangeFloat64("format_height", mrtype.RangeFloat64(filter.Height), 0, mrlib.EqualityThresholdE9),
				c.FilterAnyOf("format_status", filter.Statuses),
			)
		},
	)
}

func (re *PrintFormatPostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("format_id", mrenum.SortDirectionASC),
			)
		},
	)
}
