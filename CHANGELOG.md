# Print Shop Back Changelog
Все изменения сервиса Print Shop Back будут документироваться на этой странице.

## 2024-03-15
### Changed
- Рефакторинг кода:
    - переименование `FactoryErrService*` -> `FactoryErrUseCase*`, `errService*` -> `errUseCase*`;
    - переименование интерфейсов `*Service` -> `*UseCase`;
    - замена методов `LoadOne` на `FetchOne`;
    - методы `Create`, `Insert` теперь возвращают ID записи;
    - схема БД переименована `gscatalog` -> `sample_catalog`;
- Вся мета информация об изображениях стала необязательной (`imageUrl`, и т.д.);
- Настройки `PageSizeMax` и `PageSizeDefault` вынесены в общие настройки модулей `ModulesSettings.General`;
- Парсер `SortPage` разделён на два: `ListSorter`, `ListPager`;
- Удалены неиспользуемые параметры запросов в каждом из модулей, отсортированы по алфавиту оставшиеся;
- В логгер добавлена поддержка `IsAutoCallerOnFunc`;


- Рефакторинг API документации:
    - Добавлены компоненты:
        - `App.Response.Model.BinaryAnyFile`;
        - `App.Response.Model.BinaryImage`;
        - `App.Response.Model.BinaryMedia`;
        - `App.Response.Model.JsonFile`;
        - `App.Response.Model.SuccessModifyItem`;
        - `App.Response.Model.TextFile.yaml`;
    - Доработка описания фильтрации, сортировки при получении списков записей;
    - Доработка описания ограничений при добавлении/обновлении записей;
    - Для всех модулей поля-идентификаторы описаны как отдельные сущности;

## 2024-02-05
### Changed
- Переименованы:
    - `datetime_created` -> `created_at`;
    - `datetime_updated` -> `updated_at`;
    - `modules.Options` -> `app.Options`;
- Создание модулей переехало в `factory/modules/*`;
- Большинство юнитов было преобразовано в модули, которые объединены доменами;

## 2024-01-30
### Changed
- Внедрён новый интерфейс логгера, добавлен режим трассировки запросов;
- Для многих методов добавлен параметр `ctx context.Context`;
- Заменён устаревший интерфейс `mrcore.EventBox` на `mrsender.EventEmitter`;
- Переименован `ServiceHelper` -> `UsecaseHelper`;
- Внедрены `mrlib.CallEachFunc`, `CloseFunc` для группового закрытия ресурсов;
- Переименован `CorrelationID` на `X-Correlation-ID`;
- Объекты конфигураций/опций теперь передаются по значению (`*Config` -> `Config`, `*Options` -> `Options`);
- Внедрён `oklog/run` для управления одновременным запуском нескольких серверов (http, grpc)
- Добавлены методы для создания и инициализации всех глобальных настроек приложения
  (`CreateAppEnvironment`, `InitAppEnvironment`);
- Теперь модули собираются в рамках отдельных серверов (см. `factory.NewRestServer`);
- Изменены некоторые переменные окружения:
    - удалён `APPX_LOG_PREFIX`;
    - добавлен `APPX_LOG_TIMESTAMP=RFC3339|RFC3339Nano|DateTime|TimeOnly` (формат даты в логах);
    - добавлен `APPX_LOG_JSON=true|false` (вывод логов в json формате);
    - добавлен `APPX_LOG_COLOR=true|false` (использование цветного вывода логов в консоле);
    - переименованы:
        - `APPX_SERVICE_LISTEN_TYPE` -> `APPX_SERVER_LISTEN_TYPE`;
        - `APPX_SERVICE_LISTEN_SOCK` -> `APPX_SERVER_LISTEN_SOCK`;
        - `APPX_SERVICE_BIND` -> `APPX_SERVER_LISTEN_BIND`;
        - `APPX_SERVICE_PORT` -> `APPX_SERVER_LISTEN_PORT`;

## 2024-01-25
### Added
- Внедрены парсеры на основе интерфейсов `mrserver.RequestParserFile` и
  `mrserver.RequestParserImage` для получения файлов и изображений из `multipart` формы.
    - заменено `mrreq.File` -> `ht.parser.FormImage`;
    - в `CompanyPageLogoService` изменён тип `mrtype.File` -> `mrtype.Image`;

### Changed
- Переименовано `ConvertImageMetaToInfo` -> `ImageMetaToInfoPointer`;

### Removed
- `mrserver.RequestParserPath` удалён вместо него используется
  `mrserver.RequestParserString` и `mrserver.RequestParserParamFunc`;

## 2024-01-22
### Changed
- Расформирован объект `ClientContext` и его одноименный интерфейс, в результате:
    - Изменена сигнатура обработчиков с `func(c mrcore.ClientContext)` на `func(w http.ResponseWriter, r *http.Request) error`;
    - С помощью интерфейсов `RequestDecoder`, `ResponseEncoder` можно задавать различные форматы
      принимаемых и отправляемых данных (сейчас реализован только формат `JSON`);
    - Запросы обрабатываются встраиваемыми в обработчики объектов `mrparser.*` через интерфейсы:
      `mrserver.RequestParserPath`, `RequestParser`, `RequestParserItemStatus`, `RequestParserKeyInt32`,
      `RequestParserSortPage`, `RequestParserUUID`, `RequestParserValidate`;
    - Ответы отправляются встраиваемыми в обработчики объекты `mrresponse.*` через интерфейсы:
      `mrserver.ResponseSender`, `FileResponseSender`, `ErrorResponseSender`;
    - Вместо метода `Validate(structRequest any)` используется объект `mrparser.Validator`;
- Произведены следующие замены:
    - `HttpController.AddHandlers` -> `Handlers() []HttpHandler`.
      Убрана зависимость контроллера от роутера и секции.
      Для установки стандартных разрешений добавлены следующие методы
      `mrfactory.WithPermission`, `mrfactory.WithMiddlewareCheckAccess`;
    - `ModulesAccess` -> `AccessControl` (`modules_access` -> `access_control`) и добавлен интерфейс `mrcore.AccessControl`;
    - `ClientSection` -> `AppSection` (`client_section` -> `app_section`) удалена зависимость от `AccessControl`;
- При внедрении новой версии библиотеки `go-sysmess` было заменено:
    - `mrerr.FieldErrorList` -> `CustomErrorList`;

## 2024-01-19
### Changed
- В БД `enum` типы заменены на `int2` и удалены. Доработано, чтобы `enum` типы сохранялись в виде `int`;
- Код получения файла в обработчике заменён на `mrreq.File`;
- Переименованы методы `checkForm`, `checkLaminate` и подобные в `usecase` на более абстрактный `checkItem`
  (проверяет возможность добавления, сохранения записи);

## 2024-01-17
### Added
- Для каждой секции добавлены настройки `AuthSecret` и `AuthAudience`;
- Добавлены системные обработчики (`RegisterSystemHandlers`);
- Добавлена фильтрация полей справочников (например `Custom.Request.Query.Filter.Density*`);
- Добавлено поле `Config.AppStartedAt` для отслеживания времени запуска сервиса;

### Changed
- Переработана документация `OpenAPI`, все API разделены по секциям, настроена сборка документов.
  Добавлены новые `OpenAPI` компоненты: `App.Field*`, `App.Enum*`, `App.Measure*`,
  `App.Response.Model.*`, и другие;
- Сущности поделены на модули, каждый модуль имеет собственные настройки (`Options`),
  которые назначаются фабрикой сервиса на основе общего конфига;
- Поле `companies_pages.logo_path` заменено на `companies_pages.logo_meta` типа `jsonb`,
  в котором теперь хранится мета информация об изображении;
- В конфиге для всех таймаутов заменён тип `int32` на `time.Duration`;

## 2023-09-21
### Changed
- Обновлены зависимости библиотеки;
- Фиксация зависимостей инфраструктуры;
- Заменён адаптер `*mrpostgres.ConnAdapter` на интерфейс `mrstorage.DbConn`;
- Заменены tabs на пробелы в коде;
- Добавлен модуль страница компании с возможностью загрузки логотипа;
- Добавлен модуль получения файлов;

## 2023-09-17
### Changed
- Исправлены коды возврата в REST API для некоторых методов;
- Все объекты, которые создаются при запуске сервиса перенесены в пакет factory;

### Fixed
- Добавлено закрытие операции получения записей из БД;

## 2023-09-13
### Changed
- Все общие компоненты были вынесены в отдельные проекты: `go-sysmess`, `go-webcore`,
  `go-storage`, `go-components`, в связи с этим были полностью переработаны все связи проекта;
- Обновлены все версии библиотек, от которых зависит проект;

## 2023-08-28
### Added
- Добавлен mrlib.Helper для помощи запуска сервиса;
- Добавлен MiddlewareUserIp() для получения IP текущего пользователя;

### Changed
- Создание соединения postgres перенесено в отдельный класс фабрики;
- Переработаны константы URL адресов методов API;
- Доработано оформление кода логирования событий;
- Изменён интерфейс ClientData, метод SendResponseWithError упразднён;
- Доработана система обработки ошибок в контроллере;
- Метод mrpostgres.ConnAdapter.Close() приведён к стандартному интерфейсу;

## 2023-07-27
### Changed
- Изменена валидация поля article, теперь доступны все символы кроме пробельных;
- Переименованы `ErrHttpRequestEnumLen` -> `ErrHttpRequestParamLen` и
  `ErrHttpRequestParseEnum` -> `ErrHttpRequestParseParam` для обобщения;

### Fixed
- Добавлено забытое поле `Sides` в `entity.CatalogPaper`;
- Скорректированы значение валидаторов (min, max, lte)ж
- Добавлен пропущенный фильтр `App.Request.Query.ItemStatuses` в API документацию;
- `*Remove` -> `*RemoveURL`;
- В методах формирования списков заменён: `client.Query` -> `client.SqQuery`;

## 2023-07-23
### Add
- Добавлена возможность обновления только указанных полей структуры;
- Добавлено несколько именованных ошибок: `ErrInternalInvalidType`,
  `ErrInternalInvalidData`, `ErrDataContainer`;
- Добавлено `mrentity.EmptynullString` для сохранения NULL значений в БД;

### Changed
- В метод SqUpdate перенесена логика проверки обновления записи;

### Fixed
- Поправлена регулярка для поиска параметров в сообщении, теперь в параметрах могут быть цифры; 

## 2023-07-22
### Add
- Добавлены дополнительные ограничения в таблицы БД;
- Для модельных объектов добавлены константные имена `const ModelName*`,
  которое используется при логировании событий и ошибок связанных с этими компонентами;
- Добавлена структура mrerr.Arg для передачи какой-либо информации в логи;
- Добавлено несколько именованных ошибок: `ErrInternalNilPointer`, `ErrStorageFetchedInvalidData`;
- Добавлены методы работы с БД в связке с `squirrel` (см. `conn_squirrel.go`)

### Changed
- Доработан механизм сортировки элементов, он был оформлен в компонент `ItemOrderer`.
  Его задачи: организовывать порядок следования элементов конкретного списка и
  перемещать эти элементы в рамках этого списка;

### Removed
- Удалён метод NewWithData, вместо него нужно использовать `New(mrerr.Arg{...})`;

## 2023-07-16
### Add
- Добавлен paramName в модуль `form_data`;
- Добавлен компонент для компиляции структуры, на основе которой формируются
  пользовательские интерфейсы для калькуляции продукции;
- В интерфейс `ClientData` добавлен метод `WithContext`;
- Добавлена именованная ошибка `ErrInternalParseData`;

### Changed
- GO_TOOLS утилиты теперь не устанавливаются по умолчанию;
- В ошибке, формируемой во wrapError, теперь выводится информация о том,
  в каком методе разработчика она произошла, а не там, где её породила библиотека ядра;
- Во всех `regexp.MustCompile` добавлены raw strings во избежание двойного экранирования символов;

### Fixed
- Исправлен баг при сохранении id формы в контекст;
- Исправлен баг при добавлении поля формы, проверка `paramName` проходила раньше,
  чем этот параметр заполнялся из шаблона;
- Добавлен omitempty в описание валидации необязательных полей;
- Исправлен синтаксис sql в некоторых запросах;
- Исправлен баг с защитой от бесконечного цикла при формировании ошибки;
- В `FormElementOrderer` добавлены пропущенные проверки ошибок;

## 2023-07-14
### Added
- Доработан проект для его развёртывания утилитой Mrcmd;
- Добавлена ещё одна настройка переменной MRCMD_SHARED_PLUGINS_ENABLED
  для альтернативного варианта разработки;

### Changed
- Переименован пакет `calc-user-data-back-adm -> print-shop-back`;

## 2023-07-10
### Added
- Внедрена библиотека `squirrel` для формирования SQL запросов;

## 2023-07-10
### Added
- Внедрена библиотека `squirrel` для формирования SQL запросов;

## 2023-07-07
### Added
- Произведён рефакторинг механизма обработки ошибок, добавлена фабрика для создания ошибок;

## 2023-05-13
### Changed
- Доработан механизм обработки ошибок, logger и валидация внешних данных;

## 2023-05-12
### Added
- Загружена версия проекта с несколькими реализованными методами API аутентификации;