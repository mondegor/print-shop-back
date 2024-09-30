package validate

import (
	"net/http"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrreq"

	"github.com/mondegor/print-shop-back/pkg/controls/enum"
)

type (
	// RequestDetailingParser - comment interface.
	RequestDetailingParser interface {
		FilterElementDetailingList(r *http.Request, key string) []enum.ElementDetailing
	}

	// DetailingParser - comment struct.
	DetailingParser struct {
		defaultItems []enum.ElementDetailing
	}
)

// NewDetailingParser - создаёт объект DetailingParser.
func NewDetailingParser() *DetailingParser {
	return &DetailingParser{}
}

// NewDetailingParserWithDefault - создаёт объект DetailingParser.
func NewDetailingParserWithDefault(items []enum.ElementDetailing) *DetailingParser {
	return &DetailingParser{
		defaultItems: items,
	}
}

// FilterElementDetailingList - comment method.
func (p *DetailingParser) FilterElementDetailingList(r *http.Request, key string) []enum.ElementDetailing {
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

func (p *DetailingParser) parseList(r *http.Request, key string) ([]enum.ElementDetailing, error) {
	enumList, err := mrreq.ParseEnumList(r, key)
	if err != nil {
		return nil, err
	}

	items, err := enum.ParseElementDetailingList(enumList)
	if err != nil {
		return nil, err
	}

	return items, nil
}
