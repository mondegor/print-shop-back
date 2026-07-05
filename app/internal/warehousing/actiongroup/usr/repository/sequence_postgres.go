package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-core/mrstorage"

	"print-shop-back/internal/warehousing/module"
)

type (
	// AccountSequencePostgres - comment struct.
	AccountSequencePostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewAccountSequencePostgres - создаёт объект AccountSequencePostgres.
func NewAccountSequencePostgres(client mrstorage.DBConnManager) *AccountSequencePostgres {
	return &AccountSequencePostgres{
		client: client,
	}
}

// NextContainerID - comment method.
func (re *AccountSequencePostgres) NextContainerID(ctx context.Context, accountID uuid.UUID) (id uint64, err error) {
	sql := `
		INSERT INTO ` + module.DBTableNameAccountSequences + ` as t1
			(
				account_id
			)
		VALUES
			($1)
		ON CONFLICT (account_id) DO UPDATE
		SET
			last_container_id = t1.last_container_id + 1
		RETURNING
			t1.last_container_id;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		accountID,
	).Scan(
		&id,
	)

	return id, err
}
