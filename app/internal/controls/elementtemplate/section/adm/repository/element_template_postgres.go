package repository

import (
	"context"
	"strings"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/module"
	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// ElementTemplatePostgres - comment struct.
	ElementTemplatePostgres struct {
		client    mrstorage.DBConnManager
		sqlSelect mrstorage.SQLBuilderSelect
		sqlUpdate mrstorage.SQLBuilderUpdate
	}
)

// NewElementTemplatePostgres - создаёт объект ElementTemplatePostgres.
func NewElementTemplatePostgres(
	client mrstorage.DBConnManager,
	sqlSelect mrstorage.SQLBuilderSelect,
	sqlUpdate mrstorage.SQLBuilderUpdate,
) *ElementTemplatePostgres {
	return &ElementTemplatePostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

// NewSelectParams - comment method.
func (re *ElementTemplatePostgres) NewSelectParams(params entity.ElementTemplateParams) mrstorage.SQLSelectParams {
	return mrstorage.SQLSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Expr("deleted_at IS NULL"),
				w.FilterLikeFields([]string{"UPPER(param_name)", "UPPER(template_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("element_detailing", params.Filter.Detailing),
				w.FilterAnyOf("template_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("template_id", mrenum.SortDirectionASC),
			)
		}),
		Limit: re.sqlSelect.Limit(func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

// Fetch - comment method.
func (re *ElementTemplatePostgres) Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.ElementTemplate, error) {
	whereStr, whereArgs := params.Where.ToSQL()

	sql := `
        SELECT
            template_id,
            tag_version,
            param_name as paramName,
            template_caption as caption,
            element_type,
            element_detailing,
            template_status,
            created_at as createdAt,
			updated_at as updatedAt
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + `
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

	rows := make([]entity.ElementTemplate, 0)

	for cursor.Next() {
		var row entity.ElementTemplate

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.ParamName,
			&row.Caption,
			&row.Type,
			&row.Detailing,
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
func (re *ElementTemplatePostgres) FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSQL()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + `
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
func (re *ElementTemplatePostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.ElementTemplate, error) {
	sql := `
        SELECT
            tag_version,
            param_name,
            template_caption,
            element_type,
            element_detailing,
            element_body,
            template_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + `
        WHERE
            template_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	row := entity.ElementTemplate{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.ParamName,
		&row.Caption,
		&row.Type,
		&row.Detailing,
		&row.Body,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchStatus - comment method.
// result: mrenum.ItemStatus - exists, ErrStorageNoRowFound - not exists, error - query error.
func (re *ElementTemplatePostgres) FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            template_status
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + `
        WHERE
            template_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	var status mrenum.ItemStatus

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&status,
	)

	return status, err
}

// Insert - comment method.
func (re *ElementTemplatePostgres) Insert(ctx context.Context, row entity.ElementTemplate) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + `
            (
                param_name,
                template_caption,
                element_type,
                element_detailing,
                element_body,
                template_status
            )
        VALUES
            ($1, $2, $3, $4, $5, $6)
        RETURNING
            template_id;`

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.ParamName,
		row.Caption,
		row.Type,
		row.Detailing,
		row.Body,
		row.Status,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

// Update - comment method.
func (re *ElementTemplatePostgres) Update(ctx context.Context, row entity.ElementTemplate) (int32, error) {
	set, err := re.sqlUpdate.SetFromEntity(row)

	if err != nil || set.Empty() {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
	}

	setStr, setArgs := set.WithParam(len(args) + 1).ToSQL()

	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            ` + setStr + `
        WHERE
            template_id = $1 AND tag_version = $2 AND deleted_at IS NULL
		RETURNING
			tag_version;`

	var tagVersion int32

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		mrsql.MergeArgs(args, setArgs)...,
	).Scan(
		&tagVersion,
	)

	return tagVersion, err
}

// UpdateStatus - comment method.
func (re *ElementTemplatePostgres) UpdateStatus(ctx context.Context, row entity.ElementTemplate) (int32, error) {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
			template_status = $3
        WHERE
            template_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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
func (re *ElementTemplatePostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + `
        SET
            tag_version = tag_version + 1,
			deleted_at = NOW()
        WHERE
            template_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}
