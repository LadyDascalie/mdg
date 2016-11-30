package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"bytes"
	"sync"

	md "github.com/shurcooL/github_flavored_markdown"
)

var filePath string
var dirPath string
var skipMenu bool
var linksRegExp = regexp.MustCompile(`(?:\{\{)(.{1,})(?:\}\})`)

// default utf8 charset
var charset = []byte("<meta charset=\"UTF-8\">")

// default path for the octicons icon set in cdn.js
var octicons = []byte(`<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/octicons/4.4.0/font/octicons.css" integrity="sha256-4y5taf5SyiLyIqR9jL3AoJU216Rb8uBEDuUjBHO9tsQ=" crossorigin="anonymous" />`)

// This ensures we never run more than 12 Goroutines at the same time
// this prevents opening too many file descriptors without clearing them
var semaphore = make(chan struct{}, 12)

func main() {
	flag.StringVar(&filePath, "f", "", "mdg -f path/to/file")
	flag.StringVar(&dirPath, "d", ".", "mdg -d path/to/folder")
	flag.BoolVar(&skipMenu, "m", false, "mdg -m | Use to skip generating the menu")
	flag.Parse()

	// Get the list of markdown files in the current directory
	fileList := seekPrefixedFiles(".md")

	var wg sync.WaitGroup
	for _, file := range fileList {
		wg.Add(1)
		go process(file, fileList, &wg)
	}
	wg.Wait()
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

	// Get a file descriptor
	f, err := os.Open(file)
	if err != nil {
		log.Println("Cannot open file", err)
		return
	}

	// Read the contents of the open file
	openedFile, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("Cannot read content of file", err)
		return
	}


	var fileContents []byte

	// Skip menu if flag is set
	if skipMenu {
		fileContents = openedFile
	} else {
		fileContents = generateMenu(fileList)
		for _, v := range openedFile {
			fileContents = append(fileContents, v)
		}
	}

	// Compile and append styles
	fileContents = compileMarkdown(fileContents)
	fileContents = appendCSS(fileContents)

	// Ensure UTF-8 Encoding is properly appended to the document
	fileContents = ensureCharset(fileContents)

	// basically I need to write the file like that once it's compiled lol
	if _, err = os.Stat("html"); os.IsNotExist(err) {
		os.Mkdir("html", 0777)
	}

	// Name and save the file
	newFileName := newFileName(file)
	ioutil.WriteFile(newFileName, fileContents, 0777)

	// Move the new file into th html sub-folder
	// Overwrites are allowed
	err = os.Rename(newFileName, "html/"+ newFileName)
	if err != nil {
		log.Println(err)
		return
	}
}

func ensureCharset(file []byte) []byte {
	return append(charset, file...)
}

// Not in use currently
func ensureOcticons(file []byte) []byte {
	return append(octicons, file...)
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

// Tokenizer not in use at the moment
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
