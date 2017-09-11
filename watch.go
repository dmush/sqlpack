package main

import (
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

type watchFunc func(path string)

func watch(dir string, fn watchFunc) (err error) {
	w := watcher.New()

	w.FilterOps(watcher.Create, watcher.Write)

	go func() {
		for {
			select {
			case ev := <-w.Event:
				fn(ev.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				log.Fatalln("watcher closed")
			}
		}
	}()

	if err = w.AddRecursive(dir); err != nil {
		return
	}

	err = w.Start(time.Millisecond * 500)
	return
}
