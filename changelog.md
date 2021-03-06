# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [0.5.2] 2021-03-21

- Update dependencies

## [0.5.1] 2021-03-03

- Add MIT license

## [0.5.0] 2021-02-11

- Add Float64 options
- Argument.String() takes default value
- Fix issue with missing string values, eg. "-str -flag" fails instead
  of using "-flag" as the value.
- Fail if unknown option is given even if groups are present

## [0.4.2] 2020-11-27

- Use type Group to define new items of various types
- Rename type Action to Item

## [0.4.1] 2020-11-26

- Added Uint8,-16 and -32 for convenience

## [0.4.0] 2020-11-26

- Renamed func New to NewParser
- Renamed type CommandLine to Parser

## [0.3.0] 2020-11-26

- Add type Group and Action for grouped options
- Added CommandLine.Required and Optional methods for defining non options

## [0.2.1] 2020-11-14
## [0.2.0] 2020-11-13

- Examples on showing usage and handling --help flag

## [0.1.0] 2020-11-13

- Initial version
