package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
)

type (
	// FormElementPostgres - comment struct.
	FormElementPostgres struct {
		client    mrstorage.DBConnManager
		condition mrstorage.SQLBuilderCondition
		sqlUpdate mrstorage.SQLBuilderUpdate
	}
)

// NewFormElementPostgres - создаёт объект FormElementPostgres.
func NewFormElementPostgres(
	client mrstorage.DBConnManager,
	condition mrstorage.SQLBuilderCondition,
	sqlUpdate mrstorage.SQLBuilderUpdate,
) *FormElementPostgres {
	return &FormElementPostgres{
		client:    client,
		condition: condition,
		sqlUpdate: sqlUpdate,
	}
}

// NewOrderMeta - comment method.
func (re *FormElementPostgres) NewOrderMeta(formID uuid.UUID) mrstorage.MetaGetter {
	return mrsql.NewEntityMeta(
		module.DBSchema+"."+module.DBTableNameSubmitFormElements,
		"element_id",
		re.condition.Where(func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc {
			return w.JoinAnd(
				w.Equal("form_id", formID),
			)
		}),
	)
}

// Fetch - comment method.
func (re *FormElementPostgres) Fetch(ctx context.Context, formID uuid.UUID) ([]entity.FormElement, error) {
	sql := `
        SELECT
            fe.element_id,
            fe.tag_version,
			fe.form_id,
            fe.param_name,
            fe.element_caption,
			fe.template_id,
			fe.template_version,
			fe.element_required,
            et.element_type,
            et.element_detailing,
			et.element_body,
            fe.created_at,
			fe.updated_at
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormElements + ` fe
        JOIN
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + ` et
        ON
            fe.template_id = et.template_id
        WHERE
            fe.form_id = $1 AND fe.deleted_at IS NULL
        ORDER BY
            fe.order_index ASC, fe.element_id ASC;`

	cursor, err := re.client.Conn(ctx).Query(
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
		var row entity.FormElement

		err = cursor.Scan(
			&row.ID,
			&row.TagVersion,
			&row.FormID,
			&row.ParamName,
			&row.Caption,
			&row.TemplateID,
			&row.TemplateVersion,
			&row.Required,
			&row.Type,
			&row.Detailing,
			&row.Body,
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
func (re *FormElementPostgres) FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.FormElement, error) {
	sql := `
        SELECT
            fe.tag_version,
			fe.form_id,
            fe.param_name,
            fe.element_caption,
			fe.template_id,
			fe.template_version,
			fe.element_required,
            et.element_type,
            et.element_detailing,
			et.element_body,
            fe.created_at,
			fe.updated_at
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormElements + ` fe
        JOIN
            ` + module.DBSchema + `.` + module.DBTableNameElementTemplates + ` et
        ON
            fe.template_id = et.template_id
        WHERE
            fe.element_id = $1 AND fe.deleted_at IS NULL
        LIMIT 1;`

	row := entity.FormElement{ID: rowID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.TagVersion,
		&row.FormID,
		&row.ParamName,
		&row.Caption,
		&row.TemplateID,
		&row.TemplateVersion,
		&row.Required,
		&row.Type,
		&row.Detailing,
		&row.Body,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// FetchIDByParamName - comment method.
func (re *FormElementPostgres) FetchIDByParamName(ctx context.Context, formID uuid.UUID, paramName string) (mrtype.KeyInt32, error) {
	sql := `
        SELECT
            element_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormElements + `
        WHERE
            form_id = $1 AND param_name = $2 AND deleted_at IS NULL
        LIMIT 1;`

	var rowID mrtype.KeyInt32

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		formID,
		paramName,
	).Scan(
		&rowID,
	)

	return rowID, err
}

// IsExist - comment method.
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *FormElementPostgres) IsExist(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        SELECT
            element_id
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormElements + `
        WHERE
            element_id = $1 AND deleted_at IS NULL
        LIMIT 1;`

	return re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&rowID,
	)
}

// Insert - comment method.
func (re *FormElementPostgres) Insert(ctx context.Context, row entity.FormElement) (mrtype.KeyInt32, error) {
	sql := `
        INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameSubmitFormElements + `
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

	err := re.client.Conn(ctx).QueryRow(
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

// Update - comment method.
func (re *FormElementPostgres) Update(ctx context.Context, row entity.FormElement) (int32, error) {
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
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormElements + `
        SET
            tag_version = tag_version + 1,
			updated_at = NOW(),
            ` + setStr + `
        WHERE
            element_id = $1 AND tag_version = $2 AND deleted_at IS NULL
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

// Delete - comment method.
func (re *FormElementPostgres) Delete(ctx context.Context, rowID mrtype.KeyInt32) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormElements + `
        SET
            tag_version = tag_version + 1,
			deleted_at = NOW()
        WHERE
            element_id = $1 AND deleted_at IS NULL;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}
