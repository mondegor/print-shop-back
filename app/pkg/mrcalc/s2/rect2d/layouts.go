package rect2d

type (
	// Layouts - список схем размещения элементов.
	Layouts []Layout
)

// TotalQuantity - возвращает общее кол-во элементов во всех схемах.
func (l Layouts) TotalQuantity() (total uint64) {
	for _, layout := range l {
		total += layout.Quantity()
	}

	return total
}
