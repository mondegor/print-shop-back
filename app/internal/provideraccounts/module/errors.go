package module

import (
	"github.com/mondegor/go-sysmess/errors"
)

// ErrCompanyPageRewriteNameAlreadyExists - rewrite name already exists.
var ErrCompanyPageRewriteNameAlreadyExists = errors.NewUserProto("CompanyPageRewriteNameAlreadyExists", "rewrite name '{Name}' already exists")
