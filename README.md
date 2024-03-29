<!-- Generated by doc_test.go, DO NOT EDIT -->

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
    	os.Setenv("VERBOSE", "yes")
    	var (
    		cli		= NewBasicParser()
    		uid		= cli.Option("--uid", "Generated if not given").Int(0)
    		password	= cli.Option("-p, --password, $PASSWORD").String("")
    		verbose		= cli.Flag("-V, --verbose, $VERBOSE")
    		role		= cli.Option("-r, --role").Enum("guest",
    			"guest", "admin", "nobody",
    		)
    		url	= cli.Option("--test-host").Url("tcp://example.com:123")
    		dur	= cli.Option("--pause").Duration("200ms")
    
    		// parse and name non options
    		username	= cli.NamedArg("USERNAME").String("")
    		note		= cli.NamedArg("NOTE").String("")
    	)
    	cli.Parse()
    
    	if !verbose {
    		log.SetOutput(ioutil.Discard)
    	}
    	fmt.Fprintln(os.Stdout, uid, username, password, note, role, url, dur)
    }

