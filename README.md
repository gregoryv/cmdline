[![Build Status](https://travis-ci.org/gregoryv/cmdline.svg?branch=master)](https://travis-ci.org/gregoryv/cmdline)
[![Code coverage](https://codecov.io/gh/gregoryv/cmdline/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/cmdline)
[![Maintainability](https://api.codeclimate.com/v1/badges/3dbee57c607ffec60702/maintainability)](https://codeclimate.com/github/gregoryv/cmdlinemaintainability)


Package [cmdline](https://pkg.go.dev/pkg/github.com/gregoryv/cmdline) provides a means to parse command line arguments.
This package fixes opinionated issues with using the flag package.
- Don't hog the name flag, which is a boolean option
- Use appropriate names for arguments, options and flags
- Self documenting arguments and options are preferred
- Multiname options, e.g. -n, --dry-run map to same flag
- Skip pointer variations
- Include required arguments

