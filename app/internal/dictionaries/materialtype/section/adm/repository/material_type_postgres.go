package repository

import (
	"context"
	"strings"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/shared/repository"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// MaterialTypePostgres - comment struct.
	MaterialTypePostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
	}
)

// NewMaterialTypePostgres - создаёт объект MaterialTypePostgres.
func NewMaterialTypePostgres(client mrstorage.DBConnManager, sqlSelect mrstorage.SQLBuilderSelect) *MaterialTypePostgres {
	return &MaterialTypePostgres{
		client:    client,
		sqlSelect: sqlSelect,
	}
}

// NewSelectParams - comment method.
func (re *MaterialTypePostgres) NewSelectParams(params entity.MaterialTypeParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLike("UPPER(type_caption)", strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("type_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("type_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *MaterialTypePostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.MaterialType, error) {
	whereStr, whereArgs := params.Where.ToSQL()

	sql := `
        SELECT
            type_id,
            tag_version,
            type_caption as caption,
            type_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameMaterialTypes + `
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

	rows := make([]entity.MaterialType, 0)

	for cursor.Next() {
		var row entity.MaterialType

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

// FetchTotal - comment method.
func (re *MaterialTypePostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameMaterialTypes + `
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
func (re *MaterialTypePostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.MaterialType, error) {
	sql := `
        SELECT
            tag_version,
            type_caption,
            type_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameMaterialTypes + `
        WHERE
            type_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.MaterialType{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.Caption,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *MaterialTypePostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	return repository.MaterialTypeFetchStatusPostgres(ctx, re.client, rowID)
}

// Insert - comment method.
func (re *MaterialTypePostgres) Insert(ctx context.Context, row entity.MaterialType) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameMaterialTypes + `
            (
                type_caption,
                type_status
            )
        VALUES
            ($1, $2)
        RETURNING
            type_id;`

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *MaterialTypePostgres) Update(ctx context.Context, row entity.MaterialType) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameMaterialTypes + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            type_caption = $3
        WHERE
            type_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ID,
		row.TagVersion,
		row.Caption,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// UpdateStatus - comment method.
func (re *MaterialTypePostgres) UpdateStatus(ctx context.Context, row entity.MaterialType) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameMaterialTypes + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            type_status = $3
        WHERE
            type_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *MaterialTypePostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameMaterialTypes + `
        SET
            tag_version = tag_version + 1,
			deleted_at = NOW()
        WHERE
            type_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}
