# Добавление веса и толщищины изделия

## Надо изменить вес и толщину изделия, если есть в маршрутной карте:
- ламинация
- брошюровка


https://github.com/Printrobot/calc-api/blob/main/processes/ProcessLamination.ts

```
export type MaterialInput = {     
  // получаем из артикула материала
  thiknessFilmMicrometers: number  // Толщина ламината мкм
  lengthFilmRoll: number 
  widthFilmRoll: number
  weigthFilmRollKg: number  // вес рулона в кг
}

Для листовой продукции добавляем толщину ламината к изделию:
// меняем толщину и вес детали: 
  result.detail.thicknessMillimeters = process.detailThicknessMillimeters +
    material.thicknessFilmMicrometers / 1000 * production.laminationSides;

  result.detail.weigthGramsSqMeter = process.detailWeigthGramsSqMeter +
    roundDigits(1e6 * material.weigthFilmRollKg / (material.lengthFilmRoll * material.widthFilmRoll), 2);
```


## Общий алгоритм добавления толщины addThickness 

Входящие параметры:
- detailThickness - изначальная толщина детали [mm]
- addThickness - добавить толщины [mm]
- addQuantity - одинаковых деталей добавить [int]

Исходящие параметры:
- resultDetailThickness - новая толщина детали [mm]
resultDetailThickness = detailThickness + addThickness * addQuantity

## Общий алгоритм добавления веса addWeght 

Входящие параметры:
- detailWeight - изначальная толщина детали [g]
- addWeight - добавить вес одной детали[g]
- addQuantity - одинаковых деталей добавить [int]

Исходящие параметры:
- resultDetailWeight - новый вес детали [g]
resultDetailWeight = detailWeight + addWeight * addQuantity

===

### Листовая продукция. Если есть ламинация, то добавить вес и толщину ламината

Входящие параметры:
- detailThickness, толщина бумаги из справочника [mm]
- addThickness = толщинf ламината из справочника [mm]
- addQuantity = 1 или 2, с одной стороны или с двух сторон ламинация из формы UI [int]

- detailWeight - вес листовки [g]
- addWeight - добавить вес ламината [g]

- detailWidth ширина листовки
- detailHieght длина листовки
- material.weigthFilmRollKg вес роля ламината 
- material.lengthFilmRoll 
- material.widthFilmRoll

Вес листовок: 
function LeafletWeight (detailWidth, detailHeight, paperDensity) {
return detailWeight = detailWidth * detailHeight * paperDensity
}

addWeight = detailWidth * detailHeight * material.weigthFilmRollKg / (material.lengthFilmRoll * material.widthFilmRoll)

===  

**Тест1 Листовки**
- detailThickness = 0.2 [mm]
- addThickness = 0.032 [mm]
- addQuantity = 2 сторон ламинации
- detailWidth = 0.21 [m]
- detailHeight = 0.297 [m]
- paperDensity = 170 [g/m2]
- material.weigthFilmRollKg = 30_000 [g]
- material.lengthFilmRoll = 3000 [m]
- material.widthFilmRoll = 0.45 [m]

resultDetailThickness = detailThickness + addThickness * addQuantity
resultDetailThickness = 0.2 + 0.032 * 2 = 0.264 [mm]

detailWeight = detailWidth * detailHeight * paperDensity
detailWeight = 0.21 * 0.297 * 170 = 10.6029 [g]

resultDetailWeight = detailWeight + addWeight * addQuantity
addWeight = detailWidth * detailHeight * material.weigthFilmRollKg / (material.lengthFilmRoll * material.widthFilmRoll)

addWeight = 0.21 * 0.297 * 30_000 / (3000 * 0.45) = 1.386 [g]
resultDetailWeight = 10.6029 + 1.386 * 2 = 13.3749 [g]


### Если есть брошюровка, то к обложке добавить вес и толщину блока

Толщина передней и задней стороны обложки. Входящие параметры:
- detailThickness = 0, толщина еще не определена [mm]
- addThickness = добавить толщину передней и задней стороны обложки [mm]
- addQuantity = pages / 2, к-во листов, pages = 4 для обложки

К толщине обложки прибавить толщину блока:
- detailThickness = толщина обложки [mm]
- addThickness толщина бумаги для блока [mm]
- addQuantity = pages / 2, к-во листов, pages - страниц в блоке из формы UI 

=== 

**Тест2 Брошюры**

**Обложка:**
- detailThickness = 0.15 мм
- addThickness = 0.12 мм
- addQuantity = 2
- detailWidth = 0.21 [m]
- detailHeight = 0.297 [m]
- paperDensity = 250 [g/m2]

resultDetailThickness = 0.0 + 0.15 * 2 = 0.3 мм

resultDetailWeight = detailWeight + addWeight * addQuantity
detailWeight = detailWidth * detailHeight * paperDensity * 2

detailWeight = 0.21 * 0.297 * 250 * 2 = 31.185 [g] вес обложки

**Блок 48 стр:**
- detailThickness = 0.3 мм
- addThickness = 0.1 мм
- addQuantity = pages / 2 = 24
- detailWidth = 0.21 [m]
- detailHeight = 0.297 [m]
- paperDensity = 90 [g/m2]

resultDetailThickness = 0.3 + 0.1 * 24 = 2.7 мм

resultDetailWeight = detailWeight + addWeight * addQuantity
detailWeight = 31.185 [g] вес обложки
addWeight = detailWidth * detailHeight * paperDensity * pages / 2
addWeight = 0.21 * 0.297 * 90  = 5.6133 [g] вес одного листа блока

resultDetailWeight = 31.185 + 5.6133 * 24 = 165.9042 [g]

Если у обложки есть ламинация, то использовать увеличение веса и толщины обложки как для листовок
