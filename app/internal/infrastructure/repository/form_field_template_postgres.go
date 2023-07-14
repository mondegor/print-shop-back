package repository

import (
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"
    "context"

    "github.com/Masterminds/squirrel"
)

type FormFieldTemplate struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewFormFieldTemplate(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *FormFieldTemplate {
    return &FormFieldTemplate{
        client: client,
        builder: queryBuilder,
    }
}

func (f *FormFieldTemplate) LoadAll(ctx context.Context, listFilter *entity.FormFieldTemplateListFilter, rows *[]entity.FormFieldTemplate) error {
    tbl := f.builder.
        Select(`
            template_id,
            tag_version,
            datetime_created,
            param_name,
            template_caption,
            field_type, field_detailing,
            field_body,
            template_status`).
        From("public.form_field_templates").
        Where(squirrel.NotEq{"template_status": entity.ItemStatusRemoved}).
        OrderBy("template_caption ASC, template_id ASC")

    if len(listFilter.Statuses) > 0 {
        tbl = tbl.Where(squirrel.Eq{"template_status": listFilter.Statuses})
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
        var row entity.FormFieldTemplate

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.ParamName,
            &row.Caption,
            &row.Type,
            &row.Detailing,
            &row.Body,
            &row.Status,
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
// uses: row{Id}
// modifies: row{Version, CreatedAt, ParamName, Caption, Type, Detailing, Body, Status}
func (f *FormFieldTemplate) LoadOne(ctx context.Context, row *entity.FormFieldTemplate) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            param_name,
            template_caption,
            field_type,
            field_detailing,
            field_body,
            template_status
        FROM
            public.form_field_templates
        WHERE template_id = $1 AND template_status <> $2;`

    return f.client.QueryRow(ctx, sql, row.Id, entity.ItemStatusRemoved).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.ParamName,
        &row.Caption,
        &row.Type,
        &row.Detailing,
        &row.Body,
        &row.Status,
    )
}

// FetchStatus
// uses: row{Id, Version}
func (f *FormFieldTemplate) FetchStatus(ctx context.Context, row *entity.FormFieldTemplate) (entity.ItemStatus, error) {
    sql := `
        SELECT template_status
        FROM
            public.form_field_templates
        WHERE template_id = $1 AND tag_version = $2 AND template_status <> $3;`

    var status entity.ItemStatus

    err := f.client.QueryRow(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
    ).Scan(
        &status,
    )

    return status, err
}

// Insert
// uses: row{ParamName, Caption, Type, Detailing, Body, Status}
// modifies: row{Id}
func (f *FormFieldTemplate) Insert(ctx context.Context, row *entity.FormFieldTemplate) error {
    sql := `
        INSERT INTO public.form_field_templates
            (param_name,
             template_caption,
             field_type,
             field_detailing,
             field_body,
             template_status)
        VALUES
            ($1, $2, $3, $4, $5, $6)
        RETURNING template_id;`

    err := f.client.QueryRow(
        ctx,
        sql,
        row.ParamName,
        row.Caption,
        row.Type,
        row.Detailing,
        row.Body,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, ParamName, Caption, Type, Detailing, Body}
func (f *FormFieldTemplate) Update(ctx context.Context, row *entity.FormFieldTemplate) error {
    tbl := f.builder.
        Update("public.form_field_templates").
        SetMap(map[string]any{
            "tag_version": squirrel.Expr("tag_version + 1"),
            "param_name": row.ParamName,
            "template_caption": row.Caption,
            "field_type": row.Type,
            "field_detailing": row.Detailing,
            "field_body": row.Body,
        }).
        Where(squirrel.Eq{"template_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"template_status": entity.ItemStatusRemoved})

    sql, args, err := tbl.ToSql()

    if err != nil {
        return mrerr.ErrInternal.Wrap(err)
    }

    commandTag, err := f.client.Exec(ctx, sql, args...)

    //sql := `
    //    UPDATE public.form_field_templates
    //    SET
    //        tag_version = tag_version + 1,
    //        param_name = $4,
    //        template_caption = $5,
    //        field_type = $6,
    //        field_detailing = $7,
    //        field_body = $8
    //    WHERE template_id = $1 AND tag_version = $2 AND template_status <> $3;`
    //
    //commandTag, err := f.client.Exec(
    //    ctx,
    //    sql,
    //    row.Id,
    //    row.Version,
    //    entity.ItemStatusRemoved,
    //    row.ParamName,
    //    row.Caption,
    //    row.Type,
    //    row.Detailing,
    //    row.Body,
    //)

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrerr.ErrStorageRowsNotAffected.New()
    }

    return nil
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (f *FormFieldTemplate) UpdateStatus(ctx context.Context, row *entity.FormFieldTemplate) error {
    sql := `
        UPDATE public.form_field_templates
        SET
            tag_version = tag_version + 1,
            template_status = $4
        WHERE
            template_id = $1 AND tag_version = $2 AND template_status <> $3;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
        row.Status,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrerr.ErrStorageRowsNotAffected.New()
    }

    return nil
}

func (f *FormFieldTemplate) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.form_field_templates
        SET
            tag_version = tag_version + 1,
            template_status = $2
        WHERE
            template_id = $1 AND template_status <> $2;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        id,
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
