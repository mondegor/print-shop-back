
## Метод расчёта максимальное количества размещений элементов на внешнем формате

Возвращает максимальное количество единиц указанного внутреннего формата, которое можно разместить во внешнем указанном формате.

POST http://api.print-robot.ru/v1/calculations/algo/rect/inside-on-outside-max

Request 1:
````
{
  "inFormat": "30x30",
  "outFormat": "300x300"
}
````
Response 1:
````
{
    "fragment": {
        "byWidth": 10,
        "byHeight": 10
    },
    "total": 100
}
````
