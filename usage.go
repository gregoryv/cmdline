package cmdline

import (
	"fmt"
	"io"
)

// type Usage is for documenting the configured parser
type Usage struct {
	*Parser
}

// WriteUsageTo writes names, defaults and documentation to the given
// writer with the first line being
//
//   Usage: COMMAND [OPTIONS] ARGUMENTS...
func (me *Usage) WriteTo(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [OPTIONS]", me.args[0])
	for _, arg := range me.arguments {
		if arg.required {
			fmt.Fprintf(w, " %s", arg.name)
			continue
		}
		fmt.Fprintf(w, " [%s]", arg.name)
	}
	fmt.Fprint(w, "\n\n")
	fmt.Fprintln(w, "Options")
	me.WriteOptionsTo(w)

	indent := "    "
	for _, grp := range me.groups {
		fmt.Fprintln(w, grp.Title())
		first := grp.Items()[0]
		writeItem(w, first, me.args, indent, true)
		for _, item := range grp.Items()[1:] {
			writeItem(w, item, me.args, indent, false)
		}
	}
}

// WriteOptionsTo writes the Options section to the given writer.
func (me *Usage) WriteOptionsTo(w io.Writer) {
	me.writeOptionsTo(w, "")
}

func (me *Usage) writeOptionsTo(w io.Writer, indent string) {
	for _, opt := range me.options {
		writeOptionTo(w, opt, indent)
	}
	if len(me.options) > 0 {
		fmt.Fprintln(w)
	}
}

func writeItem(w io.Writer, me *Item, args []string, indent string, dflt bool) {
	if dflt {
		fmt.Fprintf(w, "%s%s (default)\n", indent, me.Name)
	} else {
		fmt.Fprintf(w, "%s%s\n", indent, me.Name)
	}
	extra := NewParser()
	extra.args = args
	me.Load(extra)
	extra.Usage().writeOptionsTo(w, indent)
}

func writeOptionTo(w io.Writer, opt *Option, indent string) {
	var def string
	if opt.defaultValue != "" {
		def = fmt.Sprintf(" : %v", opt.defaultValue)
	}
	if opt.quoteValue {
		def = fmt.Sprintf(" : %q", opt.defaultValue)
	}

	fmt.Fprintf(w, "%s    %s%s\n", indent, opt.names, def)
	if len(opt.doc) > 0 {
		for _, line := range opt.doc {
			fmt.Fprintf(w, "%s        %s\n", indent, line)
		}
		fmt.Fprintln(w)
	}
}
