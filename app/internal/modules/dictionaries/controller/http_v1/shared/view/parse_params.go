package view_shared

import (
	module "print-shop-back/internal/modules/dictionaries"
	"strconv"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrreq"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/mondegor/go-webcore/mrview"
)

func RawQueryParam(c mrcore.ClientContext, key string) string {
	return c.Request().URL.Query().Get(key)
}

func ParseKeyInt32(c mrcore.ClientContext, key string) mrtype.KeyInt32 {
	value, err := mrreq.ParseInt64(c.Request(), key, true)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
		return 0
	}

	return mrtype.KeyInt32(value)
}

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

func ParseFilterStatusList(c mrcore.ClientContext, key string) []mrenum.ItemStatus {
	items, err := mrreq.ParseItemStatusList(
		c.Request(),
		key,
		mrenum.ItemStatusEnabled,
	)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return items
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
