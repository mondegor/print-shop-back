package repository

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"

    "github.com/Masterminds/squirrel"
)

type FormFieldItem struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewFormFieldItem(client *mrpostgres.Connection,
                      queryBuilder squirrel.StatementBuilderType) *FormFieldItem {
    return &FormFieldItem{
        client: client,
        builder: queryBuilder,
    }
}

func (re *FormFieldItem) GetMetaData(formId mrentity.KeyInt32) usecase.ItemMetaData {
    return NewItemMetaData(
        "public.form_fields",
        "field_id",
        []Condition{
            squirrel.Eq{"form_id": formId},
            squirrel.NotEq{"field_status": entity.ItemStatusRemoved},
        },
    )
}

func (re *FormFieldItem) LoadAll(ctx context.Context, listFilter *entity.FormFieldItemListFilter, rows *[]entity.FormFieldItem) error {
    query := re.builder.
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
            ff.field_required`).
        From("public.form_fields ff").
        Join("public.form_field_templates fft ON ff.template_id = fft.template_id").
        Where(squirrel.Eq{"ff.form_id": listFilter.FormId}).
        Where(squirrel.NotEq{"ff.field_status": entity.ItemStatusRemoved}).
        OrderBy("ff.order_field ASC, ff.field_caption ASC, ff.field_id ASC")

    if len(listFilter.Detailing) > 0 {
        query = query.Where(squirrel.Eq{"fft.field_detailing": listFilter.Detailing})
    }

    cursor, err := re.client.SqQuery(ctx, query)

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
// modifies: row{TemplateId, Version, CreatedAt, ParamName, Caption, Required}
func (re *FormFieldItem) LoadOne(ctx context.Context, row *entity.FormFieldItem) error {
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
            ff.field_required
        FROM
            public.form_fields ff
        JOIN
            public.form_field_templates fft
        ON
            ff.template_id = fft.template_id
        WHERE ff.field_id = $1 AND ff.form_id = $2 AND ff.field_status <> $3;`

    err := re.client.QueryRow(
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
    )

    return err
}

func (re *FormFieldItem) FetchIdByName(ctx context.Context, formId mrentity.KeyInt32, paramName string) (mrentity.KeyInt32, error) {
    sql := `
        SELECT field_id
        FROM
            public.form_fields
        WHERE form_id = $1 AND param_name = $2;`

    var id mrentity.KeyInt32

    err := re.client.QueryRow(
        ctx,
        sql,
        formId,
        paramName,
    ).Scan(
        &id,
    )

    return id, err
}

// Insert
// uses: row{FormId, TemplateId, ParamName, Caption, Required}
// modifies: row{Id}
func (re *FormFieldItem) Insert(ctx context.Context, row *entity.FormFieldItem) error {
    sql := `
        INSERT INTO public.form_fields
            (form_id,
             template_id,
             param_name,
             field_caption,
             field_required,
             field_status)
        VALUES
            ($1, $2, $3, $4, $5, $6)
        RETURNING field_id;`

    err := re.client.QueryRow(
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
func (re *FormFieldItem) Update(ctx context.Context, row *entity.FormFieldItem) error {
    filledFields, err := mrentity.GetFilledFieldsToUpdate(row)

    if err != nil {
        if !mrentity.ErrInternalListOfFieldsIsEmpty.Is(err) {
            return err
        }
    }

    filledFields["field_required"] = row.Required

    query := re.builder.
        Update("public.form_fields").
        Set("tag_version", squirrel.Expr("tag_version + 1")).
        SetMap(filledFields).
        Where(squirrel.Eq{"field_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"field_status": entity.ItemStatusRemoved})

    return re.client.SqUpdate(ctx, query)
}

func (re *FormFieldItem) Delete(ctx context.Context, id mrentity.KeyInt32, formId mrentity.KeyInt32) error {
    sql := `
        UPDATE public.form_fields
        SET
            tag_version = tag_version + 1,
            param_name = NULL,
            prev_field_id = NULL,
            next_field_id = NULL,
            order_field = NULL,
            field_status = $3
        WHERE
            field_id = $1 AND form_id = $2 AND field_status <> $3;`

    commandTag, err := re.client.Exec(
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
