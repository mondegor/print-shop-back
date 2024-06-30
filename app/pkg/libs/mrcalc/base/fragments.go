package base

type (
	// Fragments - comment struct.
	Fragments []Fragment
)

// Total - возвращает общее кол-во единиц.
func (f *Fragments) Total() (total uint64) {
	for _, fragment := range *f {
		total += fragment.Total()
	}

	return total
}
