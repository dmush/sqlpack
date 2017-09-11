package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	dir := os.TempDir()
	done := make(chan bool)
	watched := false
	go func() {
		err := watch(dir, func(_ string) {
			watched = true
			done <- true
		})
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	f, err := ioutil.TempFile(dir, "watch_")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.WriteString("TestWatch"); err != nil {
		t.Fatal(err)
	}
	f.Close()
	go func() {
		time.Sleep(500 * time.Millisecond)
		done <- true
	}()
	<-done
	if !watched {
		t.Error("Watch failed")
	}
}
