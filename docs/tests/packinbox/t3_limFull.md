
## коробка заполнена ровно по высоте (граничное) 

{
  "product": {
    "format": "297x210",
    "thickness": 150,
    "weightM2": 100,
    "quantity": 1960
  },
  "box": {
    "format": "400x300x300",
    "thickness": 3000, 
    "margins": "0x0x0",
    "weight": 350,
    "maxWeight": 16150
  }
}  

Расчет:
boxBottomFormat = 394 х 294 мм
maxProductQuantityInStack = (300 - 6)/ 0.15 = 1960 шт
impResult.Total = 1
maxProductQuantityInBox = 1 * 1960 = 1960 шт
totalBoxQuantity = 1
restProductQuantity = 1960 - 1 * 1960 = 0 шт
boxesInnerVolume = 394 * 294 * 294 * 1 = 34_055_784 
boxesVolume = 400 * 300 * 300 * 1 = 36_000_000
productVolume = 210 * 297 * 0.15 * 1960 = 18_336_780
Weight = 0.210 * 0.297 * 1960 * 100 + 350 = 12_575 г
UnusedVolumePercent = 100 - 100 * 18_336_780 / 34_055_784  = 46.16 %