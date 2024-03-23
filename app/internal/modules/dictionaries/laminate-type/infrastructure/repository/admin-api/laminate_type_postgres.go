package repository

import (
	"context"
	module "print-shop-back/internal/modules/dictionaries/laminate-type"
	entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/admin-api"
	repository_shared "print-shop-back/internal/modules/dictionaries/laminate-type/infrastructure/repository/shared"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	LaminateTypePostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
	}
)

func NewLaminateTypePostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
) *LaminateTypePostgres {
	return &LaminateTypePostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

func (re *LaminateTypePostgres) NewSelectParams(params entity.LaminateTypeParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("type_status", mrenum.ItemStatusRemoved),
				w.FilterLike("UPPER(type_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("type_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("type_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *LaminateTypePostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.LaminateType, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            type_id,
            tag_version,
            type_caption as caption,
            type_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.laminate_types
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

	rows := make([]entity.LaminateType, 0)

	for cursor.Next() {
		var row entity.LaminateType

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

func (re *LaminateTypePostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.laminate_types
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

func (re *LaminateTypePostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.LaminateType, error) {
	sql := `
        SELECT
            tag_version,
            type_caption,
            type_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.laminate_types
        WHERE
            type_id = $1 AND type_status <> $2
        LIMIT 1;`

	row := entity.LaminateType{ID: rowID}

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

// FetchStatus
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *LaminateTypePostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository_shared.LaminateTypeFetchStatusPostgres(ctx, re.client, rowID)
}

func (re *LaminateTypePostgres) Insert(ctx context.Context, row entity.LaminateType) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.laminate_types
            (
                type_caption,
                type_status
            )
        VALUES
            ($1, $2)
        RETURNING
            type_id;`

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

func (re *LaminateTypePostgres) Update(ctx context.Context, row entity.LaminateType) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.laminate_types
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            type_caption = $4
        WHERE
            type_id = $1 AND tag_version = $2 AND type_status <> $3
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

func (re *LaminateTypePostgres) UpdateStatus(ctx context.Context, row entity.LaminateType) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.laminate_types
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            type_status = $4
        WHERE
            type_id = $1 AND tag_version = $2 AND type_status <> $3
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

func (re *LaminateTypePostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.laminate_types
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            type_status = $2
        WHERE
            type_id = $1 AND type_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}
