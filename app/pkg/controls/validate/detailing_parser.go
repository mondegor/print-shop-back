package validate

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrreq"

	"github.com/mondegor/print-shop-back/pkg/controls/enum"
)

type (
	// RequestDetailingParser - comment interface.
	RequestDetailingParser interface { // TODO: ПЕРЕНЕСТИ
		FilterElementDetailingList(r *http.Request, key string) []enum.ElementDetailing
	}

	// DetailingParser - comment struct.
	DetailingParser struct {
		logger       mrlog.Logger
		defaultItems []enum.ElementDetailing
	}
)

// NewDetailingParser - создаёт объект DetailingParser.
func NewDetailingParser(logger mrlog.Logger) *DetailingParser {
	return &DetailingParser{
		logger: logger,
	}
}

// NewDetailingParserWithDefault - создаёт объект DetailingParser.
func NewDetailingParserWithDefault(logger mrlog.Logger, items []enum.ElementDetailing) *DetailingParser {
	return &DetailingParser{
		logger:       logger,
		defaultItems: items,
	}
}

// FilterElementDetailingList - comment method.
func (p *DetailingParser) FilterElementDetailingList(r *http.Request, key string) []enum.ElementDetailing {
	items, err := p.parseList(r, key)
	if err != nil {
		p.logger.Warn(r.Context(), "DetailingParser", "error", err)

		return p.defaultItems
	}

	if len(items) == 0 {
		return p.defaultItems
	}

	return items
}

func (p *DetailingParser) parseList(r *http.Request, key string) ([]enum.ElementDetailing, error) {
	enumList, err := mrreq.ParseEnumList(r.URL.Query(), key)
	if err != nil {
		return nil, err
	}

	items, err := enum.ParseElementDetailingList(enumList)
	if err != nil {
		return nil, err
	}

	return items, nil
}
