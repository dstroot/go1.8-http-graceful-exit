// Package tmpl is a package that handles template rendering the way I like to do it.
// There are others out there.  This one is pretty good:
// https://github.com/thedevsaddam/renderer/blob/master/renderer.go
// it also handles multiple types of rendering other than just templates.
package tmpl

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/oxtoacart/bpool"
)

const (
	defaultTemplateDirectory   = "templates"
	defaultTemplateLayoutPath  = "layouts"
	defaultTemplatePartialPath = "partials"
	defaultTemplatePagePath    = "pages"
	defaultTemplateExtension   = ".html"
	defaultTemplateBaseLayout  = "layout"
)

type (
	// Options describes an option type
	Options struct {
		TemplateDirectory   string // = "templates"
		TemplateLayoutPath  string // = "layouts"
		TemplatePartialPath string // = "partials"
		TemplatePagePath    string // = "pages"
		TemplateExtension   string // = ".html"
		TemplateBaseLayout  string // = "layout"
	}

	// Render describes a renderer type
	Render struct {
		opts      Options
		templates map[string]*template.Template
		bufpool   *bpool.BufferPool
	}
)

// New return a new instance of a pointer to Render
func New(opts ...Options) *Render {
	var opt Options
	if opts != nil {
		opt = opts[0]
	}

	r := &Render{
		opts:      opt,
		templates: make(map[string]*template.Template),
		bufpool:   bpool.NewBufferPool(64),
	}

	// build options for the Render instance
	r.buildOptions()

	// if TemplateDir is not empty then call the parseTemplates
	if r.opts.TemplateDirectory != "" {
		r.parseTemplates()
	}

	return r
}

// buildOptions builds the options/Sets default values for options
func (r *Render) buildOptions() {
	if r.opts.TemplateDirectory == "" {
		r.opts.TemplateDirectory = defaultTemplateDirectory
	}

	if r.opts.TemplateLayoutPath == "" {
		r.opts.TemplateLayoutPath = defaultTemplateLayoutPath
	}

	if r.opts.TemplatePartialPath == "" {
		r.opts.TemplatePartialPath = defaultTemplatePartialPath
	}

	if r.opts.TemplatePagePath == "" {
		r.opts.TemplatePagePath = defaultTemplatePagePath
	}

	if r.opts.TemplateExtension == "" {
		r.opts.TemplateExtension = defaultTemplateExtension
	}

	if r.opts.TemplateBaseLayout == "" {
		r.opts.TemplateBaseLayout = defaultTemplateBaseLayout
	}
}

// parseTemplates parses all the templates
func (r *Render) parseTemplates() {

	// get layouts
	layouts, err := filepath.Glob(filepath.Join(r.opts.TemplateDirectory, r.opts.TemplateLayoutPath, "*"+r.opts.TemplateExtension))
	if err != nil {
		log.Fatal(err)
	}

	// get includes
	includes, err := filepath.Glob(filepath.Join(r.opts.TemplateDirectory, r.opts.TemplatePartialPath, "*"+r.opts.TemplateExtension))
	if err != nil {
		log.Fatal(err)
	}

	// get pages
	pages, err := filepath.Glob(filepath.Join(r.opts.TemplateDirectory, r.opts.TemplatePagePath, "*"+r.opts.TemplateExtension))
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map - one for each page
	files := append(layouts, includes...)
	for _, page := range pages {
		files = append(files, page)
		// TODO: add FuncMap
		// for _, fm := range r.opts.FuncMap {
		// 	tmpl.Funcs(fm)
		// }
		r.templates[filepath.Base(page)] = template.Must(template.ParseFiles(files...))
	}
}

// RenderTemplate renders our template.
// RenderTemplate is a wrapper around template.ExecuteTemplate.
// It writes into a bytes.Buffer before writing to the http.ResponseWriter to catch
// any errors resulting from populating the template.
func (r *Render) RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("the template %s does not exist", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := r.bufpool.Get()
	defer r.bufpool.Put(buf)

	// render the template and check for errors
	err := tmpl.ExecuteTemplate(buf, r.opts.TemplateBaseLayout, data)
	if err != nil {
		return err
	}

	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}
