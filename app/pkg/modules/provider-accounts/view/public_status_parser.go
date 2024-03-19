package view_shared

import (
	"net/http"
	"print-shop-back/pkg/modules/provider-accounts/enums"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

type (
	RequestPublicStatusParser interface {
		FilterPublicStatusList(r *http.Request, key string) []enums.PublicStatus
	}

	PublicStatusParser struct {
		defaultItems []enums.PublicStatus
	}
)

func NewPublicStatusParser() *PublicStatusParser {
	return &PublicStatusParser{}
}

func NewPublicStatusWithDefault(items []enums.PublicStatus) *PublicStatusParser {
	return &PublicStatusParser{
		defaultItems: items,
	}
}

func (p *PublicStatusParser) FilterPublicStatusList(r *http.Request, key string) []enums.PublicStatus {
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

func (p *PublicStatusParser) parseList(r *http.Request, key string) ([]enums.PublicStatus, error) {
	enumList, err := mrreq.ParseEnumList(r, key)

	if err != nil {
		return []enums.PublicStatus{}, err
	}

	items, err := enums.ParsePublicStatusList(enumList)

	if err != nil {
		return []enums.PublicStatus{}, err
	}

	return items, nil
}
