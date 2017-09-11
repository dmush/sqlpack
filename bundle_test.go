package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var bundleFile1Text = "create table file_1 (id uuid, name text);"
var bundleFile2Text = "create table file_2 (id uuid, name text);"
var bundleFile3Text = "create table file_3 (id uuid, name text);"
var bundleFile4Text = "create table file_4 (id uuid, name text);"

func tempFile(t *testing.T, name, text string) (file *os.File) {
	file, err := ioutil.TempFile(os.TempDir(), name)
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.Write([]byte(text))
	if err != nil {
		t.Fatal(err)
	}

	return
}

func TestBundle(t *testing.T) {
	file4 := tempFile(t, "bundle_file_4_", bundleFile4Text)
	defer os.Remove(file4.Name())

	file3 := tempFile(t, "bundle_file_3_", bundleFile3Text)
	defer os.Remove(file3.Name())

	f4path, err := filepath.Rel(os.TempDir(), file4.Name())
	if err != nil {
		t.Fatal(err)
	}
	file2 := tempFile(t, "bundle_file_2_", fmt.Sprintf(
		`%s
{{ include "%s" }}`,
		bundleFile2Text,
		f4path,
	))
	defer os.Remove(file2.Name())

	f2path, err := filepath.Rel(os.TempDir(), file2.Name())
	if err != nil {
		t.Fatal(err)
	}
	f3path, err := filepath.Rel(os.TempDir(), file3.Name())
	if err != nil {
		t.Fatal(err)
	}
	file1 := tempFile(t, "bundle_file_1_", fmt.Sprintf(
		`%s
{{ include "%s" "%s" }}`,
		bundleFile1Text,
		f2path,
		f3path,
	))
	defer os.Remove(file1.Name())

	text, err := bundle(file1.Name())
	if err != nil {
		t.Fatal(err)
	}

	bundleText := "" +
		"-- # start # " + file1.Name() +
		"\n" + bundleFile1Text +
		"\n-- # start # " + file2.Name() +
		"\n" + bundleFile2Text +
		"\n-- # start # " + file4.Name() +
		"\n" + bundleFile4Text +
		"\n-- # end # " + file4.Name() +
		"\n-- # end # " + file2.Name() +
		"\n-- # start # " + file3.Name() +
		"\n" + bundleFile3Text +
		"\n-- # end # " + file3.Name() +
		"\n-- # end # " + file1.Name()

	if text != bundleText {
		t.Errorf("invalid file content")
	}
}
