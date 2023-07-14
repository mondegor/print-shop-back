package mrpostgres

import (
    "context"
    "print-shop-back/pkg/mrcontext"
    "strings"
)

func (c *Connection) debugQuery(ctx context.Context, query string) {
    mrcontext.GetLogger(ctx).Debug("SQL Query: %s", strings.Join(strings.Fields(query), " "))
}
