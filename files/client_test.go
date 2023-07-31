package files

import (
	"os"
	"path/filepath"
	"testing"
)

var testDataDir string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testDataDir = filepath.Join(filepath.Dir(dir), "testdata")
}

func TestRandomFile(t *testing.T) {
	c := NewClient("localhost", testDataDir)

	if err := c.Init(testDataDir); err != nil {
		t.Fatal(err)
	}

	_, err := c.RandomFile()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddRemoveFolder(t *testing.T) {
	c := NewClient("localhost", testDataDir)

	if err := c.Init(testDataDir); err != nil {
		t.Fatal(err)
	}

	fullLength := len(c.Folders)

	if fullLength != 2 {
		t.Fatalf("expected length to be 2, was %d", fullLength)
	}

	if err := c.removeFolder(testDataDir); err != nil {
		t.Fatal(err)
	}

	emptyLength := len(c.Folders)

	t.Log(emptyLength)

	if emptyLength != 0 {
		t.Fatalf("expected length to be 0, was %d", emptyLength)
	}
}
