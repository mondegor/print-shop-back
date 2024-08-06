# Задать значение изделий в упаковке (листов в пачке) или макс заполнение коробок

## Два варианта заполнения коробок

В форме UI выбираем одну опцию (радиобокс):
1. кол-во листов в пачке: [поле ввода целое число от 1 до 100.000]
2. заполнить максимально коробки.

## 1. Вариант задать значение в форме UI кол-во листов в пачке

Вводим значение в поле "кол-во листов в пачке".
Расчитываем продукцию с параметрами, где толщина равна высоте пачки.

Входные параметры:
- detailThickness [mm] толщина детали (бумага + ламинат, брошюры)
- detailQuantity [int] продукции в пачке из UI
- Product.quantity [int] тираж

Исходящие параметры:
- resultDetailThickness - новая толщина детали [mm]
- resultDetailQuantity - кол-во пачек
- resultDetailThicknessRest - толщина пачки с остатком

Если detailQuantity > Product.quantity, то detailQuantity = Product.quantity

resultDetailThickness = detailThickness * detailQuantity

resultDetailQuantity = Product.quantity / detailQuantity, целое число. Если меньше 1, то

### Тест 1 пачка + остаток
- detailThickness = 0.12 [mm] толщина детали (бумага + ламинат, брошюры)
- detailQuantity = 250 [int] продукции в пачке из UI
- Product.quantity = 300 [int] тираж

resultDetailThickness = 0.12 * 250 = 30 mm
resultDetailThicknessRest = 
