package rect

type (
	// Item - прямоугольный элемент с технологическими границами.
	Item struct {
		Format
		Distance Format `json:"margins"`
	}
)

// WithDistance - возвращается прямоугольный элемент с учётом технологических границ.
func (f *Item) WithDistance() Format {
	return Format{
		Width:  f.Width + f.Distance.Width,
		Height: f.Height + f.Distance.Height,
	}
}
