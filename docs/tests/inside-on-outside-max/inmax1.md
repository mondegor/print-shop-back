
## Метод расчёта количества размещений элементов на внешнем формате

Возвращает количество единиц указанного внутреннего формата, которое можно разместить по вертикали и горизонтали во внешнем указанном формате (без использования поворотов).

POST http://api.print-robot.ru/v1/calculations/algo/rect/inside-on-outside-quantity

Request 1 ровно один элемент помещается:
````
{
  "inFormat": "100x100",
  "outFormat": "100x100"
}
````
Response 1:
````
{
    "fragments": [
        {
            "byWidth": 1,
            "byHeight": 1
        }
    ],
    "total": 1
}
````

Request 2:
````
{
  "inFormat": "100x100",
  "outFormat": "200x200"
}
````
Response 2:
````
{
    "fragments": [
        {
            "byWidth": 2,
            "byHeight": 2
        }
    ],
    "total": 4
}
````

(-) Request 3 Должно быть 9, получилось 4:
````
{
  "inFormat": "100x100",
  "outFormat": "300x300"
}
````
Response 3:
````
{
    "fragments": [
        {
            "byWidth": 2,
            "byHeight": 2
        }
    ],
    "total": 4
}
````

Request 4 не помещается:
````
{
  "inFormat": "301x300",
  "outFormat": "300x300"
}
````
Response 4:
````
{
    "fragments": null,
    "total": 0
}
````

(-) Request 4 negative. задан размер inFormat == 0, надо обработать эту ситуацию:
````
{
  "inFormat": "0x300",
  "outFormat": "300x300"
}
````
Response 4:
````
{
    "title": "Внутренняя ошибка сервера",
    "details": "DebugInfo: errCode=errUseCaseIncorrectInputData; errKind=Internal; err={[67ab5dd2-90c77390] data={In:0x0.3 Out:0.3x0.3} is incorrect: in format is not valid: 0x0.3}",
    "request": "POST /v1/calculations/algo/rect/inside-on-outside-max",
    "time": "2025-02-11T14:25:22Z",
    "errorTraceId": "67ab5dd2-90c77390"
}
````

(-) Request 5 negative. задан размер outFormat == 0, надо обработать эту ситуацию:
````
{
  "inFormat": "300x300",
  "outFormat": "0x300"
}
````
Response 5:
````
{
    "title": "Внутренняя ошибка сервера",
    "details": "DebugInfo: errCode=errUseCaseIncorrectInputData; errKind=Internal; err={[67ab5dfd-b89f1abd] data={In:0.3x0.3 Out:0x0.3} is incorrect: out format is not valid: 0x0.3}",
    "request": "POST /v1/calculations/algo/rect/inside-on-outside-max",
    "time": "2025-02-11T14:26:05Z",
    "errorTraceId": "67ab5dfd-b89f1abd"
}
````