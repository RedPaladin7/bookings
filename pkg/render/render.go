package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/RedPaladin7/bookings/pkg/config"
	"github.com/RedPaladin7/bookings/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig){
	app = a
}

//AddDefaultData adds some default data which every template needs
func AddDefaultData(td *models.TemplateData) *models.TemplateData{
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData){
	var tc map[string]*template.Template
	if app.UseCache{
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok{
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

// CreateTemplateCache creates a template cache
func CreateTemplateCache() (map[string]*template.Template, error){
	myCache := map[string]*template.Template{}

	// get all of the files name *.page.tmpl
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil{
		return myCache, err
	}

	// range through all the pages
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil{
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil{
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil{
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
