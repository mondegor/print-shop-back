package mrapp

import "context"

type Transaction interface {
    Begin(ctx context.Context) (Transaction, error)
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
}
