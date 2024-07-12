### Расчёт веса и объема листовок

```plantuml
@startuml
!theme mars

skinparam {
    MaxMessageSize 250
}

skinparam sequence {
    ParticipantPadding 125
    MessageAlign center
}

participant "User" as user order 10
participant "WebApp" as app order 20
participant "ServiceMaterial" as sm order 30
participant "ServiceBoxes" as sb order 40

app -> sm: запрос на получение фильтров материала
app <-- sm: Success 200
app -> app: отрисовка фильтров материалов

user -> app: заполняет данные о тираже и о материале
user -> app: нажимает кнопку рассчитать

app -> sm: запрос на получение толщины материала (на основе указанных фильтров материала)
app <-- sm: Success 200 + thickness


app -> sb: запрос на расчёт формулы: length, width, thickness, quantity, boxID
sb -> sb: валидация данных
sb -> sb: расчёт по формуле
app <-- sb: Success 200 + App.Response.Model.Box[]
app -> app: формирование результата расчёта

@enduml
```