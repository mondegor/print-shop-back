package module

import (
	"github.com/mondegor/go-core/errors"
)

// ErrCompanyPageRewriteNameAlreadyExists - rewrite name already exists.
var ErrCompanyPageRewriteNameAlreadyExists = errors.NewUserProto("CompanyPageRewriteNameAlreadyExists", "rewrite name '{Name}' already exists")
