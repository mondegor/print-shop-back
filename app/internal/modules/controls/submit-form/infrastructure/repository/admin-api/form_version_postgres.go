package repository

import (
	"context"
	module "print-shop-back/internal/modules/controls/submit-form"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"
	"print-shop-back/pkg/modules/controls/enums"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

type (
	FormVersionPostgres struct {
		client mrstorage.DBConn
	}
)

func NewFormVersionPostgres(
	client mrstorage.DBConn,
) *FormVersionPostgres {
	return &FormVersionPostgres{
		client: client,
	}
}

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
            ` + module.DBSchema + `.submit_form_versions
        WHERE
            form_id = $1
        ORDER BY
            version ASC;`

	cursor, err := re.client.Query(
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
		var rewriteName *string
		row := entity.FormVersion{}

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
            ` + module.DBSchema + `.submit_form_versions
        WHERE
            form_id = $1 AND version = $2
        LIMIT 1;`

	var rewriteName *string
	row := entity.FormVersion{
		ID:      primary.FormID,
		Version: primary.Version,
	}

	err := re.client.QueryRow(
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

func (re *FormVersionPostgres) FetchOneLastVersion(ctx context.Context, formID uuid.UUID) (entity.FormVersionStatus, error) {
	sql := `
		SELECT
			version,
			activity_status
		FROM
			` + module.DBSchema + `.submit_form_versions
		WHERE
			form_id = $1
		ORDER BY
			version DESC
		LIMIT 1;`

	row := entity.FormVersionStatus{FormID: formID}

	err := re.client.QueryRow(
		ctx,
		sql,
		formID,
	).Scan(
		&row.Version,
		&row.ActivityStatus,
	)

	return row, err
}

func (re *FormVersionPostgres) Insert(ctx context.Context, row entity.FormVersion) error {
	sql := `
		INSERT INTO ` + module.DBSchema + `.submit_form_versions
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

	return re.client.Exec(
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

func (re *FormVersionPostgres) Update(ctx context.Context, row entity.FormVersion) error {
	sql := `
        UPDATE
            ` + module.DBSchema + `.submit_form_versions
        SET
			updated_at = NOW(),
			rewrite_name = $4,
			form_caption = $5,
			form_detailing = $6,
			compiled_body = $7
        WHERE
            form_id = $1 AND version = $2 AND activity_status = $3;`

	return re.client.Exec(
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

func (re *FormVersionPostgres) UpdateStatus(ctx context.Context, row entity.FormVersionStatus, toStatus enums.ActivityStatus) error {
	tx, err := re.client.Begin(ctx)

	if err != nil {
		return err // :TODO:
	}

	defer func() {
		var e error

		if err != nil {
			e = tx.Rollback(ctx)
		} else {
			e = tx.Commit(ctx)
		}

		if e != nil {
			mrlog.Ctx(ctx).Error().Caller().Err(e).Msg("tx defer func()")
		}
	}()

	// архивирование старой версии со статусом toStatus, если она имеется
	sql := `
		UPDATE
			` + module.DBSchema + `.submit_form_versions
		SET
			rewrite_name = NULL,
			activity_status = $4,
			updated_at = NOW()
		WHERE
			form_id = $1 AND version < $2 AND activity_status = $3;`

	err = tx.Exec(
		ctx,
		sql,
		row.FormID,
		row.Version,
		toStatus,
		enums.ActivityStatusArchived,
	)

	// если это системная ошибка
	if err != nil && !mrcore.FactoryErrStorageRowsNotAffected.Is(err) {
		return err
	}

	// переключение только указанной версии в указанный статус
	sql = `
		UPDATE
			` + module.DBSchema + `.submit_form_versions
		SET
			activity_status = $4,
			updated_at = NOW()
		WHERE
			form_id = $1 AND version = $2 AND activity_status = $3;`

	// WARNING: err is used in defer func()
	err = tx.Exec(
		ctx,
		sql,
		row.FormID,
		row.Version,
		row.ActivityStatus,
		toStatus,
	)

	return err
}
