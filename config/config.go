package config

import "regexp"


// Flags
var FilePath string
var DirPath string
var SkipMenu bool

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
