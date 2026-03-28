package xtype

type (
	// Volume - объём контейнера, места на складе.
	Volume struct {
		Length float64 `json:"length"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	}
)

// Calc - расчёт объёма контейнера в кубических единицах.
func (v Volume) Calc() float64 {
	return v.Width * v.Height * v.Length
}
