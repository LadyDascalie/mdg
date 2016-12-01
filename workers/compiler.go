package workers

import (
	"sync"
	"os"
	"io/ioutil"
	"log"
	"github.com/ladydascalie/mdg/config"
	"github.com/ladydascalie/mdg/file/manipulate"
)

// This ensures we never run more than 12 Goroutines at the same time
// this prevents opening too many file descriptors without clearing them
var Semaphore = make(chan struct{}, 12)

func Process(file string, fileList []string, wg *sync.WaitGroup) {
	// If buffer is full, sending will be blocked
	// And tasks will wait until space is cleared in the buffer
	Semaphore <- struct{}{}

	// Once this goroutine is done, decrement the buffer by one
	defer func() {
		<-Semaphore
	}()

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
	if config.SkipMenu {
		fileContents = openedFile
	} else {
		fileContents = manipulate.GenerateMenu(fileList)
		for _, v := range openedFile {
			fileContents = append(fileContents, v)
		}
	}

	// Compile and append styles
	fileContents = manipulate.CompileMarkdown(fileContents)
	fileContents = manipulate.AppendCSS(config.CSS, fileContents)

	// Ensure UTF-8 Encoding is properly appended to the document
	fileContents = manipulate.EnsureCharset(fileContents)

	// basically I need to write the file like that once it's compiled lol
	if _, err = os.Stat("html"); os.IsNotExist(err) {
		os.Mkdir("html", 0777)
	}

	// Name and save the file
	newFileName := manipulate.NewFileName(file)
	ioutil.WriteFile(newFileName, fileContents, 0777)

	// Move the new file into th html sub-folder
	// Overwrites are allowed
	err = os.Rename(newFileName, "html/" + newFileName)
	if err != nil {
		log.Println(err)
		return
	}
}
