package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-sysmess/mrstorage"

	"print-shop-back/internal/warehousing/actiongroup/usr/entity"
	"print-shop-back/internal/warehousing/module"
)

type (
	// TerritoryPostgres - comment struct.
	TerritoryPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewTerritoryPostgres - создаёт объект TerritoryPostgres.
func NewTerritoryPostgres(client mrstorage.DBConnManager) *TerritoryPostgres {
	return &TerritoryPostgres{
		client: client,
	}
}

// FetchState - comment method.
func (re *TerritoryPostgres) FetchState(ctx context.Context, accountID uuid.UUID, rowID uint64) (row entity.TerritoryState, err error) {
	sql := `
		SELECT
			territory_id,
			tag_version,
			code_pattern
		FROM
			` + module.DBTableNameTerritories + `
		WHERE
			territory_id = $1 AND account_id = $2 AND deleted_at IS NULL
		FETCH FIRST 1 ROW ONLY;`

	err = re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rowID,
		accountID,
	).Scan(
		&row.ID,
		&row.TagVersion,
		&row.CodePattern,
	)

	return row, err
}
