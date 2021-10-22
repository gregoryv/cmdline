package cmdline_test

import (
	"path"
	"testing"

	. "github.com/gregoryv/web"
	"github.com/gregoryv/web/files"
)

func Test_generate_readme(t *testing.T) {
	project := "gregoryv/cmdline"

	body := Body(
		"<!-- Generated by doc_test.go, DO NOT EDIT! -->\n",
		githubBadge(project+"/actions/workflows/go.yml"),
		codecovBadge(project),
		codeclimateBadge(project, "d1e0ac639370f6fc982e"),
		Br(),
		P(
			"Package ", godoc(project),
			" provides a parser for command line arguments.",
		),
		P("This package is different from the builtin package flag."),
		Ol(
			Li("Options are used when parsing arguments"),
			Li("Flag is a boolean option"),
			Li("Self documenting options are preferred"),
			Li("Multiname options, e.g. -n, --dry-run map to same option"),
			Li("Easy way to default to environment variables"),
			Li("There are no pointer variations"),
			Li("Parsing non option arguments"),
			Li("Usage output with optional examples and preface"),
		),

		H2("Example"),
		Pre(
			files.MustLoadFunc("example_test.go", "ExampleNewBasicParser"),
		),
	)
	page := NewPage(
		Html(body),
	)
	page.SaveAs("README.md")
}

func godoc(project string) *Element {
	var (
		base = "https://pkg.go.dev/pkg/"
		href = base + path.Join("github.com", project)
		text = path.Base(project)
	)
	return A(Href(href), text)

}

func githubBadge(workflow string) *Element {
	var (
		base = "https://github.com/"
		href = Href(base + workflow)
		src  = Src(base + workflow + "/badge.svg")
		alt  = Alt("Go")
	)
	return A(href, Img(src, alt))
}

func travisBadge(project string) *Element {
	var (
		base = "https://travis-ci.org/"
		href = Href(base + project)
		src  = Src(base + project + ".svg?branch=master")
		alt  = Alt("Build Status")
	)
	return A(href, Img(src, alt))
}

func codecovBadge(project string) *Element {
	var (
		base = "https://codecov.io/gh/"
		href = Href(base + project)
		src  = Src(base + project + "/branch/master/graph/badge.svg")
		alt  = Alt("Code coverage")
	)
	return A(href, Img(src, alt))
}

func codeclimateBadge(project, hash string) *Element {
	var (
		base = "https://codeclimate.com/github/"
		href = Href(base + project + "/maintainability")
		// image
		api   = "https://api.codeclimate.com/v1"
		badge = path.Join("/badges/", hash, "/maintainability")
		src   = Src(api + badge)
		alt   = Alt("Maintainability")
	)
	return A(href, Img(src, alt))
}
