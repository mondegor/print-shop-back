package repository

import (
	"context"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrstorage"

	"print-shop-back/internal/warehousing/actiongroup/back/dto"
	"print-shop-back/internal/warehousing/module"
)

type (
	// ContainerPostgres - comment struct.
	ContainerPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewContainerPostgres - создаёт объект ContainerPostgres.
func NewContainerPostgres(client mrstorage.DBConnManager) *ContainerPostgres {
	return &ContainerPostgres{
		client: client,
	}
}

// FetchGroupingContainersByIDs - comment method.
func (re *ContainerPostgres) FetchGroupingContainersByIDs(ctx context.Context, rowIDs []uint64) (rows []dto.GroupingContainer, err error) {
	sql := `
		SELECT
			container_id,
			container_code,
			container_images
		FROM
			` + module.DBTableNameContainers + `
		WHERE
			container_id = ANY($1);`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		rowIDs,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for cursor.Next() {
		var row dto.GroupingContainer

		err = cursor.Scan(
			&row.ID,
			&row.Code,
			&row.Images,
		)
		if err != nil {
			return nil, err
		}

		if rows == nil {
			rows = make([]dto.GroupingContainer, 0, len(rowIDs))
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

// InsertToHistoryByIDs - переносит информацию о контейнерах в холодное хранилище.
// TODO: заюзать.
func (re *ContainerPostgres) InsertToHistoryByIDs(ctx context.Context, containerID []uint64) error {
	sql := `
		INSERT INTO ` + module.DBTableNameContainers + `_history
			(
				container_id,
				account_id,
				container_code,
				container_marker,
				container_volume,
				container_tags,
				container_images,
				created_at,
				updated_at
			)
		SELECT
			container_id,
			account_id,
			container_code,
			container_marker,
			container_volume,
			container_tags,
			container_images,
			created_at,
			updated_at
		FROM
			` + module.DBTableNameContainers + `
		WHERE
			container_id = ANY($1)
		ON CONFLICT (container_id) DO NOTHING;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		containerID,
	)
}

// UpdateGroups - обновляет информацию о групповых контейнерах.
func (re *ContainerPostgres) UpdateGroups(ctx context.Context, rows []dto.UpdateGroupContainer) error {
	if len(rows) == 0 {
		return nil
	}

	containerIDs := make([]uint64, 0, len(rows))
	tags := make([][]string, 0, len(rows))
	images := make([][]string, 0, len(rows))

	for _, row := range rows {
		containerIDs = append(containerIDs, row.ID)
		tags = append(tags, row.Tags)
		images = append(images, row.Images)
	}

	sql := `
		UPDATE
			` + module.DBTableNameContainers + ` c
		SET
			tag_version = tag_version + 1,
			container_tags = cn.container_tags,
			container_images = cn.container_images,
			updated_at = NOW()
		FROM
			(
				SELECT *
				FROM
					UNNEST($1::int8[], $2::jsonb[], $3::jsonb[])
					as t(container_id, container_tags, container_images)
	  		) cn
		WHERE
			c.container_id = cn.container_id;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		containerIDs,
		tags,
		images,
	)
}

// DeleteByIDs - удаляет информацию о контейнерах из основного хранилища.
// TODO: заюзать.
func (re *ContainerPostgres) DeleteByIDs(ctx context.Context, containerID []uint64) error {
	sql := `
		DELETE FROM
			` + module.DBTableNameContainers + `
		WHERE
			container_id = ANY($1);`

	err := re.client.Conn(ctx).Exec(
		ctx,
		sql,
		containerID,
	)
	// если это внутренняя ошибка
	if err != nil && !errors.Is(err, errors.ErrEventStorageRecordsNotAffected) {
		return err
	}

	return nil
}
