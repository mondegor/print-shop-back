package module

import (
	"github.com/mondegor/go-core/errors"
)

var (
	// ErrSubmitFormRequired - form ID is required.
	ErrSubmitFormRequired = errors.NewUserProto("SubmitFormRequired", "form ID is required")

	// ErrSubmitFormNotFound - form with ID not found.
	ErrSubmitFormNotFound = errors.NewUserProto("SubmitFormNotFound", "form with ID={Id} not found")

	// ErrSubmitFormRewriteNameAlreadyExists - rewrite name already exists.
	ErrSubmitFormRewriteNameAlreadyExists = errors.NewUserProto("SubmitFormRewriteNameAlreadyExists", "rewrite name '{Name}' already exists")

	// ErrSubmitFormParamNameAlreadyExists - param name already exists.
	ErrSubmitFormParamNameAlreadyExists = errors.NewUserProto("SubmitFormParamNameAlreadyExists", "param name '{Name}' already exists")

	// ErrSubmitFormIsDisabled - form with ID is disabled.
	ErrSubmitFormIsDisabled = errors.NewUserProto("SubmitFormIsDisabled", "form with ID={Id} is disabled")

	// ErrFormElementNotFound - form element with ID not found.
	ErrFormElementNotFound = errors.NewUserProto("FormElementNotFound", "form element with ID={Id} not found")

	// ErrFormElementParamNameAlreadyExists - param name already exists.
	ErrFormElementParamNameAlreadyExists = errors.NewUserProto("FormElementParamNameAlreadyExists", "param name '{Name}' already exists")

	// ErrFormElementDetailingNotAllowed - item detailing not allowed for form detailing.
	ErrFormElementDetailingNotAllowed = errors.NewUserProto(
		"FormElementDetailingNotAllowed", "item detailing '{Name1}' not allowed for form detailing '{Name2}'",
	)
)
