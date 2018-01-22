package info

import (
	"net/http"
	"net/http/httptest"
	"testing"
	// . "github.com/smartystreets/goconvey/convey"
)

func TestGetPort(t *testing.T) {
	// Check the port is what we expect.
	if port := getPort(); port != defaultPort {
		t.Errorf("Wrong port: got %v want %v",
			port, defaultPort)
	}
}

func TestInit(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 	// Check the response body is what we expect.
	// 	expected := `{
	//     "HostName": "Dan-iMac-DCPSC04BH3GY",
	//     "IPAddress": "192.168.15.87",
	//     "Port": "",
	//     "Program": "Info.Test",
	//     "BuildTime": "unset",
	//     "Commit": "unset",
	//     "Version": "unset",
	//     "GoVersion": "go1.9.2",
	//     "RunTime": "0s"
	// }`
	// 	if rr.Body.String() != expected {
	// 		t.Errorf("handler returned unexpected body: got %v want %v",
	// 			rr.Body.String(), expected)
	// 	}
}
