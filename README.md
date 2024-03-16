# Описание Print Shop Back v0.12.1
Этот репозиторий содержит описание сервиса Print Shop Back.

## Статус сервиса
Сервис находится в стадии разработки.

## Описание сервиса
Web сервис для расчёта стоимости и времени изготовления продукции в конкретной типографии.

> Перед запуском консольных скриптов сервиса необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

## REST API документация
- [API документация в формате OpenAPI/Swagger](./docs/api)
    - [AdminAPI](./docs/api/admin-api)
    - [PublicAPI](./docs/api/public-api)
- [Исходники, из которых собирается API документация](./docs/api-src)
- [Скрипты, которые собирают API документацию](./scripts/openapi)
- [Статья о спорных моментах REST API](https://habr.com/ru/articles/770226/)

> Более подробную информацию по OpenAPI смотрите ниже в разделе "Что такое OpenAPI?".

### Команды для сборки API документации
- `mrcmd openapi help` - помощь по командам плагина openapi;
- `mrcmd openapi build-all` - сборка документации всех API;

### Примеры запуска сборки документации из консоли Windows:
- GitBash (cmd): `"C:\Program Files\Git\git-bash.exe" --cd=d:\mrwork\tn-docs mrcmd openapi build-all`
- WSL (PowerShell): `cd D:\workdir\go-sample; wsl -d Ubuntu-20.04 -e mrcmd openapi build-all`;

## Разворачивание, установка и запуск сервиса

### Разворачивание сервиса

- Выбрать рабочую директорию, где должен быть расположен сервис
- `mkdir print-shop-back && cd print-shop-back` // создать и перейти в директорию проекта
- `git clone -b latest git@github.com:mondegor/print-shop-back.git .`
- `cp .env.dist .env`
- `mrcmd state` // проверка состояния сервиса
- `mrcmd config` // проверка установленных переменных сервиса

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

> Более подробную информацию по использованию утилиты Mrcmd
> смотрите [здесь](https://github.com/mondegor/mrcmd#readme).

## Что такое OpenAPI?
Из [OpenAPI Specification](https://github.com/OAI/OpenAPI-Specification):

> The OpenAPI Specification (OAS) defines a standard, programming language-agnostic interface
> description for HTTP APIs, which allows both humans and computers to discover and understand
> the capabilities of a service without requiring access to source code, additional documentation,
> or inspection of network traffic. When properly defined via OpenAPI, a consumer can understand
> and interact with the remote service with a minimal amount of implementation logic. Similar to
> what interface descriptions have done for lower-level programming, the OpenAPI Specification
> removes guesswork in calling a service.

### Описание OpenAPI спецификации на swagger.io
- [v3.0](https://swagger.io/specification/v3/)
- [v3.1](https://swagger.io/specification/)

### Просмотр и редактирование OpenAPI спецификации
- [JetBrains OpenAPI (Swagger) Editor](https://plugins.jetbrains.com/plugin/14837-openapi-swagger-editor)
- [Swagger Editor](https://editor.swagger.io/)
- [Insomnia](https://insomnia.rest/download)
- [OpenAPI.Tools](https://openapi.tools/)