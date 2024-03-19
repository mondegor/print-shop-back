package view

import (
	"net/http"
	"print-shop-back/pkg/modules/controls/enums"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	RequestDetailingParser interface {
		FilterElementDetailingList(r *http.Request, key string) []enums.ElementDetailing
	}

	DetailingParser struct {
		defaultItems []enums.ElementDetailing
	}
)

func NewDetailingParser() *DetailingParser {
	return &DetailingParser{}
}

func NewDetailingParserWithDefault(items []enums.ElementDetailing) *DetailingParser {
	return &DetailingParser{
		defaultItems: items,
	}
}

func (p *DetailingParser) FilterElementDetailingList(r *http.Request, key string) []enums.ElementDetailing {
	items, err := p.parseList(r, key)

	if err != nil {
		mrlog.Ctx(r.Context()).Warn().Err(err).Send()
		return p.defaultItems
	}

	if len(items) == 0 {
		return p.defaultItems
	}

	return items
}

func (p *DetailingParser) parseList(r *http.Request, key string) ([]enums.ElementDetailing, error) {
	enumList, err := mrreq.ParseEnumList(r, key)

	if err != nil {
		return []enums.ElementDetailing{}, err
	}

	items, err := enums.ParseElementDetailingList(enumList)

	if err != nil {
		return []enums.ElementDetailing{}, err
	}

	return items, nil
}
