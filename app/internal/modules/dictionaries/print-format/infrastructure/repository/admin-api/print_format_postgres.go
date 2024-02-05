package repository

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/print-format"
	entity "print-shop-back/internal/modules/dictionaries/print-format/entity/admin-api"
	repository_shared "print-shop-back/internal/modules/dictionaries/print-format/infrastructure/repository/shared"
	"strings"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PrintFormatPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewPrintFormatPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *PrintFormatPostgres {
	return &PrintFormatPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *PrintFormatPostgres) NewFetchParams(params entity.PrintFormatParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("format_status", mrenum.ItemStatusRemoved),
				w.FilterLike("UPPER(format_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterRangeInt64("format_length", params.Filter.Length, 0),
				w.FilterRangeInt64("format_width", params.Filter.Width, 0),
				w.FilterAnyOf("format_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("format_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *PrintFormatPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.PrintFormat, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            format_id,
            tag_version,
            created_at as createdAt,
			updated_at as updatedAt,
            format_caption as caption,
            format_length as length,
            format_width as width,
            format_status
        FROM
            ` + module.DBSchema + `.print_formats
        WHERE
            ` + whereStr + `
        ORDER BY
            ` + params.OrderBy.String() + params.Pager.String() + `;`

	cursor, err := re.client.Query(
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
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Caption,
			&row.Length,
			&row.Width,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *PrintFormatPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.print_formats
        WHERE
            ` + whereStr + `;`

	var totalRow int64

	err := re.client.QueryRow(
		ctx,
		sql,
		whereArgs...,
	).Scan(
		&totalRow,
	)

	return totalRow, err
}

func (re *PrintFormatPostgres) LoadOne(ctx context.Context, row *entity.PrintFormat) error {
	sql := `
        SELECT
            tag_version,
            created_at,
			updated_at,
            format_caption,
            format_length,
            format_width,
            format_status
        FROM
            ` + module.DBSchema + `.print_formats
        WHERE
            format_id = $1 AND format_status <> $2
        LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.Caption,
		&row.Length,
		&row.Width,
		&row.Status,
	)
}

func (re *PrintFormatPostgres) FetchStatus(ctx context.Context, row *entity.PrintFormat) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            format_status
        FROM
            ` + module.DBSchema + `.print_formats
        WHERE
            format_id = $1 AND format_status <> $2
        LIMIT 1;`

	var status mrenum.ItemStatus

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&status,
	)

	return status, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *PrintFormatPostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	return repository_shared.PrintFormatIsExistsPostgres(ctx, re.client, id)
}

func (re *PrintFormatPostgres) Insert(ctx context.Context, row *entity.PrintFormat) error {
	sql := `
        INSERT INTO ` + module.DBSchema + `.print_formats
            (
                format_caption,
                format_length,
                format_width,
                format_status
            )
        VALUES
            ($1, $2, $3, $4)
        RETURNING
            format_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Length,
		row.Width,
		row.Status,
	).Scan(
		&row.ID,
	)

	return err
}

func (re *PrintFormatPostgres) Update(ctx context.Context, row *entity.PrintFormat) (int32, error) {
	set, err := re.sqlUpdate.SetFromEntity(row)

	if err != nil || set.Empty() {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
	}

	setStr, setArgs := set.Param(len(args) + 1).ToSql()

	sql := `
        UPDATE
            ` + module.DBSchema + `.print_formats
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            ` + setStr + `
        WHERE
            format_id = $1 AND tag_version = $2 AND format_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err = re.client.QueryRow(
		ctx,
		sql,
		mrsql.MergeArgs(args, setArgs)...,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *PrintFormatPostgres) UpdateStatus(ctx context.Context, row *entity.PrintFormat) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.print_formats
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            format_status = $4
        WHERE
            format_id = $1 AND tag_version = $2 AND format_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
		row.Status,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *PrintFormatPostgres) Delete(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.print_formats
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            format_status = $2
        WHERE
            format_id = $1 AND format_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		id,
		mrenum.ItemStatusRemoved,
	)
}
