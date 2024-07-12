package base

// ParseFragment - возвращает результат парсинга строки вида '{byWidth}x{byHeight}' в Fragment.
func ParseFragment(str string) (Fragment, error) {
	by, err := ParseDoubleSize(str)
	if err != nil {
		return Fragment{}, err
	}

	return Fragment{
		ByWidth:  uint64(by[0]),
		ByHeight: uint64(by[1]),
	}, nil
}
