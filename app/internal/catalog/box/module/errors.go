package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrBoxNotFound - box with ID not found.
	ErrBoxNotFound = mrerr.NewKindUser("BoxNotFound", "box with ID={Id} not found")

	// ErrBoxArticleAlreadyExists - box article already exists.
	ErrBoxArticleAlreadyExists = mrerr.NewKindUser("BoxArticleAlreadyExists", "box article '{Name}' already exists")
)
