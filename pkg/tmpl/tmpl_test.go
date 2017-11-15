package tmpl

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildOptions(t *testing.T) {

	// instantiate a template renderer
	render := New()

	render.buildOptions()

	// test data
	var options = []struct {
		optName string
		value   string
	}{
		{"TemplateDirectory", "templates"},
		{"TemplateLayoutPath", "layouts"},
		{"TemplatePartialPath", "partials"},
		{"TemplatePagePath", "pages"},
		{"TemplateExtension", ".html"},
		{"TemplateBaseLayout", "layout"},
	}

	// test setting default options
	for _, o := range options {
		// check default options
		if o.optName == "TemplateDirectory" {
			if render.opts.TemplateDirectory != o.value {
				t.Error(o.optName, " returned ", render.opts.TemplateDirectory, " instead of ", o.value)
			}
		}

		if o.optName == "TemplateLayoutPath" {
			if render.opts.TemplateLayoutPath != o.value {
				t.Error(o.optName, " returned ", render.opts.TemplateDirectory, " instead of ", o.value)
			}
		}

		if o.optName == "TemplatePartialPath" {
			if render.opts.TemplatePartialPath != o.value {
				t.Error(o.optName, " returned ", render.opts.TemplateDirectory, " instead of ", o.value)
			}
		}

		if o.optName == "TemplatePagePath" {
			if render.opts.TemplatePagePath != o.value {
				t.Error(o.optName, " returned ", render.opts.TemplateDirectory, " instead of ", o.value)
			}
		}

		if o.optName == "TemplateExtension" {
			if render.opts.TemplateExtension != o.value {
				t.Error(o.optName, " returned ", render.opts.TemplateDirectory, " instead of ", o.value)
			}
		}

		if o.optName == "TemplateBaseLayout" {
			if render.opts.TemplateBaseLayout != o.value {
				t.Error(o.optName, " returned ", render.opts.TemplateDirectory, " instead of ", o.value)
			}
		}
	}
}

func TestParseTemplates(t *testing.T) {

	// instantiate a template renderer
	render := New(
		Options{
			TemplateDirectory: "../../templates",
		},
	)

	// parse templates
	render.parseTemplates()

	// make sure our templates map is not empty
	if len(render.templates) == 0 {
		t.Error("template map is empty")
	}

	for key := range render.templates {
		fmt.Println("Key:", key)
	}
}

func TestRenderTemplates(t *testing.T) {

	// instantiate a template renderer
	render := New(
		Options{
			TemplateDirectory: "../../templates",
		},
	)

	// page data to render page
	data := map[string]interface{}{
		"title": "Page 2",
		"Key":   "Value",
		"Slice": []string{"One", "Two", "Three"},
	}

	// render page template
	rr := httptest.NewRecorder()
	err := render.RenderTemplate(rr, "page.html", data)
	if err != nil {
		t.Error("rendering error")
	}

	// Check the response body contains what we expect
	expected := "Bootstrap &middot; Page 2"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	fmt.Println("Output:", rr)

}
