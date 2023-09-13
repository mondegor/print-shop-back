package usecase

import . "github.com/mondegor/go-sysmess/mrerr"

var (
    ErrFormFieldItemTemplateNotFound = NewFactory(
        "errFormFieldItemTemplateNotFound", ErrorKindUser, "template with ID={{ .id }} not found")

    ErrFormFieldItemParamNameAlreadyExists = NewFactory(
        "errFormFieldItemParamNameAlreadyExists", ErrorKindUser, "param name '{{ .name }}' is already exists")
)
