package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"

	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/module"
	"github.com/mondegor/print-shop-back/internal/calculations/queryhistory/section/pub/entity"
)

type (
	// QueryHistoryPostgres - comment struct.
	QueryHistoryPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewQueryHistoryPostgres - создаёт объект QueryHistoryPostgres.
func NewQueryHistoryPostgres(client mrstorage.DBConnManager) *QueryHistoryPostgres {
	return &QueryHistoryPostgres{
		client: client,
	}
}

// FetchOne - comment method.
func (re *QueryHistoryPostgres) FetchOne(ctx context.Context, rowID uuid.UUID) (entity.QueryHistoryItem, error) {
	sql := `
		SELECT
			query_caption,
			query_params,
			query_result,
			created_at
		FROM
			` + module.DBSchema + `.` + module.DBTableNameQueryHistory + `
		WHERE
			query_id = $1 AND deleted_at IS NULL
		LIMIT 1;`

	row := entity.QueryHistoryItem{}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
	).Scan(
		&row.Caption,
		&row.Params,
		&row.Result,
		&row.CreatedAt,
	)

	return row, err
}

// Insert - comment method.
func (re *QueryHistoryPostgres) Insert(ctx context.Context, row entity.QueryHistoryItem) (id uuid.UUID, err error) {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameQueryHistory + `
			(
				query_id,
				query_caption,
				query_params,
				query_result
			)
		VALUES
			(gen_random_uuid(), $1, $2, $3)
		RETURNING
			query_id;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		row.Caption,
		row.Params,
		row.Result,
	).Scan(
		&id,
	)

	return id, err
}

// UpdateQuantity - comment method.
func (re *QueryHistoryPostgres) UpdateQuantity(ctx context.Context, rowID uuid.UUID) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameQueryHistory + `
		SET
			used_count = used_count + 1,
			used_at = NOW()
		WHERE
			query_id = $1 AND deleted_at IS NULL
		LIMIT $2;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		rowID,
	)
}

// Delete - comment method.
func (re *QueryHistoryPostgres) Delete(ctx context.Context, expiry time.Duration, limit uint64) error {
	sql := `
		UPDATE
			` + module.DBSchema + `.` + module.DBTableNameQueryHistory + `
		SET
			deleted_at = NOW()
		WHERE
			created_at < NOW() - INTERVAL $1 SECOND AND deleted_at IS NULL
		LIMIT $2;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		expiry.Nanoseconds()/1e9, // to seconds
		limit,
	)
}
