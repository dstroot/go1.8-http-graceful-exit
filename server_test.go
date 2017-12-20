package main

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

// define a handler
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func TestNewServer(t *testing.T) {

	s := NewServer(":8000", http.HandlerFunc(hello))

	var r = reflect.TypeOf(s)
	result := "*main.Server"
	// Check the type is what we expect.
	if r.String() != result {
		t.Errorf("handler returned wrong status code: got %v want %v",
			r.String(), result)
	}

	// go func() {
	// 	s.Run()
	// }()
	//
	// time.Sleep(2 * time.Second)
	//
	// p, err := os.FindProcess(syscall.Getpid())
	// if err != nil {
	// 	t.Errorf("err")
	// }
	//
	// err = p.Signal(syscall.SIGINT)
	// if err != nil {
	// 	t.Errorf("err")
	// }
}
