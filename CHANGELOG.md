# Change Log


All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]


## [2.0.8.1] - 2021-01-09


## [2.0.8] - 2021-01-09
- Added root `cobra.Command` in container
- Added `server.CommandMode()` option for run commands
- Fixed flag resolving with `cobra.Command`
- Removed package `inject`

## [2.0.7] - 2021-01-09
- MultiLog:
    - Added `LocalMode` (colored human console) and `DebugMode` (trace log level) to `config.Application`
    - Application name in log fields
    - By default, console log enabled (added `SilentMode` in `multilog.Config` to turn it off)
- Runner:
    - All builtin providers are public
    - Added startup options:
        - `server.NoDefaultProviders()` - turn off all default providers and configurations
        - `server.LocalDebug()` - turn on colored console logging and trace level
- Other:
    - `mysql.Provide` now returns `driver.Connector` (and added `mysql.ProvideConnector` for manual registering)
    - `sql.Provide` now accepting `driver.Connector`
- Deprecated:
    - `server.NoWaitOption()` -> `server.NoWait()`
    - `server.EnvPathOption(path)` -> `server.EnvPath(path)`
    - `server.ConfigOption(cfg)` -> `server.WithConfig(cfg)`
    - Field `MySQL` in `sql.Args` structure

## [2.0.6] - 2020-11-12

## [2.0.5] - 2020-11-12

## [2.0.4] - 2020-11-12

## [2.0.3] - 2020-11-12

## [2.0.2] - 2020-11-11

## [2.0.1] - 2020-11-06

## [2.0.0] - 2020-11-06

## [1.0.2] - 2020-10-29

## [1.0.1] - 2020-10-29

## 1.0.0 - 2020-09-28

[Unreleased]: https://github.com/vseinstrumentiru/lego/compare/v2.0.8.1...HEAD
[2.0.8.1]: https://github.com/vseinstrumentiru/lego/compare/v2.0.8...v2.0.8.1
[2.0.8]: https://github.com/vseinstrumentiru/lego/compare/v2.0.7...v2.0.8
[2.0.7]: https://github.com/vseinstrumentiru/lego/compare/v2.0.6...v2.0.7
[2.0.6]: https://github.com/vseinstrumentiru/lego/compare/v2.0.5...v2.0.6
[2.0.5]: https://github.com/vseinstrumentiru/lego/compare/v2.0.4...v2.0.5
[2.0.4]: https://github.com/vseinstrumentiru/lego/compare/v2.0.3...v2.0.4
[2.0.3]: https://github.com/vseinstrumentiru/lego/compare/v2.0.2...v2.0.3
[2.0.2]: https://github.com/vseinstrumentiru/lego/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/vseinstrumentiru/lego/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/vseinstrumentiru/lego/compare/v1.0.2...v2.0.0
[1.0.2]: https://github.com/vseinstrumentiru/lego/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/vseinstrumentiru/lego/compare/v1.0.0...v1.0.1
