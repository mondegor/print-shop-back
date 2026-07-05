package measure

import (
	"github.com/mondegor/go-core/mrtype"
)

const (
	// DeltaThousand - дельта (одна десятитысячная) для сравнения float чисел.
	DeltaThousand = 0.001

	// Thousand - тысяча единиц.
	Thousand = 1000.0

	// Million - миллион единиц.
	Million = 1000000.0

	// OneThousandth - одна тысячная единицы.
	OneThousandth = 0.001

	// OneMillionth - одна миллионная единицы.
	OneMillionth = 0.000001
)

type (
	// Meter - метр (m, СИ).
	Meter float64

	// Meter2 - метр квадратный (m2, СИ).
	Meter2 float64

	// Meter3 - метр кубический (m3, СИ).
	Meter3 float64

	// Centimeter - сантиметр (cm).
	Centimeter float64

	// Millimeter - миллиметр (mm).
	Millimeter float64

	// Micrometer - микрометр (µm).
	Micrometer float64

	// Kilogram - килограмм (kg, СИ).
	Kilogram float64

	// Gram - грамм (g).
	Gram float64

	// Milligram - миллиграмм (mg).
	Milligram float64

	// GramPerMeter2 - грамм на метр квадратный (g/m2).
	GramPerMeter2 float64

	// KilogramPerMeter2 - килограмм на метр квадратный (kg/m2, СИ).
	KilogramPerMeter2 float64

	// RangeMeter - интервал (m, СИ).
	RangeMeter mrtype.RangeFloat64

	// RangeKilogram - интервал (kg, СИ).
	RangeKilogram mrtype.RangeFloat64

	// RangeKilogramPerMeter2 - интервал (kg/m2, СИ).
	RangeKilogramPerMeter2 mrtype.RangeFloat64
)
