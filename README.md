# LeGo: Go Application Bootstrap Constructor

> ⚠️ **NOTE:** this project is discontinued and have no any maintainer.
> 
> ⚠️ **Внимание!** Этот проект завершен и более не поддерживается.

### Мотивация

Для поднятия приложения часто приходится повторять одну и ту же рутинную работу: проверка конфигов, описание
http-сервера, оборачивание его метриками и чеками, подключение хранилищ, добавление логера, соединение этих всех
инфраструктурных вещей в приложении.

Эта библиотека поможет избавиться от этой рутины и поможет запустить бизнес-логику вашего сервиса в пару строчек кода.

Под капотом есть уже подготовленные http / grpc сервера с конфигами, метриками, логгерами, библиотеки для упрощения
подключения баз данных и работы с событиями

### Быстрый запуск

#### Генератор

```bash
curl -sfL https://git.io/JThPO | bash -s project-name
```

Алиас

```bash
alias lego:new="curl -sfL https://git.io/JThPO | bash -s"

lego:new project-name
```

#### Ручной запуск

```go
package main

import (
    "github.com/vseinstrumentiru/lego/v2/server"
    "github.com/vseinstrumentiru/lego/v2/transport/http"
)

// ...

func main() {
    type SomeConfig struct {
        // ...
        AnyKeyForHttp http.Config
        // ...
    }

    // will start http server with metrics and your app
    server.Run(yourpackage.App{}, server.ConfigOption(&SomeConfig{}))
}
```

## Базовые принципы

### Провайдер

Провайдер - это функция / метод - конструктор, который возвращает один или несколько типов и так же может возвращать
ошибку создания. Все возвращаемые **типы должны быть конкретными** (структура / ссылка на структуру, объявленный
интерфейс) и **не повторяться у всех провайдеров** (только один зарегистрированный провайдер может
создавать один и тот же тип).

Все арументы провайдера подружаются автоматически через DI, поэтому не нужно объявлять у провайдера аргументы с типами,
которые не зарегистрированы в DI или примитивы.

### Инстанс App

App - это первый аргумент, для запуск сервера. Структура App (может быть любое название) используется для подготовки
приложения и соединения его основных компонентов.

Все public-поля структуры App подгружаются с помощью DI.

#### Pipeline создания App

1. Метод ```Providers() []interface{}```: если у App есть этот метод, в нём можно перечислить все внутренние
   конструкторы приложения, и DI зарегистрирует их.
1. Методы ```Provide(...)``` и ```Provide[SomeSuffix...](...) (...)```: используются для быстрого объявления
   конструктора типа, объявления и подключения фонового процесса, который запустится вместе со всем приложением.
1. Создание инстанса App с подгрузкой его публичных полей
1. Методы ```Configure(...)``` и ```Configure[SomeSuffix...](...)```: используются для конфигурирования приложения:
   настройки http-роутов, сервисов, событий. **Не нужно** здесь запускать фоновые процессы.

### Конфиг приложения

Конфиг приложения - это вольная структура в которую подгрузятся данные из переменных окружения и/или из файла
конфигурации. Рутовый префикс для конфигов ```APP_``` для переменных окружения и ```app.``` для файлов конфигурации
Переменные окружения должны разделяться ```_```: для ```сfg.StructOne.StructTwo.Field``` будет переменная
окружения ```APP_STRUCTONE_STRUCTTWO_FIELD```

Все вложенные структуры регистрируются в DI и могут использоваться отдельно от базовой структуры в методах конфигурации
и провайдерах.

Для запуска встроенных компонентов LeGo, небходимо здесь же добавить поля с типами конфигурации этих компонентов: т.е.
если необходимо подключить http-сервер - в структуре конфигурации должно быть поле с типом http.Config

В DI необходимо использовать ссылку на структуры или вложенные структуры, даже если они не были объявлены таковыми в
структуре конфигурации.

Структура конфигурации и её подструктуры могут содержать (но не обязаны) два метода, которые выполнятся при загрузке
конфигурации:

#### ```SetDefaults(env config.Env)```

Используется для установки значений по умолчанию, алиасов и добавления флагов запуска приложения. Выполняются поочерёдно
от дочерних структур к родительской (т.е. родительская структура может переопределить параметры по умолчанию для детей).

#### ```Validate() error```

Используется для проверки конфига на правильность его заполнения. Выполняется после загрузки конфигурации.

## Компоненты

### Multilog Logger

Многоуровневый логгер приложения.

```go
import "github.com/vseinstrumentiru/lego/v2/multilog"
```

Настройка уровня логирования:

```go
import "github.com/vseinstrumentiru/lego/v2/multilog"

type Config struct {
    Multilog multilog.Config
}
```

```yaml
app:
  multilog:
    # trace | debug | info | warn | error
    level: error
```

Может отправлять логи в:

1. stdout
    ```go
    import "github.com/vseinstrumentiru/lego/v2/multilog/log"

    type Config struct {
        Log log.Config
    }
    ```
    ```yaml
    app:
      log:
        # Вывод в красивом виде с поддержкой цветов, удобно для локальной разработки, но НЕ ПРОД (по умолчанию - false, вывод в JSON) (true | false)
        color: false
        # Не пробрасывание сообщения в другие обработчики, если обработан этим
        stop: false
        # Глубина для отображения файла из stack trace (работает только с color = true)
        depth: -1
    ```
1. Sentry
    ```go
    import "github.com/vseinstrumentiru/lego/v2/multilog/sentry"
    type Config struct {
        Sentry sentry.Config
    }
    ```
    ```yaml
    app:
      sentry:
        # Урл от проекта в Sentry
        adrr: "http://some.sentry.com/fsfdjlkfhk3f@fsjdkfjhh"
        # Не пробрасывание сообщения в другие обработчики, если обработан этим
        stop: false
        # С каким уровнем обрабатывать сообщения (trace | debug | info | warn | error)
        level: error
    ```
1. NewRelic (в процессе)
    ```go
    import "github.com/vseinstrumentiru/lego/v2/metrics/exporters"
    type Config struct {
        NewRelic exporters.NewRelic
    }
    ```
    ```yaml
    app:
      newrelic:
        # Вкл / Выкл
        enabled: false
        # Использование телеметрии и трассировки
        telemetryEnabled: false
        # Ключ приложения
        key: "nerelic-license-key"
    ```
