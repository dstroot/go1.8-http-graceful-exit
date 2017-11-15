package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/dstroot/go1.8-http-graceful-exit/pkg/tmpl"
	"github.com/julienschmidt/httprouter"
)

/**
 * Handlers
 */

// Index handler handles GET /
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// If you want to debug your HTTP requests, then import the
	// net/http/httputil package, and invoke DumpRequest
	// with the parameter *http.Request and a boolean to specify if you
	// want to dump the request body as well. The function returns a
	// []byte, error. You could use it like this:
	dump := func(r *http.Request) {
		output, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println("Error dumping request:", err)
			return
		}
		fmt.Println(string(output))
	}

	// The function call will dump your request method, URI with
	// query parameters, headers and request body if you have one.
	dump(r)

	// page data to render page
	data := map[string]interface{}{
		"title": "The most popular HTML, CSS, and JS library in the world.",
		"Key":   "Value",
		"Slice": []string{"One", "Two", "Three"},
	}

	// render page template
	err := tmpl.RenderTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Error rendering:", err)
		return
	}
}

// Page handles page
func Page(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// page data to render page
	data := map[string]interface{}{
		"title": "Page 2",
		"Key":   "Value",
		"Slice": []string{"One", "Two", "Three"},
	}

	// render page template
	err := tmpl.RenderTemplate(w, "page.html", data)
	if err != nil {
		log.Fatalln(err)
	}
}

// Hello handler handles GET /hello/:name
func Hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
}
