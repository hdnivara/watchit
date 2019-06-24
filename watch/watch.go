package watch

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

// Op represents the file-operations to watch for. E.g., create, write,
// remove, rename, etc.
type Op uint32

const (
	// Create represents the event when a file is created in watched
	// directory.
	Create Op = iota

	// Write represents the event when a watched file is written to.
	Write

	// Remove represents the event when a watched file is deleted.
	Remove

	// Rename represents the event when a watched file is renamed.
	Rename

	// Chmod represents the event when a watched file's premissions are
	// changed.
	Chmod

	// Tbi represents the events that are yet to be supported.
	Tbi
)

var opStr = map[Op]string{
	Create: "CREATE",
	Write:  "WRITE",
	Remove: "REMOVE",
	Rename: "RENAME",
	Chmod:  "CHMOD",
	Tbi:    "TBI",
}

// Handler is a pointer to handler function that users of 'watch'
// package can register. It will be called on the event of any
// operation, 'op', on the watched files, 'fileName'.
type Handler func(Op, string)

// Watch represents watcher configuration: dirctories, commands to
// execute, file types, etc.
type Watch struct {
	// Directories to watch.
	Dirs []string

	// Commands to exectue on watched files.
	Cmds []string

	// Regex to apply to filter watched file names.
	Regex string

	// Watch directories recursively?
	Recursive bool

	// HandlerCB function to call on file-change events.
	HandlerCB Handler

	// Pointer to underlying watcher,
	watcher *watcher.Watcher
}

// Name returns the human-redable name for Op.
func (op Op) Name() string {
	if opstr, present := opStr[op]; present {
		return opstr
	}

	return "UNKNOWN"
}

// Translate 'watcher' package's event ops to our own.
func xlate(wOp watcher.Op) Op {
	switch wOp {
	case watcher.Create:
		return Create
	case watcher.Write:
		return Write
	case watcher.Remove:
		return Remove
	case watcher.Rename:
		return Rename
	case watcher.Chmod:
		return Chmod
	default:
		return Tbi
	}
}

func (w *Watch) handleEvents() {
	wr := w.watcher

	for {
		select {
		case event := <-wr.Event:
			w.HandlerCB(xlate(event.Op), event.Path)

		case err := <-wr.Error:
			fmt.Fprintf(os.Stderr, "%s", err)

		case <-wr.Closed:
			return
		}
	}
}

// New instantiates a new Watch struct.
func New(
	dirs, cmds []string,
	regex string,
	recursive bool,
	handler Handler) *Watch {
	w := Watch{dirs, cmds, regex, recursive, handler, watcher.New()}
	return &w
}

// Start starts watching the configured directories via Setup().
func (w *Watch) Start() error {
	wr := w.watcher

	// Start the event handler for the watcher before starting the
	// watcher.
	go w.handleEvents()

	// Start watching for file changes.
	if err := wr.Start(time.Second * 5); err != nil {
		return (fmt.Errorf("failed to start watching files"))
	}

	fmt.Println(wr.WatchedFiles())

	return nil
}

// Setup sets up a watcher on a directory for given set of files.
func (w *Watch) Setup() error {
	wr := w.watcher

	// Allow at most 1 event to be received per watching cycle.
	wr.SetMaxEvents(1)

	// Only monitor if the file is written.
	wr.FilterOps(watcher.Write)

	// Apply regex filters for file names.
	re, err := regexp.Compile(w.Regex)
	if err != nil {
		return (fmt.Errorf("couldn't compile regex: %s", w.Regex))
	}
	wr.AddFilterHook(watcher.RegexFilterHook(re, false))

	// Watch user-given directories for changes.
	for _, dir := range w.Dirs {
		if w.Recursive {
			if err := wr.AddRecursive(dir); err != nil {
				return (fmt.Errorf("failed to watch recursive dir: %s", dir))
			}
		} else {
			if err := wr.Add(dir); err != nil {
				return (fmt.Errorf("failed to watch dir: %s", dir))
			}
		}
	}

	return nil
}
