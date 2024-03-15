package repository

import (
	"context"
	module "print-shop-back/internal/modules/controls"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	"strings"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplatePostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewElementTemplatePostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *ElementTemplatePostgres {
	return &ElementTemplatePostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *ElementTemplatePostgres) NewFetchParams(params entity.ElementTemplateParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.NotEqual("template_status", mrenum.ItemStatusRemoved),
				w.FilterLikeFields([]string{"UPPER(param_name)", "UPPER(template_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("element_detailing", params.Filter.Detailing),
				w.FilterAnyOf("template_status", params.Filter.Statuses),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("template_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *ElementTemplatePostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.ElementTemplate, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            template_id,
            tag_version,
            created_at as createdAt,
			updated_at as updatedAt,
            param_name as paramName,
            template_caption as caption,
            element_type,
            element_detailing,
            element_body,
            template_status
        FROM
            ` + module.UnitElementTemplateDBSchema + `.element_templates
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

	rows := make([]entity.ElementTemplate, 0)

	for cursor.Next() {
		var row entity.ElementTemplate

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.ParamName,
			&row.Caption,
			&row.Type,
			&row.Detailing,
			&row.Body,
			&row.Status,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *ElementTemplatePostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.UnitElementTemplateDBSchema + `.element_templates
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

func (re *ElementTemplatePostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.ElementTemplate, error) {
	sql := `
        SELECT
            tag_version,
            created_at,
			updated_at,
            param_name,
            template_caption,
            element_type,
            element_detailing,
            element_body,
            template_status
        FROM
            ` + module.UnitElementTemplateDBSchema + `.element_templates
        WHERE
            template_id = $1 AND template_status <> $2
        LIMIT 1;`

	row := entity.ElementTemplate{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.ParamName,
		&row.Caption,
		&row.Type,
		&row.Detailing,
		&row.Body,
		&row.Status,
	)

	return row, err
}

func (re *ElementTemplatePostgres) FetchStatus(ctx context.Context, row entity.ElementTemplate) (mrenum.ItemStatus, error) {
	sql := `
        SELECT
            template_status
        FROM
            ` + module.UnitElementTemplateDBSchema + `.element_templates
        WHERE
            template_id = $1 AND template_status <> $2
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
func (re *ElementTemplatePostgres) IsExists(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        SELECT
            template_id
        FROM
            ` + module.UnitElementTemplateDBSchema + `.element_templates
        WHERE
            template_id = $1 AND template_status <> $2
        LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	).Scan(
		&rowID,
	)
}

func (re *ElementTemplatePostgres) Insert(ctx context.Context, row entity.ElementTemplate) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.UnitElementTemplateDBSchema + `.element_templates
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

	err := re.client.QueryRow(
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

func (re *ElementTemplatePostgres) Update(ctx context.Context, row entity.ElementTemplate) (int32, error) {
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
            ` + module.UnitElementTemplateDBSchema + `.element_templates
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            ` + setStr + `
        WHERE
            template_id = $1 AND tag_version = $2 AND template_status <> $3
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

func (re *ElementTemplatePostgres) UpdateStatus(ctx context.Context, row entity.ElementTemplate) (int32, error) {
	sql := `
        UPDATE
            ` + module.UnitElementTemplateDBSchema + `.element_templates
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            template_status = $4
        WHERE
            template_id = $1 AND tag_version = $2 AND template_status <> $3
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

func (re *ElementTemplatePostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.UnitElementTemplateDBSchema + `.element_templates
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            param_name = NULL,
            template_status = $2
        WHERE
            template_id = $1 AND template_status <> $2;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
		mrenum.ItemStatusRemoved,
	)
}
