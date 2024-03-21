package usecase_shared

import . "github.com/mondegor/go-sysmess/mrerr"

var (
	FactoryErrCompanyPageRewriteNameAlreadyExists = NewFactory(
		"errProviderAccountsCompanyPageRewriteNameAlreadyExists", ErrorKindUser, "rewrite name '{{ .name }}' already exists")
)
