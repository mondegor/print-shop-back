package mrredis

func (c *Connection) debugCmd(command string, key string, data any) {
    c.logger.Debug("Redis: cmd=%s, key=%s, struct=%+v", command, key, data)
}
