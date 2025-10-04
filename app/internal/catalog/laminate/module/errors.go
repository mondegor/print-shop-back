package module

import "github.com/mondegor/go-sysmess/mrerr"

var (
	// ErrLaminateNotFound - laminate with ID not found.
	ErrLaminateNotFound = mrerr.NewKindUser("LaminateNotFound", "laminate with ID={Id} not found")

	// ErrLaminateArticleAlreadyExists - laminate article already exists.
	ErrLaminateArticleAlreadyExists = mrerr.NewKindUser("LaminateArticleAlreadyExists", "laminate article '{Name}' already exists")
)
