package mrlang

import (
    "fmt"
    "strconv"
)

type NamedArg struct {
    name  string
    value any
}

func NewArg(name string, value any) NamedArg {
    return NamedArg{
        name:  name,
        value: value,
    }
}

func (n *NamedArg) valueString() string {
    switch val := n.value.(type) {
        case string:
            return val
        case int:
            return strconv.FormatInt(int64(val), 10)
        case int32:
            return strconv.FormatInt(int64(val), 10)
        case int64:
            return strconv.FormatInt(val, 10)
        default:
            return fmt.Sprintf("%v", val)
    }
}
