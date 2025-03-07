package cutting

import "github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"

// AlgoQuantity - возвращает общее количество резов необходимых совершить резательным оборудованием
// на поверхности, чтобы получить заданное кол-во элементов нужного формата.
// layouts - это список блоков прямоугольных элементов разложенных на поверхности.
// distance - это расстояние между элементами блока, для того чтобы понимать нужны ли двойные резы.
// На практике массив layouts содержит не более двух блоков (рассчитывается при спуске полос).
func AlgoQuantity(layouts []rect2d.Layout, distance rect2d.Format) (quantity uint64) {
	if len(layouts) == 0 {
		return 0
	}

	for i := range layouts {
		if layouts[i].Quantity() == 0 {
			continue
		}

		quantity += getQuantityLine(distance.Width, layouts[i].ByWidth) +
			getQuantityLine(distance.Height, layouts[i].ByHeight)
	}

	return quantity
}

func getQuantityLine(distance float64, countItems uint64) uint64 {
	if distance > 0 {
		return 2 * countItems
	}

	return countItems + 1
}
