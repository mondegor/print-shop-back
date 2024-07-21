
## коробка заполнена ровно по высоте (граничное) 

{
  "product": {
    "format": "210x297", # !change
    "thickness": 150,
    "weightM2": 100,
    "quantity": 1960
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
maxProductQuantityInStack = (300 - 6)/ 0.15 = 1960 шт
impResult.Total = 1
maxProductQuantityInBox = 1 * 1960 = 1960 шт
totalBoxQuantity = 1
restProductQuantity = 2900 - 1 * 2900 = 0 шт
BoxVolumeInternal = 394 * 296 * 296 * 1 = 69_041_408 
boxVolumeExternal = 400 * 300 * 300 * 1 = 72_000_000
productVolume = 210 * 297 * 0.1 * 3000= 18_711_000
Weight = 0.210 * 0.297 * 3000 * 100 + 350 = 19061 г
UnusedVolumePercent = 100 - 100 * 18_711_000 / 69_041_408  = 72.9 %