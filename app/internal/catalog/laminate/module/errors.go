package module

import (
	"github.com/mondegor/go-sysmess/errors"
)

var (
	// ErrLaminateNotFound - laminate with ID not found.
	ErrLaminateNotFound = errors.NewUserProto("LaminateNotFound", "laminate with ID={Id} not found")

	// ErrLaminateArticleAlreadyExists - laminate article already exists.
	ErrLaminateArticleAlreadyExists = errors.NewUserProto("LaminateArticleAlreadyExists", "laminate article '{Name}' already exists")
)
