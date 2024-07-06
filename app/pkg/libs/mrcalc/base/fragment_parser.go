package base

// ParseFragment - возвращает результат парсинга строки вида '{byWidth}x{byHeight}' в Fragment.
func ParseFragment(str string) (Fragment, error) {
	byWidth, byHeight, err := ParseDoubleSize(str)
	if err != nil {
		return Fragment{}, err
	}

	return Fragment{
		ByWidth:  uint64(byWidth),
		ByHeight: uint64(byHeight),
	}, nil
}
