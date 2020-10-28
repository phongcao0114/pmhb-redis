package klog

import (
	"io"
	"sync"
)

// This Go file holds the map of currently opening log file writers to prevent user
// open more than one writer on the same file.
// Also, we use this map to close all opening writers when klog.Close() was called
// (typically at the end of process).
var (
	openingWriters = make(map[string]io.WriteCloser)
	mux            = &sync.Mutex{}
)

func isWriterAlreadyOpened(path string) bool {
	mux.Lock()
	_, ok := openingWriters[path]
	mux.Unlock()
	return ok
}

func trackOpeningLogWriter(path string, writer io.WriteCloser) {
	mux.Lock()
	openingWriters[path] = writer
	mux.Unlock()
}

// Close closes all opening log file writers.
func Close() error {
	mux.Lock()
	defer mux.Unlock()
	for _, w := range openingWriters {
		if w != nil {
			w.Close()
		}
	}
	return nil
}
