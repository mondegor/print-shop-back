package cutting

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

// AlgoQuantity - возвращает общее количество резов необходимых совершить резательным оборудованием
// на печатном формате, чтобы получить заданное кол-во элементов нужного формата.
// fragments - это список блоков прямоугольных элементов разложенных на печатном формате.
// distance - это расстояние между элементами блока, для того чтобы понимать нужны ли двойные резы.
// На практике массив fragments содержит не более двух блоков (рассчитывается при спуске полос).
func AlgoQuantity(fragments []base.Fragment, distance rect.Format) (quantity uint64) {
	if len(fragments) == 0 {
		return 0
	}

	for i := range fragments {
		if fragments[i].Total() == 0 {
			continue
		}

		quantity += getQuantityLine(distance.Width, fragments[i].ByWidth) +
			getQuantityLine(distance.Height, fragments[i].ByHeight)
	}

	return quantity
}

func getQuantityLine(distance int64, countItems uint64) uint64 {
	if distance > 0 {
		return 2 * countItems
	}

	return countItems + 1
}
