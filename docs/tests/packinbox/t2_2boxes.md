
## Заполнено 2 коробки

{
  "product": {
    "format": "297x210", 
    "thickness": 100,
    "weightM2": 100,
    "quantity": 3000
  },
  "box": {
    "format": "400x300x300",
    "thickness": 300, # !2000
    "margins": "0x0x0",
    "weight": 350,
    "maxWeight": 16150
  }
}  

Расчет:
boxBottomFormat = 394 х 296 мм
maxProductQuantityInStack = (300 - 10)/ 0.1 = 2900 шт
impResult.Total = 1
maxProductQuantityInBox = 1 * 2900 = 2900 шт
totalBoxQuantity = 2
restProductQuantity = 3000 - 1 * 2900 = 100 шт
boxesInnerVolume = 394 * 296 * 296 * 2 = 69_041_408 
boxesVolume = 400 * 300 * 300 * 2 = 72_000_000
productVolume = 210 * 297 * 0.1 * 3000= 18_711_000
Weight = 0.210 * 0.297 * 3000 * 100 + 350 = 19061 г
UnusedVolumePercent = 100 - 100 * 18_711_000 / 69_041_408  = 72.9 %
