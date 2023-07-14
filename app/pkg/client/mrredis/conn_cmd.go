package mrredis

import "time"

func (c *Connection) SetStruct(key string, data interface{}, expiration time.Duration) error {
    c.debugCmd("set-struct", key, data)

    if err := c.conn.Set(key, data, expiration).Err(); err != nil {
        return c.wrapError(err)
    }

    return nil
}

func (c *Connection) GetStruct(key string, data interface{}) error {
    c.debugCmd("get-struct", key, data)

    err := c.conn.Get( "name").Scan(&data)

    if err != nil {
        return c.wrapError(err)
    }

    return nil
}
