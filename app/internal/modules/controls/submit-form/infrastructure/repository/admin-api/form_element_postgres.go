package repository

import (
	"context"
	module "print-shop-back/internal/modules/controls/submit-form"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FormElementPostgres struct {
		client    mrstorage.DBConn
		condition mrstorage.SqlBuilderCondition
		sqlUpdate mrstorage.SqlBuilderUpdate
	}
)

func NewFormElementPostgres(
	client mrstorage.DBConn,
	condition mrstorage.SqlBuilderCondition,
	sqlUpdate mrstorage.SqlBuilderUpdate,
) *FormElementPostgres {
	return &FormElementPostgres{
		client:    client,
		condition: condition,
		sqlUpdate: sqlUpdate,
	}
}

func (re *FormElementPostgres) NewOrderMeta(formID uuid.UUID) mrorderer.EntityMeta {
	return mrorderer.NewEntityMeta(
		module.DBSchema+".submit_form_elements",
		"element_id",
		re.condition.Where(func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc {
			return w.JoinAnd(
				w.Equal("form_id", formID),
			)
		}),
	)
}

func (re *FormElementPostgres) Fetch(ctx context.Context, formID uuid.UUID) ([]entity.FormElement, error) {
	// :TODO: избавиться от прямой зависимости element_templates
	sql := `
        SELECT
            fe.element_id,
            fe.tag_version,
            fe.created_at,
			fe.updated_at,
			fe.form_id,
            fe.param_name,
            fe.element_caption,
			fe.template_id,
			fe.template_version,
			fe.element_required,
            et.element_type,
            et.element_detailing
        FROM
            ` + module.DBSchema + `.submit_form_elements fe
        JOIN
            ` + module.DBSchema + `.element_templates et
        ON
            fe.template_id = et.template_id
        WHERE
            fe.form_id = $1
        ORDER BY
            fe.order_index ASC, fe.element_id ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		formID,
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
			&row.TemplateVersion,
			&row.Required,
			&row.Type,
			&row.Detailing,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *FormElementPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.FormElement, error) {
	// :TODO: избавиться от прямой зависимости element_templates
	sql := `
        SELECT
            fe.tag_version,
            fe.created_at,
			fe.updated_at,
			fe.form_id,
            fe.param_name,
            fe.element_caption,
			fe.template_id,
			fe.template_version,
			fe.element_required,
            et.element_type,
            et.element_detailing
        FROM
            ` + module.DBSchema + `.submit_form_elements fe
        JOIN
            ` + module.DBSchema + `.element_templates et
        ON
            fe.template_id = et.template_id
        WHERE
            fe.element_id = $1
        LIMIT 1;`

	row := entity.FormElement{ID: rowID}

	err := re.client.QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.FormID,
		&row.ParamName,
		&row.Caption,
		&row.TemplateID,
		&row.TemplateVersion,
		&row.Required,
		&row.Type,
		&row.Detailing,
	)

	return row, err
}

func (re *FormElementPostgres) FetchIdByParamName(ctx context.Context, formID uuid.UUID, paramName string) (mrtype.KeyInt32, error) {
	sql := `
        SELECT
            element_id
        FROM
            ` + module.DBSchema + `.submit_form_elements
        WHERE
            form_id = $1 AND param_name = $2
        LIMIT 1;`

	var rowID mrtype.KeyInt32

	err := re.client.QueryRow(
		ctx,
		sql,
		formID,
		paramName,
	).Scan(
		&rowID,
	)

	return rowID, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *FormElementPostgres) IsExists(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        SELECT
            element_id
        FROM
            ` + module.DBSchema + `.submit_form_elements
        WHERE
            element_id = $1
        LIMIT 1;`

	return re.client.QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&rowID,
	)
}

func (re *FormElementPostgres) Insert(ctx context.Context, row entity.FormElement) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.submit_form_elements
            (
                form_id,
                param_name,
                element_caption,
				template_id,
				template_version,
                element_required
            )
        VALUES
            ($1, $2, $3, $4, $5, $6)
        RETURNING
            element_id;`

	err := re.client.QueryRow(
		ctx,
		sql,
		row.FormID,
		row.ParamName,
		row.Caption,
		row.TemplateID,
		row.TemplateVersion,
		row.Required,
	).Scan(
		&row.ID,
	)

	return row.ID, err
}

func (re *FormElementPostgres) Update(ctx context.Context, row entity.FormElement) (int32, error) {
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
            ` + module.DBSchema + `.submit_form_elements
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

func (re *FormElementPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        DELETE FROM
            ` + module.DBSchema + `.submit_form_elements
        WHERE
            element_id = $1;`

	return re.client.Exec(
		ctx,
		sql,
		rowID,
	)
}
