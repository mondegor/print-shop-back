# Описание Print Shop Back v0.6.1
Этот репозиторий содержит описание сервиса Print Shop Back.

## Статус сервиса
Сервис находится в стадии разработки.

## Описание сервиса
Web сервис для расчёта стоимости и времени изготовления продукции в конкретной типографии.

## REST API документация
- https://github.com/mondegor/print-shop-back/blob/master/docs/

## Разворачивание, установка и запуск сервиса

### Разворачивание сервиса
> Перед разворачиванием сервиса необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

- Выбрать рабочую директорию, где должен быть расположен сервис
- `mkdir print-shop-back && cd print-shop-back` // создать и перейти в директорию проекта
- `git clone -b latest git@github.com:mondegor/print-shop-back.git .`
- `cp .env.dist .env`
- `mrcmd state` // проверка состояния сервиса
- `mrcmd config` // проверка установленных переменных сервиса

> Более подробную информацию по использованию утилиты Mrcmd смотрите [здесь](https://github.com/mondegor/mrcmd#readme).

### Установка сервиса и его первый запуск
- `mrcmd docker ps` // убеждаемся, что Docker daemon запущен
- `mrcmd install`
- `mrcmd start`
- `mrcmd docker-compose ps` // проверка всех запущенных ресурсов сервиса
- `mrcmd go-migrate up` // загрузка дампа с данными в БД
- `mrcmd go logs` // проверка, что основной сервис запущен

### Запуск и остановка сервиса
- `mrcmd start`
- `mrcmd stop`

### Остановка сервиса и удаление всех его установленных ресурсов
- `mrcmd uninstall`

### Часто используемые команды
- `mrcmd help` - помощь в контексте текущего сервиса;
- `mrcmd state` - общее состояние текущего сервиса;
- `mrcmd docker-compose ps` - текущее состояние запущенных ресурсов;
- `mrcmd docker-compose logs` - логи всех запущенных ресурсов;
- `mrcmd go help` - все команды сервиса go;
- `mrcmd go-migrate help` - все команды сервиса go-migrate;
- `mrcmd postgres help` - все команды сервиса postgres;