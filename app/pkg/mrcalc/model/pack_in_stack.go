package model

type (
	// PackInStack - результат работы алгоритма AlgoSheet.
	PackInStack struct {
		FullProduct ProductStack // полностью заполненная пачка листов
		RestProduct Product      // пачка с остатком листов
	}
)

// TotalQuantity - возвращает общее количество пачек с изделиями.
func (o *PackInStack) TotalQuantity() uint64 {
	if o.RestProduct.Weight == 0 {
		return o.FullProduct.Quantity
	}

	return o.FullProduct.Quantity + 1
}

// TotalVolume - возвращает общий объём пачек с изделиями.
func (o *PackInStack) TotalVolume() float64 {
	if o.RestProduct.Weight == 0 {
		return o.FullProduct.TotalVolume()
	}

	return o.FullProduct.TotalVolume() + o.RestProduct.Format.Volume()
}

// TotalWeight - возвращает общий вес пачек с изделиями.
func (o *PackInStack) TotalWeight() float64 {
	if o.RestProduct.Weight == 0 {
		return o.FullProduct.TotalWeight()
	}

	return o.FullProduct.TotalWeight() + o.RestProduct.Weight
}
