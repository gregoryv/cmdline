package cmdline_test

import (
	"path"
	"testing"

	. "github.com/gregoryv/web"
)

func Test_generate_readme(t *testing.T) {
	project := "gregoryv/cmdline"

	body := Body(
		travisBadge(project),
		codecovBadge(project),
		codeclimateBadge(project, "3dbee57c607ffec60702"),
		Br(),
		P(
			"Package ", godoc(project),
			" provides a means to parse command line arguments.",
		),
		P("This package fixes opinionated issues with using the flag package."),
		Ol(
			Li("Don't hog the name flag, which is a boolean option"),
			Li("Use appropriate names for arguments, options and flags"),
			Li("Self documenting arguments and options are preferred"),
			Li("Multiname options, e.g. -n, --dry-run map to same flag"),
			Li("Skip pointer variations"),
			Li("Include required arguments"),
		),
	)
	page := NewFile("README.md",
		Html(body),
	)
	page.SaveTo(".")
}

func godoc(project string) *Element {
	var (
		base = "https://pkg.go.dev/pkg/"
		href = base + path.Join("github.com", project)
		text = path.Base(project)
	)
	return A(Href(href), text)

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
		href = Href(base + project + "maintainability")

		api   = "https://api.codeclimate.com/v1"
		badge = path.Join("/badges/", hash, "/maintainability")
		src   = Src(api + badge)
		alt   = Alt("Maintainability")
	)
	return A(href, Img(src, alt))
}
