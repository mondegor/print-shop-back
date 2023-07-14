package mrpostgres

import (
    "calc-user-data-back-adm/pkg/mrcontext"
    "context"
    "strings"
)

func (c *Connection) debugQuery(ctx context.Context, query string) {
    mrcontext.GetLogger(ctx).Debug("SQL Query: %s", strings.Join(strings.Fields(query), " "))
}
