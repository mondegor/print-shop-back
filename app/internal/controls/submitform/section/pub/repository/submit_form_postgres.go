package repository

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/module"
	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/controls/enum"
)

type (
	// SubmitFormPostgres - comment struct.
	SubmitFormPostgres struct {
		client mrstorage.DBConnManager
	}
)

// NewSubmitFormPostgres - создаёт объект SubmitFormPostgres.
func NewSubmitFormPostgres(client mrstorage.DBConnManager) *SubmitFormPostgres {
	return &SubmitFormPostgres{
		client: client,
	}
}

// Fetch - comment method.
func (re *SubmitFormPostgres) Fetch(ctx context.Context, _ entity.SubmitFormParams) ([]entity.SubmitForm, error) {
	sql := `
        SELECT
            version,
			rewrite_name,
			form_caption
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
        WHERE
            activity_status = $1 AND form_detailing = $2
        ORDER BY
            form_caption ASC, created_at ASC;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		enum.ActivityStatusPublished,
		enum.ElementDetailingNormal,
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

// FetchByRewriteName - comment method.
func (re *SubmitFormPostgres) FetchByRewriteName(ctx context.Context, rewriteName string) (entity.SubmitForm, error) {
	sql := `
        SELECT
            version,
			rewrite_name,
			form_caption
        FROM
            ` + module.DBSchema + `.` + module.DBTableNameSubmitFormVersions + `
        WHERE
            rewrite_name = $1 AND activity_status = $2 AND form_detailing = $3
        LIMIT 1;`

	var row entity.SubmitForm

	err := re.client.Conn(ctx).QueryRow(
		ctx,
		sql,
		rewriteName,
		enum.ActivityStatusPublished,
		enum.ElementDetailingNormal,
	).Scan(
		&row.Version,
		&row.RewriteName,
		&row.Caption,
	)

	return row, err
}
