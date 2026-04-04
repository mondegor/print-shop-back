package locationkind

import (
	"encoding/json"
	"fmt"
)

// Возможные типы хранения элементов.
//
// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 - store
// 00010000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 - container group
// 00100000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 - container.
const (
	Store     Enum = iota // склад
	Group                 // группа    - 1 160 000 000 000 000 000
	Container             // контейнер - 2 310 000 000 000 000 000
	Undefined             // неопределено
)

const (
	enumLast = uint8(Container)
	enumName = "LocationKind"
)

type (
	// Enum - перечисление элементов.
	Enum uint8
)

var enumKeys = map[Enum]string{ //nolint:gochecknoglobals
	Store:     "STORE",
	Group:     "GROUP",
	Container: "CONTAINER",
	Undefined: "UNDEFINED",
}

// Is - сообщает, является ли ID указанного типа хранения.
func Is(id uint64, kind Enum) bool {
	return ByID(id) == kind
}

// ByID - возвращает тип хранения в зависимости от ID.
func ByID(id uint64) Enum {
	kind := uint8(id >> 60) //nolint:gosec

	if kind <= enumLast {
		return Enum(kind)
	}

	return Undefined
}

// String - возвращает значение в виде строки.
func (e Enum) String() string {
	if v, ok := enumKeys[e]; ok {
		return v
	}

	return "UNKNOWN"
}

// MarshalJSON - переводит enum значение в строковое представление.
func (e Enum) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(e.String())
	if err != nil {
		return nil, fmt.Errorf("marshal error (source='%s'): %w", enumName, err)
	}

	return bytes, nil
}
