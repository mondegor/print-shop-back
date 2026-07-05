package rect2d

type (
	// Layout - табличная схема размещения элементов.
	Layout struct {
		ByWidth  uint64 `json:"byWidth"`  // кол-во элементов по ширине
		ByHeight uint64 `json:"byHeight"` // кол-во элементов по высоте
	}
)

// Min - возвращает минимальное количество элементов размещённых по ширине или высоте.
func (l Layout) Min() uint64 {
	if l.ByWidth > l.ByHeight {
		return l.ByHeight
	}

	return l.ByWidth
}

// Max - возвращает максимальное количество элементов размещённых по ширине или высоте.
func (l Layout) Max() uint64 {
	if l.ByWidth > l.ByHeight {
		return l.ByWidth
	}

	return l.ByHeight
}

// Quantity - возвращает общее кол-во элементов.
func (l Layout) Quantity() uint64 {
	return l.ByWidth * l.ByHeight
}
