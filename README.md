[![Build Status](https:/travis-ci.org/gregoryv/cmdline.svg?branch=master)](https:/travis-ci.org/gregoryv/cmdline)
[![Code coverage](https:/codecov.io/gh/gregoryv/cmdline/branch/master/graph/badge.svg)](https:/codecov.io/gh/gregoryv/cmdline)
[![Maintainability](https:/api.codeclimate.com/v1/badges/3dbee57c607ffec60702/maintainability)](https:/codeclimate.com/github/gregoryv/cmdline/maintainability)



Package[cmdline](https:/godoc.org/pkg/github.com/gregoryv/cmdline)provides a means to parse command line arguments.

This package fixes opinionated issues with using the flag package.

- Don't hog the name flag, which is a boolean option
- Use appropriate names for arguments, options and flags
- Self documenting arguments and options are preferred
- Multiname options, e.g. -n, --dry-run map to same flag
- Skip pointer variations
- Include required arguments


## Example

    func run(w io.Writer, args ...string) {
        var (
            cli      = cmdline.New(args...)
            uid      = cli.Option("--uid", "Generated if not given").Int(0)
            password = cli.Option("-p, --password").String("")
            help     = cli.Flag("-h, --help")
    
            // parse and name non options
            username = cli.Required("USERNAME").String()
            note     = cli.Optional("NOTE").String()
        )
    
        switch {
        case help:
            cli.WriteUsageTo(w)
    
        case !cli.Ok():
            fmt.Fprintln(w, cli.Error())
            fmt.Fprintln(w, "Try --help for more information")
    
        default:
            fmt.Fprintln(w, uid, username, password, note)
        }
    }
    

Output

    Usage: adduser [OPTIONS] USERNAME [NOTE]
    
    Options
        --uid : 0
    	 Generated if not given
    
        -p, --password : ""
        -h, --help : false
    





