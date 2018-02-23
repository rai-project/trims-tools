package micro

import (
	"github.com/dc0d/dirwatch"
	"github.com/fsnotify/fsnotify"
)

type notifyFunc func(fsnotify.Event)
type watcher struct {
	*dirwatch.Watch
	roots  []string
	notify notifyFunc
}

func NewWatcher(notify notifyFunc, rootDirectories ...string) *watcher {
	watch := dirwatch.New(notify, rootDirectories...)
	return &watcher{
		Watch:  watch,
		roots:  rootDirectories,
		notify: notify,
	}
}
