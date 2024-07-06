package base

const (
	CompareTypeEqual         CompareType = iota // CompareTypeEqual - форматы равны
	CompareTypeFirstInside                      // CompareTypeFirstInside - первый формат входит во второй
	CompareTypeSecondInside                     // CompareTypeSecondInside - второй формат входит в первый
	CompareTypeNotCompatible                    // CompareTypeNotCompatible - форматы не совместимы
)

type (
	// CompareType - Тип сравнения используемый для сравнения двух прямоугольных форматов.
	CompareType uint8
)
