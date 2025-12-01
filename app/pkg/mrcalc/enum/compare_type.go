package enum

// Результаты сравнения прямоугольных форматов.
const (
	CompareTypeEqual         CompareType = iota // форматы равны
	CompareTypeFirstInside                      // первый формат входит во второй
	CompareTypeSecondInside                     // второй формат входит в первый
	CompareTypeNotCompatible                    // форматы не совместимы
)

type (
	// CompareType - Тип сравнения используемый для сравнения двух прямоугольных форматов.
	CompareType uint8
)
