package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIsDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "TestIsDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)

	if !isDir(dir) {
		t.Errorf("directory expected: %s", dir)
	}
}

func TestIsFile(t *testing.T) {
	f, err := ioutil.TempFile("", "TestIsFile")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	defer os.Remove(name)

	if !isFile(name) {
		t.Errorf("file expected: %s", name)
	}
}
