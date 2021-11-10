# Библиотека конфигурации

## Функционал

Заполняет структуру данными из переменных окружения, файлов конфигурации (yaml, toml, json. dotenv), флагов запуска
приложения и валидирует её.

Реализован на базе:

- [github.com/spf13/viper](https://github.com/spf13/viper) - заполнение из файлов и переменных окружения
- [github.com/spf13/pflag](https://github.com/spf13/pflag) - работа с флагами запуска
- [github.com/go-playground/validator/v10](https://github.com/go-playground/validator) - валидация

## Быстрый старт

```shell
go get github.com/vseinstrumentiru/lego/v3/pkg/config
```

```go
import "github.com/vseinstrumentiru/lego/v3/pkg/config"
```

```go
package main

import "github.com/vseinstrumentiru/lego/v3/pkg/config"

type Config struct {
    Application struct {
        Name    string `validator:"required"`
        MajorVersion int `flag:",v,application version"`
    } `env:",squash"`
    HTTPPort int `validator:"gte=8000"`
}

func (c *Config) SetDefaults() {
    c.Application.MajorVersion = 1
}

func main() {
    cfg := Config{}

    if err := config.New(config.WithEnvPrefix("app")).Unmarshal(&cfg); err != nil {
		panic(err)
    }

	// ... run app with validated config
}
```

#### .env
```.env
APP_NAME=my-application
APP_HTTPPORT=8081
```

#### config.yaml
```yaml
name: my-application
httpPort: 8081
```

#### Запуск с флагами
```bash
app --major-version=2
```

## Фичи
### Опции конструктора `config.New(opts ...config.Option)`

#### `config.WithEnvPrefix(prefix string)`

Устанавливает префикс для переменных окружения (напр. для `prefix = "myapp"` все переменные окружения
должны начинаться с MYAPP_)

#### `config.WithFlagSet(set *pflag.FlagSet)`

Использует кастомный FlagSet

#### `config.WithViper(v *viper.Viper)`

Использует кастомный Viper

#### `config.WithDecodeFunctions(decoders ...mapstructure.DecodeHookFunc)`

Добавляет функцию декодирования строки в свой тип

### Тэги структуры

#### Тэг `env:"[name],squash"`

Тэг мапинга из файлов или переменных окружения. Используется в случаях:
- если название поля в структуре отличается от названия поля во входящих данных (устанавливается новое имя `env:"[name]"`)
- если необходимо "объеденить" поля вложенной структуры с родительской при мапинге входящих данных `env:",squash"`

#### Тэг `flag:"[name],[short],[description]"`

Задаёт мапинг поля структуры и флага команды запуска.
Если `name` пустое, то берётся название поля структуры в kebab-case.

Может использована в вариантах:
- `flag:""`
- `flag:"[name]"`
- `flag:"[name],[description]"`
- `flag:"[name],[short],[description]"`

#### Тэг `validate:"...""`

См. [github.com/go-playground/validator](https://github.com/go-playground/validator)


