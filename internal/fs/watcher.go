package fs

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type FSWatcher struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Update   func()
}

func (f *FSWatcher) StartFSWatcher() {
	// Directories to monitor
	directories := []string{
		"./ui/html",
		"./ui/html/pages",
		"./ui/html/partials",
	}

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		f.ErrorLog.Fatal(err)
	}

	// Start listening for events.
	go func() {
		// defer func() {
		// 	if err := recover(); err != nil {
		// 		log.Print(fmt.Errorf("%s\n%s", err, debug.Stack()))
		// 	}
		// }()
		for {

			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					f.InfoLog.Println("modified file:", event.Name)
					f.Update()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				f.ErrorLog.Println("error:", err)
			}
		}
	}()

	// Add a path.
	for _, dir := range directories {
		err = watcher.Add(dir)
		f.InfoLog.Println("Monitoring dir: ", dir)
		if err != nil {
			f.ErrorLog.Fatal(err)
		}
	}
}
