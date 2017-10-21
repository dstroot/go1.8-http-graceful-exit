// https://tylerchr.blog/golang-18-whats-coming/

// It is now possible to call srv.Close() to halt an
// http.Server immediately, or srv.Shutdown(ctx) to stop
// and gracefully drain the server of connections

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

// var templates map[string]*template.Template
var bufpool *bpool.BufferPool
var templates map[string]*template.Template

// var templates *template.Template

const (
	templateLayoutPath  = "templates/layouts/"
	templatePartialPath = "templates/partials/"
	templatePagePath    = "templates/pages/"
	templateExtension   = ".html"
	templateBaseLayout  = "layout.html"
	templatePath        = "test/"
)

// create a buffer pool
func initBufferPool() {
	bufpool = bpool.NewBufferPool(64)
	log.Println("buffer allocation successful")
}

// Load templates on program initialisation
func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir := "test/"

	layouts, err := filepath.Glob(templatesDir + "layouts/*.html")
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}

	log.Println(templates)
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
// It writes into a bytes.Buffer before writing to the http.ResponseWriter to catch
// any errors resulting from populating the template.
func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.ExecuteTemplate(buf, "layout.html", data)
	if err != nil {
		return err
	}

	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}

// func loadTemplates() {
// 	// if templates == nil {
// 	// 	templates = make(map[string]*template.Template)
// 	// }
// 	//
// 	var allFiles []string
// 	var err error
// 	//
// 	// get layouts
// 	layoutFiles, err := filepath.Glob(templateLayoutPath + "*" + templateExtension)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// get partials
// 	includeFiles, err := filepath.Glob(templatePartialPath + "*" + templateExtension)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// get pages
// 	pageFiles, err := filepath.Glob(templatePagePath + "*" + templateExtension)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// append all lists
// 	allFiles = append(allFiles, includeFiles...)
// 	allFiles = append(allFiles, layoutFiles...)
// 	allFiles = append(allFiles, pageFiles...)
//
// 	log.Println(allFiles)
//
// 	// mainTemplate := template.New("main")
//
// 	// mainTemplate, err = mainTemplate.Parse(mainTmpl)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
//
// 	// // load partials
// 	// for _, file := range includeFiles {
// 	// 	fileName := filepath.Base(file)
// 	// 	// files := append(allFiles, file)
// 	// 	t, err := template.ParseFiles(file)
// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}
// 	// 	templates[fileName] = t
// 	// }
// 	//
// 	// // load layouts
// 	// for _, file := range layoutFiles {
// 	// 	fileName := filepath.Base(file)
// 	// 	// files := append(allFiles, file)
// 	// 	t, err := template.ParseFiles(file)
// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}
// 	// 	templates[fileName] = t
// 	// 	// templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
// 	// }
//
// 	// // load pages
// 	// for _, file := range pageFiles {
// 	// 	fileName := filepath.Base(file)
// 	// 	// t, err := template.ParseFiles(file)
// 	// 	// if err != nil {
// 	// 	// 	log.Fatal(err)
// 	// 	// }
// 	// 	// templates[fileName] = t
// 	// 	// files := append(allFiles, file)
// 	// 	// templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
// 	// 	//
// 	// 	templates[fileName] = template.Must(templates[fileName].ParseFiles(allFiles...))
// 	// }
//
// 	// log.Println(allFiles)
//
// 	// for key, value := range templates {
// 	// 	fmt.Println("Key:", key, "Value:", value)
// 	// }
// 	// log.Println("templates loading successful")
//
// 	// var allFiles []string
// 	// files, err := ioutil.ReadDir("./templates")
// 	// if err != nil {
// 	//     fmt.Println(err)
// 	// }
// 	// for _, file := range files {
// 	//     filename := file.Name()
// 	//     if strings.HasSuffix(filename, ".tmpl") {
// 	//         allFiles = append(allFiles, "./templates/"+filename)
// 	//     }
// 	// }
//
// 	templates, err = template.ParseFiles(allFiles...) // parses all .tmpl files in the 'templates' folder
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// }

// func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
// 	// // get the template from the templates map
// 	// tmpl, ok := templates[name]
// 	// if !ok {
// 	// 	http.Error(w, fmt.Sprintf("The template %s does not exist.", name),
// 	// 		http.StatusInternalServerError)
// 	// }
//
// 	// tmpl := templates.Lookup("layout.html")
// 	// if tmpl == nil {
// 	// 	http.Error(w, fmt.Sprintf("The template %s does not exist.", name),
// 	// 		http.StatusInternalServerError)
// 	// }
// 	//
// 	// buf := bufpool.Get()
// 	// defer bufpool.Put(buf)
// 	//
// 	// // err := tmpl.Execute(buf, data)
// 	// err := templates.ExecuteTemplate(buf, "layout", data)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// }
// 	//
// 	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	// buf.WriteTo(w)
// }
