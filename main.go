package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/ladydascalie/mdg/config"
	"github.com/ladydascalie/mdg/file/manipulate"
	"github.com/ladydascalie/mdg/workers"
)

func init() {
	// Preload the CSS so it is immediately available
	loadCSS()
}

func main() {
	flag.StringVar(&config.DirPath, "d", ".", "mdg -d path/to/folder")
	flag.BoolVar(&config.SkipMenu, "m", false, "mdg -m | Use to skip generating the menu")
	flag.Parse()

	// Get the list of markdown files in the current directory
	fileList := manipulate.FindFilesOfType(config.FileExtensions)

	if len(fileList) == 0 {
		// Abort if no files found
		fmt.Println("No markdown files found in folder.\nAborting...")
		return
	}

	var wg sync.WaitGroup
	for _, file := range fileList {
		wg.Add(1)
		go workers.Process(file, fileList, &wg)
	}
	wg.Wait()
	close(workers.Semaphore)
}

func loadCSS() {
	var err error
	config.CSS, err = Asset("github-markdown.html")
	if err != nil {
		panic(err)
	}
}
