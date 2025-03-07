package model

type (
	// PackInBox - результат работы алгоритма Algo.
	PackInBox struct {
		FullBox FilledBoxStack // полностью заполненная стопка коробок
		RestBox FilledBox      // коробка с остатком изделий
	}
)

// BoxesQuantity - возвращает общее количество коробок.
func (m *PackInBox) BoxesQuantity() uint64 {
	if m.RestBox.IsEmpty() {
		return m.FullBox.Quantity
	}

	return m.FullBox.Quantity + 1
}

// BoxesVolume - возвращает общий внешний объём коробок.
func (m *PackInBox) BoxesVolume() float64 {
	if m.RestBox.IsEmpty() {
		return m.FullBox.TotalVolume()
	}

	return m.FullBox.TotalVolume() + m.RestBox.Box.Volume()
}

// BoxesInnerVolume - возвращает общий внутренний объём коробок.
func (m *PackInBox) BoxesInnerVolume() float64 {
	if m.RestBox.IsEmpty() {
		return m.FullBox.TotalInnerVolume()
	}

	return m.FullBox.TotalInnerVolume() + m.RestBox.Box.InnerVolume()
}

// BoxesWeight - возвращает общий вес коробок с изделиями.
func (m *PackInBox) BoxesWeight() float64 {
	if m.RestBox.IsEmpty() {
		return m.FullBox.TotalWeight()
	}

	return m.FullBox.TotalWeight() + m.RestBox.Weight()
}

// ProductsVolume - возвращает общий объём изделий.
func (m *PackInBox) ProductsVolume() float64 {
	if m.RestBox.IsEmpty() {
		return m.FullBox.Product.TotalVolume()
	}

	return m.FullBox.Product.TotalVolume() + m.RestBox.Product.TotalVolume()
}
