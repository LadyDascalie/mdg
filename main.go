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
	"sync"
	"bytes"
)

var filePath string
var dirPath string
var shouldMenu string
var linksRegExp = regexp.MustCompile(`(?:\{\{)(.{1,})(?:\}\})`)
var charset = []byte("<meta charset=\"UTF-8\">")

// This ensures we never run more than 12 Goroutines at the same time
// this prevents opening too many file descriptors without clearing them
var semaphore = make(chan struct{}, 12)
func main() {
	flag.StringVar(&filePath, "f", "", "mdg -f path/to/file")
	flag.StringVar(&dirPath, "d", ".", "mdg -d path/to/folder")
	flag.StringVar(&shouldMenu, "m", "true", "mdg -m false")
	flag.Parse()

	// list is a []string of markdown files
	fileList := seekPrefixedFiles(".md")

	var wg sync.WaitGroup
	for _, file := range fileList {
		wg.Add(1)
		go process(file, fileList, &wg)
	}
	wg.Wait();
	close(semaphore)
}

func process(file string, fileList []string, wg *sync.WaitGroup) {
	// If buffer is full, sending will be blocked
	// And tasks will wait until space is cleared in the buffer
	semaphore <- struct{}{}

	// Once this goroutine is done, decrement the buffer by one
	defer func() { <-semaphore }()

	// defer until after the semaphore is read from
	defer wg.Done()
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

	var fileMenu []byte
	if shouldMenu == "true" {
		fileMenu = generateMenu(fileList)
		for _, v := range fileContent {
			fileMenu = append(fileMenu, v)
		}
	} else {
		fileMenu = fileContent
	}

	fileMenu = compileMarkdown(fileMenu)
	fileMenu = appendCSS(fileMenu)

	// Ensure UTF-8 Encoding is properly appended to the document
	fileMenu = ensureCharset(fileMenu)

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
	wg.Done()
}

func ensureCharset(file []byte) []byte {
	return append(charset, file...)
}

func generateMenu(fileList []string) []byte {
	if len(fileList) > 40 {
		return []byte("")
	}

	menu := "#### Menu\n"
	for _, file := range fileList {
		f := strings.TrimSuffix(file, ".md")
		menu += fmt.Sprintf("- [%s](%s.html)\n", f, f)
	}
	menu += "\n---\n\n"

	return []byte(menu)
}

func newFileName(name string) string {
	name = strings.TrimSuffix(name, ".md")
	return name + ".html"
}

func appendCSS(stream []byte) []byte {
	css, err := Asset("github-markdown.html")
	if err != nil {
		panic(err)
	}

	stream = append(css, stream...)

	return stream
}

func replaceTokens(stream []byte) []byte {
	buffer := bytes.NewBuffer(stream)
	temp := buffer.Bytes()

	for _, match := range linksRegExp.FindAll(temp, -1) {
		stripped := bytes.Trim(match, "{}")
		linked := []byte(string(stripped) + ".html")
		temp = bytes.Replace(temp, match, linked, -1)
	}

	return temp
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
