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
func (u *Usage) Preface(lines ...string) {
	for _, line := range lines {
		u.preface.WriteString(line)
		u.preface.WriteString("\n")
	}
}

// Example adds an examples section. The examples are placed last
// after options and named arguments. Examples are plain text and not
// evaluated in any way.
func (u *Usage) Example(lines ...string) {
	for _, line := range lines {
		if u.examples.Len() > 0 {
			u.examples.WriteString("\n")
		}
		u.examples.WriteString(indent)
		u.examples.WriteString(line)
	}
}

// WriteUsageTo writes names, defaults and documentation to the given
// writer with the first line being
//
//	Usage: COMMAND [OPTIONS] ARGUMENTS...
func (u *Usage) WriteTo(w io.Writer) (int64, error) {
	p, err := nexus.NewPrinter(w)
	p.Printf("Usage: %s [OPTIONS]", u.args[0])
	// Named arguments
	for _, arg := range u.arguments {
		if arg.required {
			p.Printf(" %s", arg.name)
			continue
		}
		p.Printf(" [%s]", arg.name)
	}
	// Preface
	p.Print("\n\n")
	u.writePreface(p)
	// Options
	p.Println("Options")
	u.WriteOptionsTo(p)
	if len(u.options) > 0 {
		fmt.Fprintln(w)
	}
	u.writeGroups(p)
	u.writeExamples(p)

	p.Println()
	return p.Written, *err
}

const indent = "    "

func (u *Usage) writeGroups(p *nexus.Printer) {
	if len(u.groups) == 0 {
		return
	}
	for _, grp := range u.groups {
		p.Println(grp.Title())
		first := grp.Items()[0]
		writeItem(p, first, u.args, indent, true)
		for _, item := range grp.Items()[1:] {
			writeItem(p, item, u.args, indent, false)
		}
	}
}

// WriteOptionsTo writes the Options section to the given writer.
func (u *Usage) WriteOptionsTo(w io.Writer) {
	u.writeOptionsTo(w, "")
}

func (u *Usage) writePreface(p *nexus.Printer) {
	if u.preface.Len() == 0 {
		return
	}
	p.Println(u.preface.String())
}

func (u *Usage) writeExamples(p *nexus.Printer) {
	if u.examples.Len() == 0 {
		return
	}
	if len(u.groups) > 0 {
		p.Println()
	}
	p.Println("Examples")
	p.Print(u.examples.String())
}

func (u *Usage) writeOptionsTo(w io.Writer, indent string) {
	for _, opt := range u.options {
		writeOptionTo(w, opt, indent)
	}
}

func writeItem(w io.Writer, m *Item, args []string, indent string, dflt bool) {
	if dflt {
		fmt.Fprintf(w, "%s%s (default)\n", indent, m.Name)
	} else {
		fmt.Fprintf(w, "%s%s\n", indent, m.Name)
	}
	extra := NewParser()
	extra.args = args
	m.Load(extra)
	extra.Usage().writeOptionsTo(w, indent)
}

func writeOptionTo(w io.Writer, opt *Option, indent string) {
	var def string
	val := opt.defaultValue
	if opt.hidden {
		val = "********"
	}
	if val != "" {
		def = fmt.Sprintf(" : %v", val)
	}
	if opt.quoteValue {
		def = fmt.Sprintf(" : %q", val)
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
