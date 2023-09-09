package server

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type Templates struct {
	Dir   string
	store *template.Template
}

func (tmpl *Templates) Refresh() error{
	tmpl.store = template.New("")
	return tmpl.Load()
}

func (tmpl *Templates) Load() error {
	err := filepath.WalkDir(tmpl.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			relPath, err := filepath.Rel(tmpl.Dir, path)
			if err != nil {
				return err
			}
			tmpl.store, err = tmpl.store.ParseFiles(path)
			if err != nil {
				return err
			}
			fmt.Printf("Loaded template: %s\n", relPath)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
