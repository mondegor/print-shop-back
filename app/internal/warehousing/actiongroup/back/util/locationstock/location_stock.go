package locationstock

import (
	"iter"

	"github.com/mondegor/print-shop-back/internal/warehousing/actiongroup/back/dto"
)

// ChunkByLocationID - возвращает чанки указанного массива, разбитых по LocationID.
// Указанный массив обязан быть отсортированным.
func ChunkByLocationID(s []dto.LocationStock) iter.Seq[[]dto.LocationStock] {
	return func(yield func([]dto.LocationStock) bool) {
		if len(s) == 0 {
			return
		}

		i := 0
		locationID := s[0].LocationID

		for j := 1; j < len(s); j++ {
			if s[j].LocationID == locationID {
				continue
			}

			if !yield(s[i:j:j]) {
				return
			}

			i = j
			locationID = s[j].LocationID
		}

		end := len(s)

		if i != end {
			yield(s[i:end:end])
		}
	}
}

// CalcContainersVolume - возвращает суммарный объём указанных контейнеров.
func CalcContainersVolume(s []dto.LocationStock) (total float64) {
	if len(s) == 0 {
		return 0
	}

	for i := range s {
		total += s[i].ContainerVolume * float64(s[i].ContainerQuantity)
	}

	return total
}
