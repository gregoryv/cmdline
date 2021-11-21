# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [unreleased]

- Replace Parser.Optional and Parser.Required with Parser.Argument

## [0.8.1] 2021-11-13

- Fix BasicParser -h if given

## [0.8.0] 2021-11-13

- Replace Parser.SetArgs and SetExit with SetShell
- Add types Shell, ShellOS and ShellT

## [0.7.0] 2021-10-21

- Add example command
- Add optional examples section to usage
- Add optional preface to usage
- Remove method Parser.WriteUsageTo, use Parser.Usage().WriteTo
- Add type Usage

## [0.6.1] 2021-10-17

- Use os.Stdout to write parsing errors
- Fix BasicParser.Parse

## [0.6.0] 2021-10-17

- Add NewBasicParser to simplify parsers using -h, --help for usage
- Add method Parser.Parse with default error checking
- Remove func Parse
- NewParser defaults to using os.Args if none are given
- Option names can include $ENV_NAME to parse environment variabls
- Usage no longer writes " : false" for flag options, it's implied

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
