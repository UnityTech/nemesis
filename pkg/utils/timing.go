package utils

import (
	"fmt"
	"time"
)

// Elapsed sets up a goroutine to indicate how long it took to run a segment of code
func Elapsed(what string) func() {
	if *flagDebug {
		start := time.Now()
		return func() {
			fmt.Printf("%s took %v\n", what, time.Since(start))
		}
	}
	return func() {}
}
