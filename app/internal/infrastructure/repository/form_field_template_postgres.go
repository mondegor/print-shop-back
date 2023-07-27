package repository

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"

    "github.com/Masterminds/squirrel"
)

type FormFieldTemplate struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewFormFieldTemplate(client *mrpostgres.Connection,
                          queryBuilder squirrel.StatementBuilderType) *FormFieldTemplate {
    return &FormFieldTemplate{
        client: client,
        builder: queryBuilder,
    }
}

func (re *FormFieldTemplate) LoadAll(ctx context.Context, listFilter *entity.FormFieldTemplateListFilter, rows *[]entity.FormFieldTemplate) error {
    query := re.builder.
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
        query = query.Where(squirrel.Eq{"template_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

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
func (re *FormFieldTemplate) LoadOne(ctx context.Context, row *entity.FormFieldTemplate) error {
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

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        entity.ItemStatusRemoved,
    ).Scan(
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
func (re *FormFieldTemplate) FetchStatus(ctx context.Context, row *entity.FormFieldTemplate) (entity.ItemStatus, error) {
    sql := `
        SELECT template_status
        FROM
            public.form_field_templates
        WHERE template_id = $1 AND tag_version = $2 AND template_status <> $3;`

    var status entity.ItemStatus

    err := re.client.QueryRow(
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
func (re *FormFieldTemplate) Insert(ctx context.Context, row *entity.FormFieldTemplate) error {
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
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, ParamName, Caption, Type, Detailing, Body}
func (re *FormFieldTemplate) Update(ctx context.Context, row *entity.FormFieldTemplate) error {
    filledFields, err := mrentity.GetFilledFieldsToUpdate(row)

    if err != nil {
        return err
    }

    query := re.builder.
        Update("public.form_field_templates").
        Set("tag_version", squirrel.Expr("tag_version + 1")).
        SetMap(filledFields).
        Where(squirrel.Eq{"template_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"template_status": entity.ItemStatusRemoved})

    return re.client.SqUpdate(ctx, query)
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *FormFieldTemplate) UpdateStatus(ctx context.Context, row *entity.FormFieldTemplate) error {
    sql := `
        UPDATE public.form_field_templates
        SET
            tag_version = tag_version + 1,
            template_status = $4
        WHERE
            template_id = $1 AND tag_version = $2 AND template_status <> $3;`

    commandTag, err := re.client.Exec(
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

func (re *FormFieldTemplate) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.form_field_templates
        SET
            tag_version = tag_version + 1,
            param_name = NULL,
            template_status = $2
        WHERE
            template_id = $1 AND template_status <> $2;`

    commandTag, err := re.client.Exec(
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
