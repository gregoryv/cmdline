[![Build Status](https://travis-ci.org/gregoryv/cmdline.svg?branch=master)](https://travis-ci.org/gregoryv/cmdline)
[![codecov](https://codecov.io/gh/gregoryv/cmdline/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/cmdline)
[![Maintainability](https://api.codeclimate.com/v1/badges/3dbee57c607ffec60702/maintainability)](https://codeclimate.com/github/gregoryv/cmdline/maintainability)


[cmdline](https://godoc.org/pkg/github.com/gregoryv/cmdline) -
package cmdline provides a way to parse command line arguments

This package fixes a few problems that come with using the flag package.

  1. Don't hog the name flag, which is a boolean option
  2. Use appropriate names for arguments, options and flags
  3. Optional documentation, self documenting options are preferred
  4. Simplify multiname options, e.g. -n, --dry-run map to same flag
  5. Skip pointer variations
  6. Include required options TODO

Example:

    func main() {
        cli := cmdline.New(os.Args)
        uid, opt := cli.Option("--uid").IntOpt(0)
        opt.Doc(
                "user id to set on the new account",
                "If not given, one is generated",
        )
        username := cli.Option("-u, --username").String("john")
        password := cli.Option("-p, --password").String("")
        dryrun := cli.Flag("-n, --dry-run")
        cli.CheckOptions()

        // ...
    }
