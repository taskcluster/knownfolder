//+build !windows

package main

import (
	"log"
	"runtime"
)

func main() {
	log.Printf("knownfolder only runs on Windows (you are running on %v)", runtime.GOOS)
}
