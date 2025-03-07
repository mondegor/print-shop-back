package model

type (
	// BoxStack - стопка коробок.
	BoxStack struct {
		Box
		Quantity uint64
	}
)

// TotalVolume - возвращает общий объём коробок.
func (m BoxStack) TotalVolume() float64 {
	return m.Format.Volume() * float64(m.Quantity)
}
