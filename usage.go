package cmdline

import (
	"fmt"
	"io"
	"strings"

	"github.com/gregoryv/nexus"
)

// type Usage is for documenting the configured parser
type Usage struct {
	*Parser

	preface  strings.Builder
	examples strings.Builder
}

// Preface adds lines just before the options section
func (me *Usage) Preface(lines ...string) {
	for _, line := range lines {
		me.preface.WriteString(line)
		me.preface.WriteString("\n")
	}
}

// Example adds an examples section. The examples are placed last
// after options and named arguments. Examples are plain text and not
// evaluated in any way.
func (me *Usage) Example(lines ...string) {
	for _, line := range lines {
		me.examples.WriteString(indent)
		me.examples.WriteString(line)
		me.examples.WriteString("\n")
	}
	me.examples.WriteString("\n")
}

// WriteUsageTo writes names, defaults and documentation to the given
// writer with the first line being
//
//   Usage: COMMAND [OPTIONS] ARGUMENTS...
func (me *Usage) WriteTo(w io.Writer) (int64, error) {
	p, err := nexus.NewPrinter(w)
	p.Printf("Usage: %s [OPTIONS]", me.args[0])
	// Named arguments
	for _, arg := range me.arguments {
		if arg.required {
			p.Printf(" %s", arg.name)
			continue
		}
		p.Printf(" [%s]", arg.name)
	}
	// Preface
	p.Print("\n\n")
	me.writePreface(p)
	// Options
	p.Println("Options")
	me.WriteOptionsTo(p)

	for _, grp := range me.groups {
		p.Println(grp.Title())
		first := grp.Items()[0]
		writeItem(p, first, me.args, indent, true)
		for _, item := range grp.Items()[1:] {
			writeItem(p, item, me.args, indent, false)
		}
	}
	me.writeExamples(p)
	return p.Written, *err
}

const indent = "    "

// WriteOptionsTo writes the Options section to the given writer.
func (me *Usage) WriteOptionsTo(w io.Writer) {
	me.writeOptionsTo(w, "")
}

func (me *Usage) writePreface(p *nexus.Printer) {
	if me.preface.Len() == 0 {
		return
	}
	p.Println(me.preface.String())
}

func (me *Usage) writeExamples(p *nexus.Printer) {
	if me.examples.Len() == 0 {
		return
	}
	if len(me.groups) == 0 {
		p.Println()
	}
	p.Println("Examples")
	p.Print(me.examples.String())
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
	var enum string
	if len(opt.enumerated) > 0 {
		enum = fmt.Sprintf(" %v", opt.enumerated)
	}

	fmt.Fprintf(w, "%s    %s%s%v\n", indent, opt.names, def, enum)
	writeDocTo(w, opt, indent)
}

func writeDocTo(w io.Writer, opt *Option, indent string) {
	if len(opt.doc) > 0 {
		for _, line := range opt.doc {
			fmt.Fprintf(w, "%s        %s\n", indent, line)
		}
		fmt.Fprintln(w)
	}
}
