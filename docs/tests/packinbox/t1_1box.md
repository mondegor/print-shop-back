
## Заполнена 1 коробка

{
  "product": {
    "format": "210x297",
    "thickness": 300,
    "weightM2": 300,
    "quantity": 1000
  },
  "box": {
    "format": "630x320x340",
    "thickness": 500, # !5000
    "margins": "0x0x0",
    "weight": 500,
    "maxWeight": 25000
  }
}

Расчет:

boxBottomFormat = 620 х 310 мм
maxProductQuantityInStack = (340 - 10)/ 0.3 = 1100 шт
impResult.Total = 2
maxProductQuantityInBox = 2 * 1100 = 2200 шт
totalBoxQuantity = 1
restProductQuantity = 1000 - 0 * 2200 = 1000
BoxVolumeInternal = 620 * 310 * 330 * 1 = 63_426_000
boxVolumeExternal = 630 * 320 * 340 * 1 = 68_544_000
productVolume = 210 * 297 * 0.3 * 1000= 18_711_000
Weight = 0.210 * 0.297 * 1000 * 300 + 500 = 19211 г
UnusedVolumePercent = 100 - 100 * 18_711_000 / 63_426_000 = 70.5 %
