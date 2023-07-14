package repository

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/client/mrpostgres"
    "calc-user-data-back-adm/pkg/mrentity"
    "calc-user-data-back-adm/pkg/mrerr"
    "context"

    "github.com/Masterminds/squirrel"
)

type FormFieldItem struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewFormFieldItem(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *FormFieldItem {
    return &FormFieldItem{
        client: client,
        builder: queryBuilder,
    }
}

func (f *FormFieldItem) LoadAll(ctx context.Context, listFilter *entity.FormFieldItemListFilter, rows *[]entity.FormFieldItem) error {
    tbl := f.builder.
        Select(`
            ff.field_id,
            ff.template_id,
            ff.tag_version,
            ff.datetime_created,
            ff.param_name,
            ff.field_caption,
            fft.field_type,
            fft.field_detailing,
            fft.field_body,
            ff.field_required,
            ff.order_field`).
        From("public.form_fields ff").
        Join("public.form_field_templates fft ON ff.template_id = fft.template_id").
        Where(squirrel.Eq{"ff.form_id": listFilter.FormId}).
        Where(squirrel.NotEq{"ff.field_status": entity.ItemStatusRemoved}).
        OrderBy("ff.order_field ASC, ff.field_caption ASC, ff.field_id ASC")

    if len(listFilter.Detailing) > 0 {
        tbl = tbl.Where(squirrel.Eq{"fft.field_detailing": listFilter.Detailing})
    }

    sql, args, err := tbl.ToSql()

    if err != nil {
        return mrerr.ErrInternal.Wrap(err)
    }

    cursor, err := f.client.Query(ctx, sql, args...)

    if err != nil {
        return err
    }

    for cursor.Next() {
        row := entity.FormFieldItem{FormId: listFilter.FormId}

        err = cursor.Scan(
            &row.Id,
            &row.TemplateId,
            &row.Version,
            &row.CreatedAt,
            &row.ParamName,
            &row.Caption,
            &row.Type,
            &row.Detailing,
            &row.Body,
            &row.Required,
            &row.OrderField,
        )

        if err != nil {
            return mrerr.ErrStorageFetchDataFailed.Wrap(err)
        }

        *rows = append(*rows, row)
    }

    if err = cursor.Err(); err != nil {
        return mrerr.ErrStorageFetchDataFailed.Wrap(err)
    }

    return nil
}

// LoadOne
// uses: row{Id, FormId}
// modifies: row{TemplateId, Version, CreatedAt, ParamName, Caption, Required, OrderField}
func (f *FormFieldItem) LoadOne(ctx context.Context, row *entity.FormFieldItem) error {
    sql := `
        SELECT
            ff.template_id,
            ff.tag_version,
            ff.datetime_created,
            ff.param_name,
            ff.field_caption,
            fft.field_type,
            fft.field_detailing,
            fft.field_body,
            ff.field_required,
            ff.order_field
        FROM
            public.form_fields ff
        JOIN
            public.form_field_templates fft
        ON
            ff.template_id = fft.template_id
        WHERE ff.field_id = $1 AND ff.form_id = $2 AND ff.field_status <> $3;`

    err := f.client.QueryRow(
        ctx,
        sql,
        row.Id,
        row.FormId,
        entity.ItemStatusRemoved,
    ).Scan(
        &row.TemplateId,
        &row.Version,
        &row.CreatedAt,
        &row.ParamName,
        &row.Caption,
        &row.Type,
        &row.Detailing,
        &row.Body,
        &row.Required,
        &row.OrderField,
    )

    return err
}

// FetchIdByName
// uses: row{FormId, ParamName}
func (f *FormFieldItem) FetchIdByName(ctx context.Context, row *entity.FormFieldItem) (mrentity.KeyInt32, error) {
    sql := `
        SELECT field_id
        FROM
            public.form_fields
        WHERE form_id = $1 AND param_name = $2;`

    var id mrentity.KeyInt32

    err := f.client.QueryRow(ctx, sql, row.FormId, row.ParamName).Scan(
        &id,
    )

    return id, err
}

// Insert
// uses: row{FormId, TemplateId, ParamName, Caption, Required}
// modifies: row{Id}
func (f *FormFieldItem) Insert(ctx context.Context, row *entity.FormFieldItem) error {
    sql := `
        INSERT INTO public.form_fields
            (form_id,
             template_id,
             param_name,
             field_caption,
             field_required,
             field_status,
             order_field)
        VALUES
            ($1, $2, $3, $4, $5, $6, NULL)
        RETURNING field_id;`

    err := f.client.QueryRow(
        ctx,
        sql,
        row.FormId,
        row.TemplateId,
        row.ParamName,
        row.Caption,
        row.Required,
        entity.ItemStatusEnabled,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, FormId, Version, ParamName, Caption, Required}
func (f *FormFieldItem) Update(ctx context.Context, row *entity.FormFieldItem) error {
    sql := `
        UPDATE public.form_fields
        SET
            tag_version = tag_version + 1,
            param_name = $5,
            field_caption = $6,
            field_required = $7,
        WHERE field_id = $1 AND form_id = $2 AND tag_version = $3 AND field_status <> $4;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        row.Id,
        row.FormId,
        row.Version,
        entity.ItemStatusRemoved,
        row.ParamName,
        row.Caption,
        row.Required,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrerr.ErrStorageRowsNotAffected.New()
    }

    return nil
}

func (f *FormFieldItem) Delete(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) error {
    sql := `
        UPDATE public.form_fields
        SET
            tag_version = tag_version + 1,
            param_name = NULL,
            order_field = NULL,
            next_field_id = NULL,
            field_status = $3
        WHERE
            field_id = $1 AND form_id = $2 AND field_status <> $3;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        id,
        formId,
        entity.ItemStatusRemoved,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrerr.ErrStorageRowsNotAffected.New()
    }

    return nil
}
