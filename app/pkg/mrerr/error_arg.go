package mrerr

import (
    "fmt"
)

type Arg map[string]any

func (a Arg) String() string {
    var buf []byte
    firstItem := true

    buf = append(buf, '[')

    for key, value := range a {
        if firstItem {
            firstItem = false
        } else {
            buf = append(buf, ',', ' ')
        }

        buf = append(buf, key...)
        buf = append(buf, ':', ' ')
        buf = append(buf, fmt.Sprintf("%v", value)...)
    }

    buf = append(buf, ']')

    return string(buf)
}
