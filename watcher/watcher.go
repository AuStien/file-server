package watcher

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

type Watcher struct {
	*fsnotify.Watcher
}

// NewWatcher creates a new file watcher.
//
// Remember to close it!
func NewWatcher(rootPath string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			log.WithField("folder", path).Info("Adding folder to watcher")
			return watcher.Add(path)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Watcher{watcher}, nil
}

func (w *Watcher) Watch(folder chan string, errChan chan error) {
	for {
		select {
		case event := <-w.Events:

			switch event.Op {
			// NOTE: when a file is deleted (at least) on Linux a RENAME is emitted instead of DELETE
			case fsnotify.Remove, fsnotify.Rename:
				log.WithField("event", event).Debug("Handling event")
				folder <- event.Name

			case fsnotify.Create:
				log.WithField("event", event).Debug("Handling event")
				folder <- event.Name
				fi, err := os.Stat(event.Name)
				if err != nil {
					errChan <- err
				}
				if fi.IsDir() {
					if err := w.Add(event.Name); err != nil {
						if err == fsnotify.ErrNonExistentWatch {
							log.WithField("path", event.Name).Warn("Tried to add non existent path to watcher")
						} else {
							errChan <- err
						}
					}
				}

			default:
				log.WithField("event", event).Debug("Ignoring event")
			}

		case err := <-w.Errors:
			errChan <- err
		}
	}
}
