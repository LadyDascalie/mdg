package config

import "regexp"

// DirPath is the CLI flag storing the location of the folder you want to run mdg against
var DirPath string

// SkipMenu is the CLI flag indicating wether or not mdg should skip generating a menu
var SkipMenu bool

// CSS is used to store the preloaded github-markdown.html asset
var CSS []byte

// LinksRegExp matches link tokens in the following format
// [link]({{content}})
var LinksRegExp = regexp.MustCompile(`(?:\{\{)(.{1,})(?:\}\})`)

// Charset is used to embed the meta tag in the final file
var Charset = []byte("<meta charset=\"UTF-8\">")

// Octicons path in cdn.js
var Octicons = []byte(`<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/octicons/4.4.0/font/octicons.css" integrity="sha256-4y5taf5SyiLyIqR9jL3AoJU216Rb8uBEDuUjBHO9tsQ=" crossorigin="anonymous" />`)

// FileExtensions is the list of possible prefixes to look for
var FileExtensions = []string{".md", ".markdown"}
