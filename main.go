package main

import (
	"log"
	"os"
	"runtime/trace"
	"time"
)

var _x int

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := trace.Start(f); err != nil {
		log.Fatal(err)
	}
	go func() {
		// runtime.PreemptNS(1000)
		for {
			_x++
		}
	}()
	<-time.After(3 * time.Second)
	defer trace.Stop()
}
