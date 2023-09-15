package components

import (
	"errors"
	"html/template"
	"io"
	"lommix/wichtelbot/server/store"
	"os"
	"path/filepath"
)

const (
	ComponentsDir = "./templates/components"
	PagesDir      = "./templates/pages"
)

type TemplateContext struct {
	Snippets map[string]interface{}
	User     store.User
}

func (ctx *TemplateContext) IsLoggedIn() bool {
	return ctx.User.Id > 0
}

type Templates struct {
	shared *template.Template
	pages  map[string]*template.Template
}

func (tmpl *Templates) Render(writer io.Writer, name string, data any) error {
	template := tmpl.pages[name]
	if template != nil {
		return template.ExecuteTemplate(writer, name, data)
	}

	template = tmpl.shared.Lookup(name)
	if template != nil {
		return template.Execute(writer, data)
	}

	return errors.New("Template not found")
}

// Load templates form dir
func (tmpl *Templates) Load() error {
	tmpl.shared = template.New("")
	paths, err := getValidTemplatesInDir(ComponentsDir)
	if err != nil {
		return err
	}
	_, err = tmpl.shared.ParseFiles(paths...)

	if err != nil {
		return err
	}

	tmpl.pages = make(map[string]*template.Template)

	paths, err = getValidTemplatesInDir(PagesDir)
	if err != nil {
		return err
	}

	for _, path := range paths {
		t, err := tmpl.shared.Clone()
		if err != nil {
			return err
		}
		_, err = t.ParseFiles(path)
		if err != nil {
			return err
		}

		tmpl.pages[filepath.Base(path)] = t
	}

	return nil
}

func getValidTemplatesInDir(path string) ([]string, error) {
	var out []string
	files, err := os.ReadDir(path)
	if err != nil {
		return out, err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".html" {
			out = append(out, path+"/"+f.Name())
		}
	}
	return out, nil
}
