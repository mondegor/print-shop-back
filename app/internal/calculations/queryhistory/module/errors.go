package module

import (
	"github.com/mondegor/go-core/errors"
)

// ErrQueryHistoryNotFound - query with ShortLink not found.
var ErrQueryHistoryNotFound = errors.NewUserProto("QueryHistoryNotFound", "query {ShortLink} not found")
