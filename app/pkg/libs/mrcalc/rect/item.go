package rect

type (
	// Item - прямоугольный элемент с технологическими границами.
	Item struct {
		Format
		Border Format `json:"border"`
	}
)

// WithBorder - возвращается прямоугольный элемент с учётом технологических границ.
func (f *Item) WithBorder() Format {
	return Format{
		Width:  f.Width + f.Border.Width,
		Height: f.Height + f.Border.Height,
	}
}
