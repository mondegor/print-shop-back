package repository

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// FormVersionPostgres - comment struct.
	// FormVersionPostgres - comment struct.
	FormVersionPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewFormVersionPostgres - создаёт объект FormVersionPostgres.
func NewFormVersionPostgres(client mrstorage.DBConnManager) *FormVersionPostgres {
	return &FormVersionPostgres{
		client: client,
	}
}

// Fetch - comment method.
func (re *FormVersionPostgres) Fetch(ctx context.Context, formID uuid.UUID) ([]entity.FormVersion, error) {
	sql := `
        SELECT
            version,
            rewrite_name,
			form_caption,
            form_detailing,
            activity_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
        WHERE
            form_id = $1
        ORDER BY
            version ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		formID,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.FormVersion, 0)

	for cursor.Next() {
		var (
			rewriteName *string
			row         entity.FormVersion
		)

		err = cursor.Scan(
			&row.Version,
			&rewriteName,
			&row.Caption,
			&row.Detailing,
			&row.ActivityStatus,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if rewriteName != nil {
			row.RewriteName = *rewriteName
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

// FetchOne - comment method.
func (re *FormVersionPostgres) FetchOne(ctx context.Context, primary entity.PrimaryKey) (entity.FormVersion, error) {
	sql := `
        SELECT
            rewrite_name,
            form_caption,
            form_detailing,
            compiled_body,
            activity_status,
            created_at,
			updated_at
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
        WHERE
            form_id = $1 AND version = $2
        LIMIT 1;`

	var rewriteName *string

	row := entity.FormVersion{
		ID:      primary.FormID,
		Version: primary.Version,
	}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		primary.FormID,
		primary.Version,
	).Scan(
		&rewriteName,
		&row.Caption,
		&row.Detailing,
		&row.Body,
		&row.ActivityStatus,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if rewriteName != nil {
		row.RewriteName = *rewriteName
	}

	return row, err
}

// FetchOneLastVersion - comment method.
func (re *FormVersionPostgres) FetchOneLastVersion(ctx context.Context, formID uuid.UUID) (entity.FormVersionStatus, error) {
	sql := `
		SELECT
			version,
			activity_status
		FROM
			` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
		WHERE
			form_id = $1
		ORDER BY
			version DESC
		LIMIT 1;`

	row := entity.FormVersionStatus{FormID: formID}

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		formID,
	).Scan(
		&row.Version,
		&row.ActivityStatus,
	)

	return row, err
}

// Insert - comment method.
func (re *FormVersionPostgres) Insert(ctx context.Context, row entity.FormVersion) error {
	sql := `
		INSERT INTO ` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
			(
				form_id,
				version,
				rewrite_name,
				form_caption,
				form_detailing,
				compiled_body,
				activity_status
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7);`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		row.ID,
		row.Version,
		row.RewriteName,
		row.Caption,
		row.Detailing,
		row.Body,
		row.ActivityStatus,
	)
}

// Update - comment method.
func (re *FormVersionPostgres) Update(ctx context.Context, row entity.FormVersion) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
        SET
			updated_at = NOW(),
			rewrite_name = $4,
			form_caption = $5,
			form_detailing = $6,
			compiled_body = $7
        WHERE
            form_id = $1 AND version = $2 AND activity_status = $3;`

	return re.client.Conn(ctx).Exec(
		ctx,
		sql,
		row.ID,
		row.Version,
		row.ActivityStatus,
		row.RewriteName,
		row.Caption,
		row.Detailing,
		row.Body,
	)
}

// UpdateStatus - comment method.
func (re *FormVersionPostgres) UpdateStatus(ctx context.Context, row entity.FormVersionStatus, toStatus enum.ActivityStatus) error {
	return re.client.Do(ctx, func(ctx context.Context) error {
		conn := re.client.Conn(ctx)

		// архивирование старой версии со статусом toStatus, если она имеется
		sql := `
			UPDATE
				` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
			SET
				rewrite_name = NULL,
				activity_status = $4,
				updated_at = NOW()
			WHERE
				form_id = $1 AND version < $2 AND activity_status = $3;`

		err := conn.Exec(
			ctx,
			sql,
			row.FormID,
			row.Version,
			toStatus,
			enum.ActivityStatusArchived,
		)
		if err != nil {
			// если это системная ошибка
			if !mrcore.ErrStorageRowsNotAffected.Is(err) {
				return err
			}
		}

		// переключение только указанной версии в указанный статус
		sql = `
			UPDATE
				` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
			SET
				activity_status = $4,
				updated_at = NOW()
			WHERE
				form_id = $1 AND version = $2 AND activity_status = $3;`

		return conn.Exec(
			ctx,
			sql,
			row.FormID,
			row.Version,
			row.ActivityStatus,
			toStatus,
		)
	})
}
