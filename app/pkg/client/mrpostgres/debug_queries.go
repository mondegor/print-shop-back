package mrpostgres

import (
    "print-shop-back/pkg/mrcontext"
    "context"
    "strings"
)

func (c *Connection) debugQuery(ctx context.Context, query string) {
    mrcontext.GetLogger(ctx).Debug("SQL Query: %s", strings.Join(strings.Fields(query), " "))
}
