package files

import (
	"log"
	"time"
	"github.com/fsnotify/fsnotify"
)

type Set map[string]bool

/**
	Execute a function when files are created
*/
func WatchNewFiles(fn func(file string)) *fsnotify.Watcher{
	fileSet := make(Set)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	
	/**
		The file watcher generates multiple values when a file is created, 
		so we need to debounce the events by savng them into a set
	*/
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fileSet[event.Name] = true
					
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	
	/**
		Use a second loop to get the debounced values provided by the file watcher
	*/
	go func(){
		for range time.Tick(3000*time.Millisecond) {
			for k, _ := range fileSet{
				delete(fileSet, k)
				fn(k)
			}
		}
	}()
	
	return watcher
}