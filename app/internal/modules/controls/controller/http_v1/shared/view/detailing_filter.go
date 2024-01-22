package view_shared

import (
	"net/http"
	entity_shared "print-shop-back/internal/modules/controls/entity/shared"

	"github.com/mondegor/go-webcore/mrserver/mrreq"
)

func parseFilterDetailingList(r *http.Request, key string, defaultItems []entity_shared.ElementDetailing) ([]entity_shared.ElementDetailing, error) {
	def := func(defaultItems []entity_shared.ElementDetailing) []entity_shared.ElementDetailing {
		if len(defaultItems) == 0 {
			return []entity_shared.ElementDetailing{}
		}

		return defaultItems
	}

	enums, err := mrreq.ParseEnumList(r, key)

	if err != nil {
		return def(defaultItems), err
	}

	items, err := entity_shared.ParseElementDetailingList(enums)

	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return def(defaultItems), nil
	}

	return items, nil
}
