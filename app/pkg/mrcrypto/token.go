package mrcrypto

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "strings"
)

func GenTokenBase64(length int) string {
    return base64.StdEncoding.EncodeToString(GenToken(length))
}

func GenTokenHex(length int) string {
    return hex.EncodeToString(GenToken(length))
}

func GenTokenHexWithDelimiter(length int, repeat int) string {
    if repeat < 1 {
        panic("param repeat < 1")
    }

    var s []string

    for i := 0; i < repeat; i++ {
        s = append(s, hex.EncodeToString(GenToken(length)))
    }

    return strings.Join(s, "-")
}

func GenToken(length int) []byte {
    if length < 1 {
        panic("param length < 1")
    }

    value := make([]byte, length)

    _, err := rand.Read(value)

    if err != nil {
        return []byte{}
    }

    return value
}
