package repository

import (
	"context"

	"github.com/mondegor/go-core/mrstorage"

	"print-shop-back/internal/warehousing/actiongroup/back/dto"
	"print-shop-back/internal/warehousing/module"
)

type (
	// StorePostgres - comment struct.
	StorePostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewStorePostgres - создаёт объект StorePostgres.
func NewStorePostgres(client mrstorage.DBConnManager) *StorePostgres {
	return &StorePostgres{
		client: client,
	}
}

// UpdateContainersVolume - comment method.
func (re *StorePostgres) UpdateContainersVolume(ctx context.Context, rows []dto.LocationContainersVolume) error {
	if len(rows) == 0 {
		return nil
	}

	storeIDs := make([]uint64, 0, len(rows))
	totalVolumes := make([]float64, 0, len(rows))

	for _, row := range rows {
		storeIDs = append(storeIDs, row.LocationID)
		totalVolumes = append(totalVolumes, row.TotalVolume)
	}

	sql := `
		UPDATE
			` + module.DBTableNameStores + ` s
		SET
			containers_volume = sn.containers_volume,
			updated_at = NOW()
		FROM
		  	(
				SELECT *
				FROM
					UNNEST($1::int8[], $2::double precision[])
					as t(store_id, containers_volume)
	  	  	) sn
		WHERE
			s.store_id = sn.store_id;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		storeIDs,
		totalVolumes,
	)
}
