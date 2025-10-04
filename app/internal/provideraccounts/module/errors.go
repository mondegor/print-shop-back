package module

import "github.com/mondegor/go-sysmess/mrerr"

// ErrCompanyPageRewriteNameAlreadyExists - rewrite name already exists.
var ErrCompanyPageRewriteNameAlreadyExists = mrerr.NewKindUser("CompanyPageRewriteNameAlreadyExists", "rewrite name '{Name}' already exists")
