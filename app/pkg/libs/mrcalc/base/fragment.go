package base

type (
	// Fragment - фрагмент прямоугольного формата с количеством размещенных на нём элементов.
	Fragment struct {
		ByWidth  uint64 `json:"byWidth"`  // кол-во элементов по ширине
		ByHeight uint64 `json:"byHeight"` // кол-во элементов по высоте
	}
)

// Max - возвращает максимальное количество элементов по ширине или высоте.
func (f *Fragment) Max() uint64 {
	if f.ByWidth > f.ByHeight {
		return f.ByWidth
	}

	return f.ByHeight
}

// Total - возвращает общее кол-во элементов.
func (f *Fragment) Total() uint64 {
	return f.ByWidth * f.ByHeight
}
