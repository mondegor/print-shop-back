package repository

import (
    "calc-user-data-back-adm/internal/entity"
    "calc-user-data-back-adm/pkg/client/mrpostgres"
    "calc-user-data-back-adm/pkg/mrentity"
    "calc-user-data-back-adm/pkg/mrerr"
    "context"
    "fmt"

    "github.com/Masterminds/squirrel"
)

type FormFieldItemOrderer struct {
    client *mrpostgres.Connection
    builder squirrel.StatementBuilderType
}

func NewFormFieldItemOrderer(client *mrpostgres.Connection, queryBuilder squirrel.StatementBuilderType) *FormFieldItemOrderer {
    return &FormFieldItemOrderer{
        client: client,
        builder: queryBuilder,
    }
}

func (f *FormFieldItemOrderer) LoadNode(ctx context.Context, row *entity.Node) error {
    sql := `
        SELECT prev_field_id, next_field_id, order_field
        FROM public.form_fields_test
        WHERE field_id = $1;`

    return f.client.QueryRow(ctx, sql, row.Id).Scan(&row.PrevId, &row.NextId, &row.OrderField)
}

func (f *FormFieldItemOrderer) LoadFirstNode(ctx context.Context, row *entity.Node) error {
    sql := `
        SELECT MIN(order_field)
        FROM public.form_fields_test;`

    err := f.client.QueryRow(ctx, sql).Scan(&row.OrderField)

    if err != nil {
        return err
    }

    err = f.loadNodeByOrderField(ctx, row)

    if row.PrevId > 0 {
        return mrerr.ErrStorageFetchDataFailed.NewWithData("row.Id=%d; row.PrevId=%d", row.Id, row.PrevId)
    }

    return nil
}

func (f *FormFieldItemOrderer) LoadLastNode(ctx context.Context, row *entity.Node) error {
    sql := `
        SELECT MAX(order_field)
        FROM public.form_fields_test;`

    err := f.client.QueryRow(ctx, sql).Scan(&row.OrderField)

    if err != nil {
        return err
    }

    err = f.loadNodeByOrderField(ctx, row)

    if row.NextId > 0 {
        return mrerr.ErrStorageFetchDataFailed.NewWithData("row.Id=%d; row.NextId=%d", row.Id, row.NextId)
    }

    return nil
}

func (f *FormFieldItemOrderer) UpdateNode(ctx context.Context, row *entity.Node) error {
    sql := `
        UPDATE public.form_fields_test
        SET
            prev_field_id = $2,
            next_field_id = $3,
            order_field = $4
        WHERE field_id = $1;`

    commandTag, err := f.client.Exec(ctx, sql, row.Id, row.PrevId, row.NextId, row.OrderField)

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() != 1 {
        return fmt.Errorf("field_id = %d not updated (full)", row.Id)
    }

    return nil
}

func (f *FormFieldItemOrderer) UpdateNodePrevId(ctx context.Context, id mrentity.KeyInt32, prevId mrentity.ZeronullInt32) error {
    sql := `
        UPDATE public.form_fields_test
        SET
            prev_field_id = $2
        WHERE field_id = $1;`

    commandTag, err := f.client.Exec(ctx, sql, id, prevId)

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() != 1 {
        return fmt.Errorf("field_id = %d not updated (prev)", id)
    }

    return nil
}

func (f *FormFieldItemOrderer) UpdateNodeNextId(ctx context.Context, id mrentity.KeyInt32, nextId mrentity.ZeronullInt32) error {
    sql := `
        UPDATE public.form_fields_test
        SET
            next_field_id = $2
        WHERE field_id = $1;`

    commandTag, err := f.client.Exec(ctx, sql, id, nextId)

    if err != nil {
        return err
    }

    if commandTag.RowsAffected() != 1 {
        return fmt.Errorf("field_id = %d not updated (next)", nextId)
    }

    return nil
}

func (f *FormFieldItemOrderer) RecalcOrderField(ctx context.Context, minBorder mrentity.ZeronullInt64, step mrentity.ZeronullInt64) error {
    sql := `
        UPDATE public.form_fields_test
        SET
            order_field = order_field + $2
        WHERE order_field > $1;`

    _, err := f.client.Exec(ctx, sql, minBorder, step)

    return err
}

func (f *FormFieldItemOrderer) loadNodeByOrderField(ctx context.Context, row *entity.Node) error {
    sql := `
        SELECT field_id, prev_field_id, next_field_id
        FROM public.form_fields_test
        WHERE order_field = $1
        FETCH FIRST 1 ROWS ONLY;`

    return f.client.QueryRow(ctx, sql, row.OrderField).Scan(&row.Id, &row.PrevId, &row.NextId)
}
