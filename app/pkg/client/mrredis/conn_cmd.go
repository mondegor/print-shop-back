package mrredis

import (
    "context"
    "strings"
    "time"
)

func (c *Connection) GetStruct(ctx context.Context, key string, data any) error {
    err := c.conn.Get(ctx, key).Scan(data)

    if err != nil {
        return c.wrapError(err)
    }

    c.debugCmd("get-struct", key, data)

    return nil
}

func (c *Connection) SetStruct(ctx context.Context, key string, data any, expiration time.Duration) error {
    c.debugCmd("set-struct", key, data)

    err := c.conn.Set(ctx, key, data, expiration).Err()

    if err != nil {
        return c.wrapError(err)
    }

    return nil
}

func (c *Connection) Delete(ctx context.Context, key ...string) error {
    c.debugCmd("delete-row", strings.Join(key, ", "), nil)

    err := c.conn.Del(ctx, key...).Err()

    if err != nil {
        return c.wrapError(err)
    }

    return nil
}
