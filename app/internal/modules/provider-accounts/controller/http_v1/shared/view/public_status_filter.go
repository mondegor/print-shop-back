package view_shared

import (
	"net/http"
	entity_shared "print-shop-back/internal/modules/provider-accounts/entity/shared"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
	"github.com/mondegor/go-webcore/mrreq"
)

func ParseFilterPublicStatusList(c mrcore.ClientContext, key string) []entity_shared.PublicStatus {
	items, err := parseFilterPublicStatusList(
		c.Request(),
		key,
		[]entity_shared.PublicStatus{
			entity_shared.PublicStatusPublished,
			entity_shared.PublicStatusPublishedShared,
		},
	)

	if err != nil {
		mrctx.Logger(c.Context()).Warn(err)
	}

	return items
}

func parseFilterPublicStatusList(r *http.Request, key string, defaultItems []entity_shared.PublicStatus) ([]entity_shared.PublicStatus, error) {
	def := func(defaultItems []entity_shared.PublicStatus) []entity_shared.PublicStatus {
		if len(defaultItems) == 0 {
			return []entity_shared.PublicStatus{}
		}

		return defaultItems
	}

	enums, err := mrreq.ParseEnumList(r, key)

	if err != nil {
		return def(defaultItems), err
	}

	items, err := entity_shared.ParsePublicStatusList(enums)

	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return def(defaultItems), nil
	}

	return items, nil
}
