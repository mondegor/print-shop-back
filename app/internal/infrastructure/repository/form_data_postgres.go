package repository

import (
    "context"
    "print-shop-back/internal/entity"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrentity"
    "print-shop-back/pkg/mrerr"

    "github.com/Masterminds/squirrel"
)

type FormData struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewFormData(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *FormData {
    return &FormData{
        client: client,
        builder: queryBuilder,
    }
}

func (f *FormData) LoadAll(ctx context.Context, listFilter *entity.FormDataListFilter, rows *[]entity.FormData) error {
    tbl := f.builder.
        Select(`
            form_id,
            tag_version,
            datetime_created,
            param_name,
            form_caption,
            form_detailing,
            form_status`).
        From("public.form_data").
        Where(squirrel.NotEq{"form_status": entity.ItemStatusRemoved}).
        OrderBy("form_caption ASC, form_id ASC")

    if len(listFilter.Detailing) > 0 {
        tbl = tbl.Where(squirrel.Eq{"form_detailing": listFilter.Detailing})
    }

    if len(listFilter.Statuses) > 0 {
        tbl = tbl.Where(squirrel.Eq{"form_status": listFilter.Statuses})
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
        var row entity.FormData

        err = cursor.Scan(
            &row.Id,
            &row.Version,
            &row.CreatedAt,
            &row.ParamName,
            &row.Caption,
            &row.Detailing,
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
// modifies: row{Version, CreatedAt, Caption, Detailing, Body, Status}
func (f *FormData) LoadOne(ctx context.Context, row *entity.FormData) error {
    sql := `
        SELECT
            tag_version,
            datetime_created,
            param_name,
            form_caption,
            form_detailing,
            form_body_compiled,
            form_status
        FROM
            public.form_data
        WHERE form_id = $1 AND form_status <> $2;`

    return f.client.QueryRow(ctx, sql, row.Id, entity.ItemStatusRemoved).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.ParamName,
        &row.Caption,
        &row.Detailing,
        &row.Body,
        &row.Status,
    )
}

// FetchIdByName
// uses: row{FormId, ParamName}
func (f *FormData) FetchIdByName(ctx context.Context, row *entity.FormData) (mrentity.KeyInt32, error) {
    sql := `
        SELECT form_id
        FROM
            public.form_fields
        WHERE param_name = $1;`

    var id mrentity.KeyInt32

    err := f.client.QueryRow(ctx, sql, row.ParamName).Scan(
        &id,
    )

    return id, err
}

// FetchStatus
// uses: row{Id, Version}
func (f *FormData) FetchStatus(ctx context.Context, row *entity.FormData) (entity.ItemStatus, error) {
    sql := `
        SELECT form_status
        FROM
            public.form_data
        WHERE form_id = $1 AND tag_version = $2 AND form_status <> $3;`

    var status entity.ItemStatus

    err := f.client.QueryRow(ctx, sql, row.Id, row.Version, entity.ItemStatusRemoved).Scan(
        &status,
    )

    return status, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (f *FormData) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.form_data
        WHERE form_id = $1 AND form_status <> $2;`

    return f.client.QueryRow(ctx, sql, id, entity.ItemStatusRemoved).Scan(&id)
}

// Insert
// uses: row{Caption, Detailing, Status}
// modifies: row{Id}
func (f *FormData) Insert(ctx context.Context, row *entity.FormData) error {
    sql := `
        INSERT INTO public.form_data
            (param_name,
             form_caption,
             form_detailing,
             form_body_compiled,
             form_status)
        VALUES
            ($1, $2, $3, '[]', $4)
        RETURNING form_id;`

    err := f.client.QueryRow(
        ctx,
        sql,
        row.ParamName,
        row.Caption,
        row.Detailing,
        row.Status,
    ).Scan(
        &row.Id,
    )

    return err
}

// Update
// uses: row{Id, Version, Caption, Detailing}
func (f *FormData) Update(ctx context.Context, row *entity.FormData) error {
    sql := `
        UPDATE public.form_data
        SET
            tag_version = tag_version + 1,
            param_name = $4,
            form_caption = $5,
            form_detailing = $6
        WHERE form_id = $1 AND tag_version = $2 AND form_status <> $3;`

    commandTag, err := f.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        entity.ItemStatusRemoved,
        row.ParamName,
        row.Caption,
        row.Detailing,
    )

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
func (f *FormData) UpdateStatus(ctx context.Context, row *entity.FormData) error {
    sql := `
        UPDATE public.form_data
        SET
            tag_version = tag_version + 1,
            form_status = $4
        WHERE
            form_id = $1 AND tag_version = $2 AND form_status <> $3;`

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

func (f *FormData) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.form_data
        SET
            tag_version = tag_version + 1,
            form_status = $2
        WHERE
            form_id = $1 AND form_status <> $2;`

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
