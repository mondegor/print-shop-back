package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrtype/sortdirection"

	"print-shop-back/internal/provideraccounts/module"
	"print-shop-back/internal/provideraccounts/section/adm/entity"
)

type (
	// CompanyPagePostgres - comment struct.
	CompanyPagePostgres struct {
		client        mrstorage.DBConnManager
		sqlBuilder    mrstorage.SQLBuilder
		repoTotalRows db.TotalRowsFetcher[int]
	}
)

// NewCompanyPagePostgres - создаёт объект CompanyPagePostgres.
func NewCompanyPagePostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *CompanyPagePostgres {
	return &CompanyPagePostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoTotalRows: db.NewTotalRowsFetcher[int](
			client,
			module.DBTableNameCompaniesPages,
		),
	}
}

// FetchWithTotal - comment method.
func (re *CompanyPagePostgres) FetchWithTotal(ctx context.Context, params entity.CompanyPageParams) (rows []entity.CompanyPage, countRows int, err error) {
	condition := re.sqlBuilder.Condition().BuildFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.FilterLikeFields([]string{"UPPER(rewrite_name)", "UPPER(page_title)", "UPPER(site_url)"}, strings.ToUpper(params.Filter.SearchText)),
				c.FilterAnyOf("page_status", params.Filter.Statuses),
			)
		},
	)

	total, err := re.repoTotalRows.Fetch(ctx, condition)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	if params.Pager.Size > total {
		params.Pager.Size = total
	}

	orderBy := re.sqlBuilder.OrderBy().BuildFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(params.Sorter.Column, params.Sorter.Direction),
				o.Field("account_id", sortdirection.ASC),
			)
		},
	)
	limit := re.sqlBuilder.Limit().Build(params.Pager.Index, params.Pager.Size)

	rows, err = re.fetch(ctx, condition, orderBy, limit, params.Pager.Size)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

// Fetch - comment method.
func (re *CompanyPagePostgres) fetch(
	ctx context.Context,
	condition mrstorage.SQLPart,
	orderBy mrstorage.SQLPart,
	limit mrstorage.SQLPart,
	maxRows int,
) ([]entity.CompanyPage, error) {
	whereStr, whereArgs := condition.WithPrefix(" WHERE ").ToSQL()

	sql := `
        SELECT
            account_id,
            rewrite_name as rewriteName,
            page_title as pageTitle,
            COALESCE(logo_meta ->> 'path', '') as logoUrl,
            site_url as siteUrl,
            page_status,
			created_at as createdAt,
            updated_at as updatedAt
        FROM
            ` + module.DBTableNameCompaniesPages + `
		` + whereStr + `
        ORDER BY
            ` + mrstorage.ToSQL(orderBy) + mrstorage.ToSQL(limit) + `;`

	cursor, err := re.client.Conn(ctx).Query(
		ctx,
		sql,
		whereArgs...,
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	rows := make([]entity.CompanyPage, 0, maxRows)

	for cursor.Next() {
		var row entity.CompanyPage

		err = cursor.Scan(
			&row.AccountID,
			&row.RewriteName,
			&row.PageTitle,
			&row.LogoURL,
			&row.SiteURL,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, cursor.Err()
}
