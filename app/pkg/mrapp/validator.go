package mrapp

import "context"

type ValidatorTagFunc func() any

type Validator interface {
    Register(tag string, fn ValidatorTagFunc) error
    Validate(ctx context.Context, structure any) error
}

type UserErrorList interface {
    Add(id string, err error)
    Error() string
}
