//go:build generate

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
)

func main() {
	var buf bytes.Buffer

	data := make(map[string]any)
	data["Package"] = os.Args[1]
	data["Variants"] = []struct {
		Public, Private string
	}{
		{"String", "string"},
		{"Int", "int"},
		{"Uint8", "uint8"},
		{"Uint16", "uint16"},
		{"Float64", "float64"},
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal(err)
	}
	out := "generated.go"
	os.WriteFile(out, buf.Bytes(), 0644)
	fmt.Println("generated", out)
}

var tmpl = template.Must(template.New("").Parse(
	`// GENERATED CODE, DO NOT EDIT

package {{.Package}}

{{range .Variants}}
// {{.Public}}Var same as {{.Public}} but uses dst as default and destination.
func (opt *Option) {{.Public}}Var(dst *{{.Private}}) {{.Private}} {
	*dst = opt.{{.Public}}(*dst)
	return *dst
}
{{end}}
`))
