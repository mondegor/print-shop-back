package mrentity

var (
    boolFalse = newFalse()
    boolTrue = newTrue()
)

func newFalse() *bool {
    value := false
    return &value
}

func newTrue() *bool {
    value := true
    return &value
}

func BoolPointer(value bool) *bool {
    if value {
        return boolTrue
    }

    return boolFalse
}
