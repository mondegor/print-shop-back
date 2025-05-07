
## Метод расчёта спуска полос

POST http://api.print-robot.ru/v1/calculations/algo/rect/imposition


Надо в description API более подробные описания сделать в Response:
"layout": {
"width": 1,
"height": 1
}
description:
layout - минимально достаточный формат для размещения всех элементов с учетом поля itemDistance.

Небходимо добавить в garbage поля между элементами itemDistance:
Горизонт полоска: itemDistance[0] * layout.width
Вертикальная полоска: itemDistance[1] * layout.height
Пересечение полосок количество: fragments.byWidth + fragments.byHeight - 2
garbageDistance = (itemDistance[0] * layout.width + itemDistance[1] * layout.height) * (fragments.byWidth + fragments.byHeight - 2)

Если есть остаток с поворотом, то там будет свой размер layout и fragments.

Test:
{
"itemFormat": "50x50",
"itemDistance": "5x8",
"outFormat": "400x300",
"disableRotation": false,
"useMirror": false
}

{
"layout": {
"width": 0.38,
"height": 0.282
},
"fragments": [
{
"byWidth": 7,
"byHeight": 5
}
],
"total": 35,
"garbage": 0.01284
}

garbageDistance = (5 * 380 + 8 * 282) * (7 + 5 - 2) / 1_000_000 = 0.04156 кв.м
garbage = garbage + garbageDistance

Request 1 помещается ровно 100 штук без двойных резов получаем "garbage": 0:
````
{
  "itemFormat": "10x10",
  "itemDistance": "0x0",
  "outFormat": "100x100",
  "disableRotation": false,
  "useMirror": false
}
````
Response 1:
````
{
    "layout": {
        "width": 0.1,
        "height": 0.1
    },
    "fragments": [
        {
            "byWidth": 10,
            "byHeight": 10
        }
    ],
    "total": 100,
    "garbage": 0
}
````

(-) Request 2 есть двойные резы по вертикале - "garbage" не должен быть 0
поля, которые получаются между элементами надо считать как garbage:
+ Как посчитать garbage?
````
{
  "itemFormat": "10x10",
  "itemDistance": "5x0",
  "outFormat": "100x100",
  "disableRotation": false,
  "useMirror": false
}
````
Response 2:
````
{
    "layout": {
        "width": 0.1,
        "height": 0.1
    },
    "fragments": [
        {
            "byWidth": 7,
            "byHeight": 10
        }
    ],
    "total": 70,
    "garbage": 0
}
````

Request 3 проверка зеркально по вертикальной оси false:
````
{
  "itemFormat": "10x5",
  "itemDistance": "0x0",
  "outFormat": "90x100",
  "disableRotation": false,
  "useMirror": false
}
````
Response 3:
````
{
    "layout": {
        "width": 0.09,
        "height": 0.1
    },
    "fragments": [
        {
            "byWidth": 9,
            "byHeight": 20
        }
    ],
    "total": 180,
    "garbage": 0
}
````

Request 4 проверка зеркально по вертикальной оси true:
````
{
  "itemFormat": "10x5",
  "itemDistance": "0x0",
  "outFormat": "90x100",
  "disableRotation": false,
  "useMirror": true
}
````
Response 4:
````
{
    "layout": {
        "width": 0.08,
        "height": 0.1
    },
    "fragments": [
        {
            "byWidth": 8,
            "byHeight": 20
        }
    ],
    "total": 160,
    "garbage": 0.001
}
````

Request 5 граничное значение - itemFormat == outFormat:
````
{
  "itemFormat": "100x100",
  "itemDistance": "0x0",
  "outFormat": "100x100",
  "disableRotation": false,
  "useMirror": false
}
````
Response 5:
````
{
    "layout": {
        "width": 0.1,
        "height": 0.1
    },
    "fragments": [
        {
            "byWidth": 10,
            "byHeight": 10
        }
    ],
    "total": 100,
    "garbage": 0
}
````

Request 6 negative - itemFormat > outFormat:
````
{
  "itemFormat": "101x100",
  "itemDistance": "0x0",
  "outFormat": "100x100",
  "disableRotation": false,
  "useMirror": false
}
````
Response 6:
````
{
    "layout": {
        "width": 0,
        "height": 0
    },
    "fragments": null,
    "total": 0,
    "garbage": 0
}
````

Request 7 negative - itemFormat > outFormat:
````
{
  "itemFormat": "101x100",
  "itemDistance": "0x0",
  "outFormat": "100x100",
  "disableRotation": false,
  "useMirror": false
}
````
Response 7:
````
{
    "layout": {
        "width": 0,
        "height": 0
    },
    "fragments": null,
    "total": 0,
    "garbage": 0
}
````

Request 8:
````
{
  "itemFormat": "100x100",
  "itemDistance": "1x0",
  "outFormat": "100x100",
  "disableRotation": false,
  "useMirror": false
}
````
Response 8:
````
{
    "layout": {
        "width": 0.1,
        "height": 0.1
    },
    "fragments": [
        {
            "byWidth": 1,
            "byHeight": 1
        }
    ],
    "total": 1,
    "garbage": 0
}
````

Request 9 помещается ровно 16 элементов:
````
{
  "itemFormat": "100x100",
  "itemDistance": "0x0",
  "outFormat": "400x400",
  "disableRotation": false,
  "useMirror": false
}
````
Response 9:
````
{
    "layout": {
        "width": 0.4,
        "height": 0.4
    },
    "fragments": [
        {
            "byWidth": 4,
            "byHeight": 4
        }
    ],
    "total": 16,
    "garbage": 0
}
````

(-) Request 10. помещается 4 элемента, но должно быть 9
это случается при outFormat кратно 3 или 6, например "outFormat": "300x300", "600x600", "1200x900":
````
{
  "itemFormat": "100x100",
  "itemDistance": "0x0",
  "outFormat": "300x300",
  "disableRotation": false,
  "useMirror": false
}
````
Response 10:
````
{
    "layout": {
        "width": 0.2,
        "height": 0.2
    },
    "fragments": [
        {
            "byWidth": 2,
            "byHeight": 2
        }
    ],
    "total": 4,
    "garbage": 0.05
}
````

Request 11. симметрия вдоль короткой стороны:
````
{
  "itemFormat": "100x100",
  "itemDistance": "0x0",
  "outFormat": "700x500",
  "disableRotation": false,
  "useMirror": true
}
````
Response 11:
````
{
    "layout": {
        "width": 0.6,
        "height": 0.5
    },
    "fragments": [
        {
            "byWidth": 6,
            "byHeight": 5
        }
    ],
    "total": 30,
    "garbage": 0.05
}
````

Request 12. симметрия вдоль длинной стороны:
````
{
  "itemFormat": "100x100",
  "itemDistance": "0x0",
  "outFormat": "500x700",
  "disableRotation": false,
  "useMirror": true
}
````
Response 12:
````
{
    "layout": {
        "width": 0.4,
        "height": 0.7
    },
    "fragments": [
        {
            "byWidth": 4,
            "byHeight": 7
        }
    ],
    "total": 28,
    "garbage": 0.07
}
````

Request 13. возможно разложение остатка с поворотом элемента:
````
{
  "itemFormat": "100x50",
  "itemDistance": "0x0",
  "outFormat": "290x200",
  "disableRotation": false,
  "useMirror": false
}
````
Response 13:
````
{
    "layout": {
        "width": 0.25,
        "height": 0.2
    },
    "fragments": [
        {
            "byWidth": 2,
            "byHeight": 4
        },
        {
            "byWidth": 2,
            "byHeight": 1
        }
    ],
    "total": 10,
    "garbage": 0.008
}
````

Request 14. нет itemDistance, помещается 2х2:
````
{
  "itemFormat": "50x50",
  "itemDistance": "0x0",
  "outFormat": "100x100",
  "disableRotation": false,
  "useMirror": true
}
````
Response 14:
````
{
    "layout": {
        "width": 0.1,
        "height": 0.1
    },
    "fragments": [
        {
            "byWidth": 2,
            "byHeight": 2
        }
    ],
    "total": 4,
    "garbage": 0
}
````

Request 15. есть itemDistance, помещается 1х2, но т.к. useMirror = true, получаем "total": 0:
````
{
  "itemFormat": "50x50",
  "itemDistance": "1x0",
  "outFormat": "100x100",
  "disableRotation": false,
  "useMirror": true
}
````
Response 15:
````
{
    "layout": {
        "width": 0,
        "height": 0
    },
    "fragments": null,
    "total": 0,
    "garbage": 0
}
````