package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var pathToTemplates = "./web/templates"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *TemplateData) error {
	var tc map[string]*template.Template
	tc, _ = CreateTemplate()

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser:", err)
		return err
	}

	return nil
}

func CreateTemplate() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*_page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*_layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*_layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
