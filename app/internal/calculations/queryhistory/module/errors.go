package module

import "github.com/mondegor/go-sysmess/mrerr"

// ErrQueryHistoryNotFound - query with ShortLink not found.
var ErrQueryHistoryNotFound = mrerr.NewKindUser("QueryHistoryNotFound", "query {ShortLink} not found")
