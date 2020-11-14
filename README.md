[![Build Status](https://travis-ci.org/gregoryv/cmdline.svg?branch=master)](https://travis-ci.org/gregoryv/cmdline)
[![codecov](https://codecov.io/gh/gregoryv/cmdline/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/cmdline)
[![Maintainability](https://api.codeclimate.com/v1/badges/3dbee57c607ffec60702/maintainability)](https://codeclimate.com/github/gregoryv/cmdline/maintainability)


Package [cmdline](https://godoc.org/pkg/github.com/gregoryv/cmdline)
provides a way to parse command line arguments

This package fixes opinionated issues with using the flag package.

  1. Don't hog the name flag, which is a boolean option
  2. Use appropriate names for arguments, options and flags
  3. Optional documentation, self documenting options are preferred
  4. Simplify multiname options, e.g. -n, --dry-run map to same flag
  5. Skip pointer variations
  6. Include required arguments

Example:

    func main() {
        var (
            cli      = cmdline.New(args...)
            uid      = cli.Option("--uid", "Generated if not given").Int(0)
            password = cli.Option("-p, --password").String("")
            help     = cli.Flag("-h, --help")
    
            // parse and name non options
            username = cli.Required("USERNAME").String()
        )

        switch {
        case !cli.Ok():
            fmt.Println(cli.Error())
            fmt.Println("Try --help for more information")

        case help:
            cli.WriteUsageTo(os.Stdout)

        default:
            // ...
    }

Usage is written as

    Usage: adduser [OPTIONS] USERNAME
    
    Options
        --uid : 0
             Generated if not given
    
        -p, --password : ""
        -h, --help : false

