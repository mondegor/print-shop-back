package repository

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrpostgres/db"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/provideraccounts/module"
	"github.com/mondegor/print-shop-back/internal/provideraccounts/section/adm/entity"
)

type (
	// CompanyPagePostgres - comment struct.
	CompanyPagePostgres struct {
		client        mrstorage.DBConnManager
		sqlBuilder    mrstorage.SQLBuilder
		repoTotalRows db.TotalRowsFetcher[uint64]
	}
)

// NewCompanyPagePostgres - создаёт объект CompanyPagePostgres.
func NewCompanyPagePostgres(client mrstorage.DBConnManager, sqlBuilder mrstorage.SQLBuilder) *CompanyPagePostgres {
	return &CompanyPagePostgres{
		client:     client,
		sqlBuilder: sqlBuilder,
		repoTotalRows: db.NewTotalRowsFetcher[uint64](
			client,
			module.DBTableNameCompaniesPages,
		),
	}
}

// FetchWithTotal - comment method.
func (re *CompanyPagePostgres) FetchWithTotal(ctx context.Context, params entity.CompanyPageParams) (rows []entity.CompanyPage, countRows uint64, err error) {
	condition := re.sqlBuilder.Condition().Build(re.fetchCondition(params.Filter))

	total, err := re.repoTotalRows.Fetch(ctx, condition)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	if params.Pager.Size > total {
		params.Pager.Size = total
	}

	orderBy := re.sqlBuilder.OrderBy().Build(re.fetchOrderBy(params.Sorter))
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
	maxRows uint64,
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
            ` + orderBy.String() + limit.String() + `;`

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

func (re *CompanyPagePostgres) fetchCondition(filter entity.CompanyPageListFilter) mrstorage.SQLPartFunc {
	return re.sqlBuilder.Condition().HelpFunc(
		func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc {
			return c.JoinAnd(
				c.FilterLikeFields([]string{"UPPER(rewrite_name)", "UPPER(page_title)", "UPPER(site_url)"}, strings.ToUpper(filter.SearchText)),
				c.FilterAnyOf("page_status", filter.Statuses),
			)
		},
	)
}

func (re *CompanyPagePostgres) fetchOrderBy(sorter mrtype.SortParams) mrstorage.SQLPartFunc {
	return re.sqlBuilder.OrderBy().HelpFunc(
		func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc {
			return o.JoinComma(
				o.Field(sorter.FieldName, sorter.Direction),
				o.Field("account_id", mrenum.SortDirectionASC),
			)
		},
	)
}
