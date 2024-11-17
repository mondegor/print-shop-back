package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

// ErrCompanyPageRewriteNameAlreadyExists - rewrite name already exists.
var ErrCompanyPageRewriteNameAlreadyExists = mrerr.NewProto(
	"provideraccounts.errCompanyPageRewriteNameAlreadyExists", mrerr.ErrorKindUser, "rewrite name '{{ .name }}' already exists")
