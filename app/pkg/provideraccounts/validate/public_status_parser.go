package validate

import (
	"net/http"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrreq"

	"github.com/mondegor/print-shop-back/pkg/provideraccounts/enum"
)

type (
	// RequestPublicStatusParser - comment interface.
	RequestPublicStatusParser interface { // TODO: ПЕРЕНЕСТИ
		FilterPublicStatusList(r *http.Request, key string) []enum.PublicStatus
	}

	// PublicStatusParser - comment struct.
	PublicStatusParser struct {
		logger       mrlog.Logger
		defaultItems []enum.PublicStatus
	}
)

// NewPublicStatusParser - создаёт объект PublicStatusParser.
func NewPublicStatusParser(logger mrlog.Logger) *PublicStatusParser {
	return &PublicStatusParser{
		logger: logger,
	}
}

// NewPublicStatusParserWithDefault - создаёт объект PublicStatusParser.
func NewPublicStatusParserWithDefault(logger mrlog.Logger, items []enum.PublicStatus) *PublicStatusParser {
	return &PublicStatusParser{
		logger:       logger,
		defaultItems: items,
	}
}

// FilterPublicStatusList - comment method.
func (p *PublicStatusParser) FilterPublicStatusList(r *http.Request, key string) []enum.PublicStatus {
	items, err := p.parseList(r, key)
	if err != nil {
		p.logger.Warn(r.Context(), "PublicStatusParser", "error", err)

		return p.defaultItems
	}

	if len(items) == 0 {
		return p.defaultItems
	}

	return items
}

func (p *PublicStatusParser) parseList(r *http.Request, key string) ([]enum.PublicStatus, error) {
	enumList, err := mrreq.ParseEnumList(r.URL.Query(), key)
	if err != nil {
		return nil, err
	}

	items, err := enum.ParsePublicStatusList(enumList)
	if err != nil {
		return nil, err
	}

	return items, nil
}
