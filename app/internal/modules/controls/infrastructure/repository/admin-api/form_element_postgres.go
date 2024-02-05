package repository

import (
	"context"
	module "print-shop-back/internal/modules/controls"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"
	"strings"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FormElementPostgres struct {
		client    mrstorage.DBConn
		sqlSelect mrstorage.SqlBuilderSelect
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewFormElementPostgres(
	client mrstorage.DBConn,
	sqlSelect mrstorage.SqlBuilderSelect,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *FormElementPostgres {
	return &FormElementPostgres{
		client:    client,
		sqlSelect: sqlSelect,
		sqlUpdate: sqlUpdate,
	}
}

func (re *FormElementPostgres) GetMetaData(formID mrtype.KeyInt32) mrorderer.EntityMeta {
	return mrorderer.NewEntityMeta(
		module.UnitFormElementsDBSchema+".form_elements",
		"element_id",
		re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.Equal("form_id", formID),
			)
		}),
	)
}

func (re *FormElementPostgres) NewFetchParams(params entity.FormElementParams) mrstorage.SqlSelectParams {
	return mrstorage.SqlSelectParams{
		Where: re.sqlSelect.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.Equal("fe.form_id", params.FormID),
				w.FilterLikeFields([]string{"UPPER(fe.param_name)", "UPPER(fe.element_caption)"}, strings.ToUpper(params.Filter.SearchText)),
				w.FilterAnyOf("et.element_detailing", params.Filter.Detailing),
			)
		}),
		OrderBy: re.sqlSelect.OrderBy(func(s mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc {
			return s.Join(
				s.Field(params.Sorter.FieldName, params.Sorter.Direction),
				s.Field("fe.element_id", mrenum.SortDirectionASC),
			)
		}),
		Pager: re.sqlSelect.Pager(func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc {
			return p.OffsetLimit(params.Pager.Index, params.Pager.Size)
		}),
	}
}

func (re *FormElementPostgres) Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.FormElement, error) {
	whereStr, whereArgs := params.Where.ToSql()

	sql := `
        SELECT
            fe.element_id,
            fe.tag_version,
            fe.created_at as createdAt,
			fe.updated_at as updatedAt,
			fe.form_id,
            fe.param_name as paramName,
            fe.element_caption as caption,
			fe.template_id,
			fe.element_required,
            et.element_type,
            et.element_detailing,
            et.element_body
        FROM
            ` + module.UnitFormElementsDBSchema + `.form_elements fe
        JOIN
            ` + module.UnitElementTemplateDBSchema + `.element_templates et
        ON
            fe.template_id = et.template_id
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

	rows := make([]entity.FormElement, 0)

	for cursor.Next() {
		row := entity.FormElement{}

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.FormID,
			&row.ParamName,
			&row.Caption,
			&row.TemplateID,
			&row.Required,
			&row.Type,
			&row.Detailing,
			&row.Body,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *FormElementPostgres) FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error) {
	whereStr, whereArgs := where.ToSql()

	sql := `
        SELECT
            COUNT(*)
        FROM
            ` + module.UnitFormElementsDBSchema + `.form_elements fe
        JOIN
            ` + module.UnitElementTemplateDBSchema + `.element_templates et
        ON
            fe.template_id = et.template_id
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

func (re *FormElementPostgres) LoadOne(ctx context.Context, row *entity.FormElement) error {
	sql := `
        SELECT
            fe.tag_version,
            fe.created_at,
			fe.updated_at,
			fe.form_id,
            fe.param_name,
            fe.element_caption,
			fe.template_id,
			fe.element_required,
            et.element_type,
            et.element_detailing,
            et.element_body
        FROM
            ` + module.UnitFormElementsDBSchema + `.form_elements fe
        JOIN
            ` + module.UnitElementTemplateDBSchema + `.element_templates et
        ON
            fe.template_id = et.template_id
        WHERE
            fe.element_id = $1
        LIMIT 1;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.ID,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.FormID,
		&row.ParamName,
		&row.Caption,
		&row.TemplateID,
		&row.Required,
		&row.Type,
		&row.Detailing,
		&row.Body,
	)

	return err
}

func (re *FormElementPostgres) FetchIdByName(ctx context.Context, formID mrtype.KeyInt32, paramName string) (mrtype.KeyInt32, error) {
	sql := `
        SELECT
            element_id
        FROM
            ` + module.UnitFormElementsDBSchema + `.form_elements
        WHERE
            form_id = $1 AND param_name = $2
        LIMIT 1;`

	var id mrtype.KeyInt32

	err := re.client.QueryRow(
		ctx,
		sql,
		formID,
		paramName,
	).Scan(
		&id,
	)

	return id, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *FormElementPostgres) IsExists(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
        SELECT
            1
        FROM
            ` + module.UnitFormElementsDBSchema + `.form_elements
        WHERE
            element_id = $1
        LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		id,
	).Scan(
		&id,
	)
}

func (re *FormElementPostgres) Insert(ctx context.Context, row *entity.FormElement) error {
	sql := `
        INSERT INTO ` + module.UnitFormElementsDBSchema + `.form_elements
            (
                form_id,
                param_name,
                element_caption,
				template_id,
                element_required
            )
        VALUES
            ($1, $2, $3, $4, $5)
        RETURNING
            element_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.FormID,
		row.ParamName,
		row.Caption,
		row.TemplateID,
		row.Required,
	).Scan(
		&row.ID,
	)

	return err
}

func (re *FormElementPostgres) Update(ctx context.Context, row *entity.FormElement) (int32, error) {
	set, err := re.sqlUpdate.SetFromEntityWith(row, func(s mrstorage.SqlBuilderSet) mrstorage.SqlBuilderPartFunc {
		return s.Field("element_required", row.Required)
	})

	if err != nil {
		return 0, err
	}

	args := []any{
		row.ID,
		row.TagVersion,
	}

	setStr, setArgs := set.Param(len(args) + 1).ToSql()

	sql := `
        UPDATE
            ` + module.UnitFormElementsDBSchema + `.form_elements
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            ` + setStr + `
        WHERE
            element_id = $1 AND tag_version = $2
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

func (re *FormElementPostgres) Delete(ctx context.Context, id mrtype.KeyInt32) error {
	sql := `
        DELETE FROM
            ` + module.UnitFormElementsDBSchema + `.form_elements
        WHERE
            element_id = $1;`

	return re.client.Exec(
		ctx,
		sql,
		id,
	)
}
