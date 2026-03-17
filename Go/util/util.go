package util

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	ch := make(chan struct{})

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-ch:
		return false
	case <-time.After(timeout):
		return true
	}
}

func SourceInfo() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "<unknown>"
	}
	return fmt.Sprintf("%v:%v", filepath.Base(file), line)
}
