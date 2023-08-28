package mrlib

import (
    "context"
    "fmt"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrcontext"
    "print-shop-back/pkg/mrerr"
    "reflect"
    "strings"

    "github.com/go-playground/validator/v10"
)

// go get -u github.com/go-playground/validator/v10

type Validator struct {
    validate *validator.Validate
}

// Make sure the Validator conforms with the mrapp.Validator interface
var _ mrapp.Validator = (*Validator)(nil)

func NewValidator() *Validator {
    validate := validator.New()

    validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }

        return name
    })

    return &Validator {
        validate: validate,
    }
}

func (v *Validator) Register(tag string, fn mrapp.ValidatorTagFunc) error {
    if vfn, ok := fn().(func (fl validator.FieldLevel) bool); ok {
        err := v.validate.RegisterValidation(tag, vfn)

        if err != nil {
            return err
        }

        return nil
    }

    return mrerr.ErrInternalTypeAssertion.New("func (fl validator.FieldLevel) bool", fn())
}

func (v *Validator) Validate(ctx context.Context, structure any) error {
    err := v.validate.Struct(structure)

    if err == nil {
        return nil
    }

    if _, ok := err.(*validator.InvalidValidationError); ok {
        return mrerr.ErrInternal.Wrap(err)
    }

    errorList := NewUserErrorList()

    for _, errField := range err.(validator.ValidationErrors) {
        errorList.Add(errField.Field(), v.createAppError(errField))

        v.logger(ctx).Debug(
            "Namespace: %s\n" +
            "Field: %s\n" +
            "StructNamespace: %s\n" +
            "StructField: %s\n" +
            "Tag: %s\n" +
            "ActualTag: %s\n" +
            "Kind: %v\n" +
            "Type: %v\n" +
            "Value: %v\n" +
            "Param: %s",
            errField.Namespace(),
            errField.Field(),
            errField.StructNamespace(),
            errField.StructField(),
            errField.Tag(),
            errField.ActualTag(),
            errField.Kind(),
            errField.Type(),
            errField.Value(),
            errField.Param(),
        )
    }

    return errorList
}

func (v *Validator) createAppError(field validator.FieldError) *mrerr.AppError {
    id := []byte("errValidation")
    tag := []byte(field.Tag())

    if len(tag) == 0 {
        return mrerr.New(mrerr.ErrorCode(id), string(id))
    }

    tag[0] -= 32 // to uppercase first char
    id = append(id, tag...)

    return mrerr.New(
        mrerr.ErrorCode(id),
        fmt.Sprintf("%s: value='{{ .value }}'", id),
        field.Value(),
    )
}

func (v *Validator) logger(ctx context.Context) mrapp.Logger {
    return mrcontext.GetLogger(ctx)
}
