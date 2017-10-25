package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/oxtoacart/bpool"
)

/**
 * Template Handling
 */

var bufpool *bpool.BufferPool
var templates map[string]*template.Template

const (
	templateLayoutPath  = "layouts"
	templatePartialPath = "partials"
	templatePagePath    = "pages"
	templateExtension   = ".html"
	templateBaseLayout  = "layout"
	templateDirectory   = "templates"
)

// create a buffer pool
func initBufferPool() {
	bufpool = bpool.NewBufferPool(64)
	// log.Println("buffer allocation successful")
}

// Load templates on program initialisation
func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layouts, err := filepath.Glob(templateDirectory + "/" + templateLayoutPath + "/*" + templateExtension)
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(templateDirectory + "/" + templatePartialPath + "/*" + templateExtension)
	if err != nil {
		log.Fatal(err)
	}

	pages, err := filepath.Glob(templateDirectory + "/" + templatePagePath + "/*" + templateExtension)
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, page := range pages {
		files := append(layouts, includes...)
		files = append(files, page)
		templates[filepath.Base(page)] = template.Must(template.ParseFiles(files...))
	}

	// for key := range templates {
	// 	fmt.Println("Page: ", key)
	// }
	// log.Println("templates loading successful")
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
// It writes into a bytes.Buffer before writing to the http.ResponseWriter to catch
// any errors resulting from populating the template.
func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("the template %s does not exist", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	// render the template
	err := tmpl.ExecuteTemplate(buf, templateBaseLayout, data)
	if err != nil {
		return err
	}

	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
