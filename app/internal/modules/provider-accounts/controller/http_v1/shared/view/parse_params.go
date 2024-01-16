package view_shared

import (
	module "print-shop-back/internal/modules/provider-accounts"
	"strconv"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

func ParseKeyInt32FromPath(c mrcore.ClientContext, key string) mrtype.KeyInt32 {
	value, err := strconv.ParseInt(c.ParamFromPath(key), 10, 32)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
		return 0
	}

	return mrtype.KeyInt32(value)
}

func ParseFilterString(c mrcore.ClientContext, key string) string {
	value, err := mrreq.ParseStr(c.Request(), key, false)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return value
}

func ParseFilterRangeInt64(c mrcore.ClientContext, key string) mrtype.RangeInt64 {
	str, err := mrreq.ParseRangeInt64(c.Request(), key)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return str
}

func ParseSortParams(c mrcore.ClientContext, sorter mrview.ListSorter) mrtype.SortParams {
	value, err := mrreq.ParseSortParams(
		c.Request(),
		module.ParamNameSortField,
		module.ParamNameSortDirection,
	)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	if value.FieldName == "" {
		return sorter.DefaultSort()
	}

	if !sorter.CheckField(value.FieldName) {
		mrctx.Logger(c.Context()).Warning("sort field '%s' is not registered", value.FieldName)
		return sorter.DefaultSort()
	}

	return value
}

func ParsePageParams(c mrcore.ClientContext) mrtype.PageParams {
	value, err := mrreq.ParsePageParams(
		c.Request(),
		module.ParamNamePageIndex,
		module.ParamNamePageSize,
	)

	if err != nil || value.Size < 1 || value.Size > module.PageSizeMax {
		if err != nil {
			mrctx.Logger(c.Context()).Warn(err)
		}

		return mrtype.PageParams{
			Size: module.PageSizeDefault,
		}
	}

	return value
}
