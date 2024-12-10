# Описание Print Shop Back v0.17.4
Этот репозиторий содержит описание сервиса Print Shop Back.

## Статус сервиса
Сервис находится в стадии разработки.

## Описание сервиса
Web сервис для расчёта стоимости и времени изготовления продукции.

Подробнее смотри [документацию проекта](./docs/README.md)

> Перед запуском консольных скриптов сервиса необходимо скачать и установить утилиту Mrcmd.\
> Инструкция по её установке находится [здесь](https://github.com/mondegor/mrcmd#readme)

### Команды для сборки API документации v0.4.4
- `mrcmd openapi help` - помощь по командам плагина openapi;
- `mrcmd openapi build-all` - сборка документации всех API;

### Примеры запуска сборки документации из консоли Windows:
- GitBash (cmd): `"C:\Program Files\Git\git-bash.exe" --cd=d:\mrwork\print-shop-back mrcmd openapi build-all`;
- WSL (PowerShell): `cd D:\workdir\print-shop-back; wsl -d Ubuntu-20.04 -e mrcmd openapi build-all`;

## Разворачивание, установка и запуск сервиса

### Разворачивание сервиса
- Выбрать рабочую директорию, где должен быть расположен сервис
- `mkdir print-shop-back && cd print-shop-back` // создать и перейти в директорию проекта
- `git clone git@github.com:mondegor/print-shop-back.git .`
- `cp env.dist-stand .env`
- `mrcmd state` // проверка состояния сервиса
- `mrcmd config` // проверка установленных переменных сервиса
- `mrcmd go deps` // загрузка зависимостей проекта
- Для работы утилит `gofumpt`, `goimports` (и других) необходимо в `.env` проверить
  значения переменных `GO_TOOLS_INSTALL_*` и запустить `mrcmd go install-tools`

### Установка сервиса и его первый запуск
- `mrcmd docker ps` // убеждаемся, что Docker daemon запущен
- `mrcmd install`
- `mrcmd start`
- `mrcmd docker-compose ps` // проверка всех запущенных ресурсов сервиса
- `mrcmd go logs` // проверка запуска основного сервиса

### Запуск и остановка сервиса
- `mrcmd start`
- `mrcmd stop`

### Остановка сервиса и удаление всех его установленных ресурсов
- `mrcmd uninstall`

### Часто используемые команды
- `mrcmd help` // помощь в контексте текущего сервиса;
- `mrcmd state` // общее состояние текущего сервиса;
- `mrcmd docker-compose conf` // отображает список `.yaml` файлов из которых собрана конфигурация и саму конфигурацию;
- `mrcmd docker-compose ps` // текущее состояние запущенных ресурсов;
- `mrcmd docker-compose logs` // логи всех запущенных ресурсов;
- `mrcmd go-migrate help` // все команды сервиса go-migrate;
- `mrcmd go-migrate up` // загрузка дампа с данными в БД;
- `mrcmd postgres help` // все команды сервиса postgres;
- `mrcmd go help` // выводит список всех доступных go команд;
- `mrcmd go gofumpt-fix` // исправляет форматирование кода (`gofumpt -l -w -extra ./`);
- `mrcmd go goimports-fix` // исправляет imports, если это требуется (`goimports -d -local ${GO_IMPORTS_LOCAL_PREFIXES} ./`);
- `mrcmd go test` // запускает тестов библиотеки;
- `mrcmd go test-report` // запускает тесты библиотеки с формированием отчёта о покрытии кода (`test-coverage-full.html`);
- `mrcmd golangci-lint check` // запускает линтеров для проверки кода (на основе `.golangci.yaml`);
- `mrcmd plantuml build-all` // генерирует файлы изображений из `.puml` [подробнее](https://github.com/mondegor/mrcmd-plugins/blob/master/plantuml/README.md#%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%B0-%D1%81-%D0%B4%D0%BE%D0%BA%D1%83%D0%BC%D0%B5%D0%BD%D1%82%D0%B0%D1%86%D0%B8%D0%B5%D0%B9-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%B0-markdown--plantuml);

> Более подробную информацию по использованию утилиты Mrcmd
> смотрите [здесь](https://github.com/mondegor/mrcmd#readme).

#### Короткий вариант выше приведённых команд (Makefile)
- `make build` // аналог `mrcmd install`
- `make build-api` // аналог `mrcmd openapi build-all`
- `make deps` // аналог `mrcmd go deps`
- `make migrate` // аналог `mrcmd go-migrate up`
- `make generate` // аналог `mrcmd go generate`
- `make fmt` // аналог `mrcmd go fmt`
- `make fmti` // аналог `mrcmd go fmti`
- `make lint` // аналог `mrcmd golangci-lint check`
- `make test` // аналог `mrcmd go test`
- `make test-report` // аналог `mrcmd go test-report`
- `make plantuml` // аналог `mrcmd plantuml build-all`
- `make app-conf` // аналог `mrcmd config && mrcmd docker-compose conf`
- `make app-start` // аналог `mrcmd start`
- `make app-state` // аналог `mrcmd docker-compose ps`
- `make app-logs` // аналог `mrcmd docker-compose logs`
- `make app-stop` // аналог `mrcmd stop`

> Чтобы расширить список команд, необходимо создать Makefile.mk и добавить
> туда дополнительные команды, все они будут добавлены в единый список команд make утилиты.

## Панели управления развёрнутой инфраструктуры
- TRAEFIK: http://traefik.local/ (admin 12345678);
- API: http://api.print-shop.local/;
- HEALTH: http://print-shop.internal/health;
- INFO: http://print-shop.internal/system-info;
- METRICS: http://print-shop.internal/metrics;

### Использование локальных доменов
Необходимо в hosts добавить следующие записи:
- `127.0.0.1 traefik.local`
- `127.0.0.1 print-shop.local`
- `127.0.0.1 api.print-shop.local`
- `127.0.0.1 print-shop.internal`