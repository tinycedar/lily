package main

import (
	"fmt"
	"log"
	"os"

	"bufio"

	"github.com/fsnotify/fsnotify"
)

const (
	hostFile = "C:/Windows/System32/drivers/etc/hosts"
)

func main() {
	addFileWatcher()
}

func addFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event)
					// readFile()
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	if err = watcher.Add(hostFile); err != nil {
		log.Fatal(err)
	}
	<-done
}

func readFile() {
	file, err := os.Open(hostFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	fmt.Println("============================== Reading file =====================================")
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
