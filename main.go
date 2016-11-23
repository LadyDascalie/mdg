package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	md "github.com/shurcooL/github_flavored_markdown"
)

var filePath string
var dirPath string

func main() {
	flag.StringVar(&filePath, "f", "", "mdg -f path/to/file")
	flag.StringVar(&dirPath, "d", ".", "mdg -d path/to/folder")
	flag.Parse()

	// list is a []string of markdown files
	fileList := seekPrefixedFiles(".md")

	for _, file := range fileList {
		f, err := os.Open(file)
		if err != nil {
			log.Println("Cannot open file", err)
			return
		}

		fileContent, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println("Cannot read content of file", err)
			return
		}

		fileMenu := generateMenu(fileList)
		for _, v := range fileContent {
			fileMenu = append(fileMenu, v)
		}
		fileMenu = replaceTokens(fileMenu)
		fileMenu = compileMarkdown(fileMenu)
		fileMenu = appendCSS(fileMenu)

		// basically I need to write the file like that once it's compiled lol
		if _, err := os.Stat("html"); os.IsNotExist(err) {
			os.Mkdir("html", 0777)
		}

		newFname := newFileName(file)
		ioutil.WriteFile(newFname, fileMenu, 0777)

		err = os.Rename(newFname, "html/" + newFname)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func generateMenu(fileList []string) []byte {
	menu := "#### Menu\n"
	for _, file := range fileList {
		f := strings.TrimSuffix(file, ".md")
		menu += fmt.Sprintf("[%s]({{%s}})\n", f, f)
	}
	menu += "\n---\n\n"

	return []byte(menu)
}

func newFileName(name string) string {
	name = strings.TrimSuffix(name, ".md")
	return name + ".html"
}

func appendCSS(stream []byte) []byte {
	temp := string(stream)
	css, err := Asset("github-markdown.css")
	if err != nil {
		panic(err)
	}
	temp += fmt.Sprintf("<style>%s</style>", string(css))

	return []byte(temp)
}

func replaceTokens(stream []byte) []byte {
	temp := string(stream)

	var re = regexp.MustCompile(`(?:\{\{)(.{1,})(?:\}\})`)

	for _, match := range re.FindAllString(temp, -1) {
		stripped := strings.Trim(match, "{}")
		linked := fmt.Sprint(stripped + ".html")
		temp = strings.Replace(temp, match, linked, -1)
	}

	return []byte(temp)
}

func compileMarkdown(text []byte) []byte {
	return md.Markdown(text)
}

func seekPrefixedFiles(prefix string) []string {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	var list []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), prefix) {
			list = append(list, file.Name())
		}
	}

	return list
}
