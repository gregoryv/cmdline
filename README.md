<!-- Generated by doc_test.go, DO NOT EDIT! -->

[![Go](https://github.com/gregoryv/cmdline/actions/workflows/go.yml/badge.svg)](https://github.com/gregoryv/cmdline/actions/workflows/go.yml)
[![Code coverage](https://codecov.io/gh/gregoryv/cmdline/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/cmdline)
[![Maintainability](https://api.codeclimate.com/v1/badges/d1e0ac639370f6fc982e/maintainability)](https://codeclimate.com/github/gregoryv/cmdline/maintainability)


Package [cmdline](https://pkg.go.dev/pkg/github.com/gregoryv/cmdline) provides a parser for command line arguments.
This package is different from the builtin package flag.
- Options are used when parsing arguments
- Flag is a boolean option
- Self documenting options are preferred
- Multiname options, e.g. -n, --dry-run map to same option
- Easy way to default to environment variables
- There are no pointer variations
- Parsing non option arguments
- Usage output with optional examples and preface

## Example

    func ExampleNewBasicParser() {
    	var (
    		cli		= cmdline.NewBasicParser()
    		uid		= cli.Option("--uid", "Generated if not given").Int(0)
    		password	= cli.Option("-p, --password, $PASSWORD").String("")
    
    		username	= cli.Argument("USERNAME").String("")
    		note		= cli.Argument("NOTE").String("")
    	)
    	cli.Parse()
    
    	fmt.Fprintln(os.Stdout, uid, username, password, note)
    }

