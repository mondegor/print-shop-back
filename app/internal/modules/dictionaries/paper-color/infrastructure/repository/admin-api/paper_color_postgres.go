package repository

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/paper-color"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/admin-api"
	repository_shared "print-shop-back/internal/modules/dictionaries/paper-color/infrastructure/repository/shared"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperColorPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
	}
)

func NewPaperColorPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
) *PaperColorPostgres {
	return &PaperColorPostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

func (re *PaperColorPostgres) NewSelectParams(params entity.PaperColorParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("color_status", mrenum.ItemStatusRemoved),
				w.FilterLike("UPPER(color_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("color_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("color_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *PaperColorPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.PaperColor, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            color_id,
            tag_version,
            color_caption as caption,
            color_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.paper_colors
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

	rows := make([]entity.PaperColor, 0)

	for cursor.Next() {
		var row entity.PaperColor

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.Caption,
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

func (re *PaperColorPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.paper_colors
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

func (re *PaperColorPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PaperColor, error) {
	sql := `
        SELECT
            tag_version,
            color_caption,
            color_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.paper_colors
        WHERE
            color_id = $1 AND color_status <> $2
        LIMIT 1;`

	row := entity.PaperColor{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

func (re *PaperColorPostgres) FetchStatus(ctx context.Context, row entity.PaperColor) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            color_status
        FROM
            ` + module.DBSchema + `.paper_colors
        WHERE
            color_id = $1 AND color_status <> $2
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
func (re *PaperColorPostgres) IsExists(ctx context.Context, rowID mrtype.KeyInt32) error {
	return repository_shared.PaperColorIsExistsPostgres(ctx, re.client, rowID)
}

func (re *PaperColorPostgres) Insert(ctx context.Context, row entity.PaperColor) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.paper_colors
            (
                color_caption,
                color_status
            )
        VALUES
            ($1, $2)
        RETURNING
            color_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

func (re *PaperColorPostgres) Update(ctx context.Context, row entity.PaperColor) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.paper_colors
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            color_caption = $4
        WHERE
            color_id = $1 AND tag_version = $2 AND color_status <> $3
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		mrenum.ItemStatusRemoved,
		row.Caption,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

func (re *PaperColorPostgres) UpdateStatus(ctx context.Context, row entity.PaperColor) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.paper_colors
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            color_status = $4
        WHERE
            color_id = $1 AND tag_version = $2 AND color_status <> $3
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

func (re *PaperColorPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.paper_colors
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            color_status = $2
        WHERE
            color_id = $1 AND color_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}
