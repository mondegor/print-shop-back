package module

import (
	"github.com/mondegor/go-core/errors"
)

var (
	// ErrBoxNotFound - box with ID not found.
	ErrBoxNotFound = errors.NewUserProto("BoxNotFound", "box with ID={Id} not found")

	// ErrBoxArticleAlreadyExists - box article already exists.
	ErrBoxArticleAlreadyExists = errors.NewUserProto("BoxArticleAlreadyExists", "box article '{Name}' already exists")
)
