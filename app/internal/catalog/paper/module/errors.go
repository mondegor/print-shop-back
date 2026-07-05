package module

import (
	"github.com/mondegor/go-sysmess/errors"
)

var (
	// ErrPaperNotFound - paper with ID not found.
	ErrPaperNotFound = errors.NewUserProto("PaperNotFound", "paper with ID={Id} not found")

	// ErrPaperArticleAlreadyExists - paper article already exists.
	ErrPaperArticleAlreadyExists = errors.NewUserProto("PaperArticleAlreadyExists", "paper article '{Name}' already exists")
)
