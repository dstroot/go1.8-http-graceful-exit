package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestIndex(t *testing.T) {

	router := httprouter.New()
	router.GET("/", Index)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check content type
	if ctype := rr.Header().Get("Content-Type"); ctype != "text/html; charset=utf-8" {
		t.Errorf("content type header does not match: got %v want %v",
			ctype, "text/html; charset=utf-8")
	}

	// // Check the response body is what we expect.
	// if expected := `{"alive": true}`; rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}

func TestPage(t *testing.T) {

	router := httprouter.New()
	router.GET("/page", Page)

	req, _ := http.NewRequest("GET", "/page", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check content type
	if ctype := rr.Header().Get("Content-Type"); ctype != "text/html; charset=utf-8" {
		t.Errorf("content type header does not match: got %v want %v",
			ctype, "text/html; charset=utf-8")
	}

	// // Check the response body is what we expect.
	// if expected := `{"alive": true}`; rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}

func TestHello(t *testing.T) {

	r := httprouter.New()
	r.GET("/hello/:name", Hello)

	req, _ := http.NewRequest("GET", "/hello/Dan", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check content type
	if ctype := rr.Header().Get("Content-Type"); ctype != "text/plain; charset=utf-8" {
		t.Errorf("content type header does not match: got %v want %v",
			ctype, "text/plain; charset=utf-8")
	}

	// Check the response body is what we expect.
	if expected := "Hello, Dan!\n"; rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// func TestIndexHandler(t *testing.T) {
// 	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
// 	// pass 'nil' as the third parameter.
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
//
// 	// wrap your handler so that it can be accessed as an http.HandlerFunc
// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		Index(w, r, httprouter.Params{})
// 	})
//
// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)
//
// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
//
// 	// Check content type
// 	if ctype := rr.Header().Get("Content-Type"); ctype != "text/html; charset=utf-8" {
// 		t.Errorf("content type header does not match: got %v want %v",
// 			ctype, "text/html; charset=utf-8")
// 	}
//
// 	fmt.Println(rr.Body.String())
//
// 	// Check the response body is what we expect.
// 	if expected := `{"alive": true}`; rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
