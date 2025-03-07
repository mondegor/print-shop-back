package model

type (
	// FilledBox - коробка заполненная изделиями.
	FilledBox struct {
		Box     Box          // исходная коробка
		Product ProductStack // изделия размещённые в коробке
	}
)

// IsEmpty - проверяется заполнена ли коробка.
func (m FilledBox) IsEmpty() bool {
	return m.Product.Quantity == 0
}

// Weight - возвращает общий вес коробки с изделиями.
func (m FilledBox) Weight() float64 {
	return m.Box.Weight + m.Product.TotalWeight()
}

// UnusedVolumePercent - возвращает незаполненный объём коробки в %.
func (m FilledBox) UnusedVolumePercent() float64 {
	innerFormat := m.Box.NarrowedFormat()

	if !innerFormat.IsValid() {
		return 0
	}

	return 100 - 100*m.Product.TotalVolume()/innerFormat.Volume()
}
