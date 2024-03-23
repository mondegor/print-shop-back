package repository

import (
	"context"
	module "print-shop-back/internal/modules/controls/submit-form"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/public-api"
	"print-shop-back/pkg/modules/controls/enums"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	SubmitFormPostgres struct {
		client mrstorage.DBConn
	}
)

func NewSubmitFormPostgres(
	client mrstorage.DBConn,
) *SubmitFormPostgres {
	return &SubmitFormPostgres{
		client: client,
	}
}

func (re *SubmitFormPostgres) Fetch(ctx context.Context, params entity.SubmitFormParams) ([]entity.SubmitForm, error) {
	sql := `
        SELECT
            version,
			rewrite_name,
			form_caption
        FROM
            ` + module.DBSchema + `.submit_form_versions
        WHERE
            activity_status = $1 AND form_detailing = $2
        ORDER BY
            form_caption ASC, created_at ASC;`

	cursor, err := re.client.Query(
		ctx,
		sql,
		enums.ActivityStatusPublished,
		enums.ElementDetailingNormal,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.SubmitForm, 0)

	for cursor.Next() {
		var row entity.SubmitForm

		err = cursor.Scan(
			&row.Version,
			&row.RewriteName,
			&row.Caption,
		)

		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}

func (re *SubmitFormPostgres) FetchByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error) {
	sql := `
        SELECT
            version,
			rewrite_name,
			form_caption
        FROM
            ` + module.DBSchema + `.submit_form_versions
        WHERE
            rewrite_name = $1 AND activity_status = $2 AND form_detailing = $3
        LIMIT 1;`

	var row entity.SubmitForm

	err := re.client.QueryRow(
		ctx,
		sql,
		rewriteName,
		enums.ActivityStatusPublished,
		enums.ElementDetailingNormal,
	).Scan(
		&row.Version,
		&row.RewriteName,
		&row.Caption,
	)

	return row, err
}
