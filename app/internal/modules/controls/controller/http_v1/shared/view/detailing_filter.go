package view_shared

import (
	"net/http"
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrreq"
)

func ParseFilterElementDetailingList(c mrcore.ClientContext, key string) []entity_shared.ElementDetailing {
	items, err := parseFilterDetailingList(
		c.Request(),
		key,
		entity_shared.ElementDetailingNormal,
	)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return items
}

func parseFilterDetailingList(r *http.Request, key string, defaultItem entity_shared.ElementDetailing) ([]entity_shared.ElementDetailing, error) {
	def := func(defaultItem entity_shared.ElementDetailing) []entity_shared.ElementDetailing {
		if defaultItem == 0 {
			return []entity_shared.ElementDetailing{}
		}

		return []entity_shared.ElementDetailing{defaultItem}
	}

	enums, err := mrreq.ParseEnumList(r, key)

	if err != nil {
		return def(defaultItem), err
	}

	items, err := entity_shared.ParseElementDetailingList(enums)

	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return def(defaultItem), nil
	}

	return items, nil
}
