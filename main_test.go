package main

import (
	"os"
	"testing"
)

// this doesnt really test anything but is necessary for
// `make cover` to work since each file needs a corresponding
// test file.
func TestMain(m *testing.M) {
	// setup()
	code := m.Run()
	// shutdown()
	os.Exit(code)
}
