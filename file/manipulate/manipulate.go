package manipulate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ladydascalie/mdg/config"
	md "github.com/shurcooL/github_flavored_markdown"
)

// EnsureCharset appends the HTML meta tag
func EnsureCharset(file []byte) []byte {
	return append(config.Charset, file...)
}

// EnsureOcticons is not yet in use
func EnsureOcticons(file []byte) []byte {
	return append(config.Octicons, file...)
}

// GenerateMenu creates a simple list menu
// when less than 40 items are in the folder
func GenerateMenu(fileList []string) []byte {
	if len(fileList) > 40 {
		return []byte("")
	}

	menu := "#### Menu\n"
	for _, file := range fileList {
		suffix, err := extractSuffix(file)
		if err != nil {
			panic(err)
		}
		f := strings.TrimSuffix(file, suffix)
		menu += fmt.Sprintf("- [%s](%s.html)\n", f, f)
	}
	menu += "\n---\n\n"

	return []byte(menu)
}

func extractSuffix(name string) (string, error) {
	for _, ext := range config.FileExtensions {
		if strings.Contains(name, ext) {
			return ext, nil
		}
	}
	return "", fmt.Errorf("%s", "Could not derive suffix")
}

// NewFileName ...
func NewFileName(name string) string {
	suffix, err := extractSuffix(name)
	if err != nil {
		panic(err)
	}

	name = strings.TrimSuffix(name, suffix)
	return name + ".html"
}

// Tokenizer not in use at the moment
func replaceTokens(stream []byte) []byte {
	buffer := bytes.NewBuffer(stream)
	temp := buffer.Bytes()

	for _, match := range config.LinksRegExp.FindAll(temp, -1) {
		stripped := bytes.Trim(match, "{}")
		linked := []byte(string(stripped) + ".html")
		temp = bytes.Replace(temp, match, linked, -1)
	}

	return temp
}

// AppendCSS appends any CSS byte stream to the file
func AppendCSS(css, stream []byte) []byte {
	stream = append(css, stream...)
	return stream
}

// CompileMarkdown compiles the markdown to html
func CompileMarkdown(text []byte) []byte {
	return md.Markdown(text)
}

// FindFilesOfType loops through the provided file extensions
// and returns a slice containing the names of matching files
func FindFilesOfType(extensions []string) []string {
	files, err := ioutil.ReadDir(config.DirPath)
	if err != nil {
		panic(err)
	}

	var list []string
	for _, file := range files {
		for _, ext := range extensions {
			if strings.HasSuffix(file.Name(), ext) {
				list = append(list, file.Name())
			}
		}
	}

	return list
}
