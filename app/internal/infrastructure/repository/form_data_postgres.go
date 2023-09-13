package repository

import (
    "context"
    "print-shop-back/internal/entity"

    "github.com/Masterminds/squirrel"
    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrpostgres"
    "github.com/mondegor/go-webcore/mrcore"
)

type FormData struct {
    client *mrpostgres.ConnAdapter
    builder squirrel.StatementBuilderType
}

func NewFormData(client *mrpostgres.ConnAdapter,
                 queryBuilder squirrel.StatementBuilderType) *FormData {
    return &FormData{
        client: client,
        builder: queryBuilder,
    }
}

func (re *FormData) LoadAll(ctx context.Context, listFilter *entity.FormDataListFilter, rows *[]entity.FormData) error {
    query := re.builder.
        Select(`
            form_id,
            tag_version,
            datetime_created,
            param_name,
            form_caption,
            form_detailing,
            form_status`).
        From("public.form_data").
        Where(squirrel.NotEq{"form_status": mrcom.ItemStatusRemoved}).
        OrderBy("form_caption ASC, form_id ASC")

    if len(listFilter.Detailing) > 0 {
        query = query.Where(squirrel.Eq{"form_detailing": listFilter.Detailing})
    }

    if len(listFilter.Statuses) > 0 {
        query = query.Where(squirrel.Eq{"form_status": listFilter.Statuses})
    }

    cursor, err := re.client.SqQuery(ctx, query)

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
            return mrcore.FactoryErrStorageFetchDataFailed.Wrap(err)
        }

        *rows = append(*rows, row)
    }

    if err = cursor.Err(); err != nil {
        return mrcore.FactoryErrStorageFetchDataFailed.Wrap(err)
    }

    return nil
}

// LoadOne
// uses: row{Id}
// modifies: row{Version, CreatedAt, Caption, Detailing, Body, Status}
func (re *FormData) LoadOne(ctx context.Context, row *entity.FormData) error {
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

    return re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        mrcom.ItemStatusRemoved,
    ).Scan(
        &row.Version,
        &row.CreatedAt,
        &row.ParamName,
        &row.Caption,
        &row.Detailing,
        &row.Body,
        &row.Status,
    )
}

func (re *FormData) FetchIdByName(ctx context.Context, paramName string) (mrentity.KeyInt32, error) {
    sql := `
        SELECT form_id
        FROM
            public.form_fields
        WHERE param_name = $1;`

    var id mrentity.KeyInt32

    err := re.client.QueryRow(
        ctx,
        sql,
        paramName,
    ).Scan(
        &id,
    )

    return id, err
}

// FetchStatus
// uses: row{Id, Version}
func (re *FormData) FetchStatus(ctx context.Context, row *entity.FormData) (mrcom.ItemStatus, error) {
    sql := `
        SELECT form_status
        FROM
            public.form_data
        WHERE form_id = $1 AND tag_version = $2 AND form_status <> $3;`

    var status mrcom.ItemStatus

    err := re.client.QueryRow(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
    ).Scan(
        &status,
    )

    return status, err
}

// IsExists
// result: nil - exists, ErrStorageNoRowFound - not exists, error - query error
func (re *FormData) IsExists(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        SELECT 1
        FROM
            public.form_data
        WHERE form_id = $1 AND form_status <> $2;`

    return re.client.QueryRow(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    ).Scan(
        &id,
    )
}

// Insert
// uses: row{Caption, Detailing, Status}
// modifies: row{Id}
func (re *FormData) Insert(ctx context.Context, row *entity.FormData) error {
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

    err := re.client.QueryRow(
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
func (re *FormData) Update(ctx context.Context, row *entity.FormData) error {
    filledFields, err := mrentity.FilledFieldsToUpdate(row)

    if err != nil {
        return err
    }

    query := re.builder.
        Update("public.form_data").
        Set("tag_version", squirrel.Expr("tag_version + 1")).
        SetMap(filledFields).
        Where(squirrel.Eq{"form_id": row.Id}).
        Where(squirrel.Eq{"tag_version": row.Version}).
        Where(squirrel.NotEq{"form_status": mrcom.ItemStatusRemoved})

    return re.client.SqUpdate(ctx, query)
}

// UpdateStatus
// uses: row{Id, Version, Status}
func (re *FormData) UpdateStatus(ctx context.Context, row *entity.FormData) error {
    sql := `
        UPDATE public.form_data
        SET
            tag_version = tag_version + 1,
            form_status = $4
        WHERE
            form_id = $1 AND tag_version = $2 AND form_status <> $3;`

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        row.Id,
        row.Version,
        mrcom.ItemStatusRemoved,
        row.Status,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}

func (re *FormData) Delete(ctx context.Context, id mrentity.KeyInt32) error {
    sql := `
        UPDATE public.form_data
        SET
            tag_version = tag_version + 1,
            param_name = NULL,
            form_status = $2
        WHERE
            form_id = $1 AND form_status <> $2;`

    commandTag, err := re.client.Exec(
        ctx,
        sql,
        id,
        mrcom.ItemStatusRemoved,
    )

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() < 1 {
        return mrcore.FactoryErrStorageRowsNotAffected.New()
    }

    return nil
}
