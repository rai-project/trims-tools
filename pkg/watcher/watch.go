package watcher

import (
	"github.com/dc0d/dirwatch"
	"github.com/fsnotify/fsnotify"
)

type notifyFunc func(fsnotify.Event)
type Watcher struct {
	*dirwatch.Watch
	roots  []string
	notify notifyFunc
}

func New(notify notifyFunc, rootDirectories ...string) *Watcher {
	watch := dirwatch.New(notify, rootDirectories...)
	return &Watcher{
		Watch:  watch,
		roots:  rootDirectories,
		notify: notify,
	}
}
