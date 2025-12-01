package model

type (
	// FilledBoxStack - стопка коробок заполненных изделиями.
	FilledBoxStack struct {
		FilledBox
		Quantity uint64 // количество коробок
	}
)

// TotalVolume - возвращает общий объём коробок с изделиями.
func (m FilledBoxStack) TotalVolume() float64 {
	return m.Box.Volume() * float64(m.Quantity)
}

// TotalInnerVolume - возвращает общий объём коробок с изделиями.
func (m FilledBoxStack) TotalInnerVolume() float64 {
	return m.Box.InnerVolume() * float64(m.Quantity)
}

// TotalWeight - возвращает общий вес коробок с изделиями.
func (m FilledBoxStack) TotalWeight() float64 {
	return m.Weight() * float64(m.Quantity)
}
