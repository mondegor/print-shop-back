package module

import (
	"github.com/mondegor/go-sysmess/mrerr"
)

var (
	// ErrSubmitFormRequired - form ID is required.
	ErrSubmitFormRequired = mrerr.NewKindUser("SubmitFormRequired", "form ID is required")

	// ErrSubmitFormNotFound - form with ID not found.
	ErrSubmitFormNotFound = mrerr.NewKindUser("SubmitFormNotFound", "form with ID={Id} not found")

	// ErrSubmitFormRewriteNameAlreadyExists - rewrite name already exists.
	ErrSubmitFormRewriteNameAlreadyExists = mrerr.NewKindUser("SubmitFormRewriteNameAlreadyExists", "rewrite name '{Name}' already exists")

	// ErrSubmitFormParamNameAlreadyExists - param name already exists.
	ErrSubmitFormParamNameAlreadyExists = mrerr.NewKindUser("SubmitFormParamNameAlreadyExists", "param name '{Name}' already exists")

	// ErrSubmitFormIsDisabled - form with ID is disabled.
	ErrSubmitFormIsDisabled = mrerr.NewKindUser("SubmitFormIsDisabled", "form with ID={Id} is disabled")

	// ErrFormElementNotFound - form element with ID not found.
	ErrFormElementNotFound = mrerr.NewKindUser("FormElementNotFound", "form element with ID={Id} not found")

	// ErrFormElementParamNameAlreadyExists - param name already exists.
	ErrFormElementParamNameAlreadyExists = mrerr.NewKindUser("FormElementParamNameAlreadyExists", "param name '{Name}' already exists")

	// ErrFormElementDetailingNotAllowed - item detailing not allowed for form detailing.
	ErrFormElementDetailingNotAllowed = mrerr.NewKindUser("FormElementDetailingNotAllowed", "item detailing '{Name1}' not allowed for form detailing '{Name2}'")
)
