package module

import "github.com/mondegor/go-sysmess/mrerr"

var (
	// ErrPaperNotFound - paper with ID not found.
	ErrPaperNotFound = mrerr.NewKindUser("PaperNotFound", "paper with ID={Id} not found")

	// ErrPaperArticleAlreadyExists - paper article already exists.
	ErrPaperArticleAlreadyExists = mrerr.NewKindUser("PaperArticleAlreadyExists", "paper article '{Name}' already exists")
)
