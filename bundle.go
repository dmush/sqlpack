package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

func bundle(path string) (text string, err error) {
	tmpl := filepath.Base(path)
	funcMap := template.FuncMap{
		"include": include(path),
	}
	t, err := template.New(tmpl).Funcs(funcMap).ParseFiles(path)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	err = t.ExecuteTemplate(buf, tmpl, struct {
		FileName string
	}{path})
	if err != nil {
		return
	}
	text = buf.String()
	text = fmt.Sprintf("-- # start # %s\n%s\n-- # end # %s", path, text, path)
	return
}

func include(parent string) func(...string) (string, error) {
	dir := filepath.Dir(parent)
	ext := filepath.Ext(parent)
	return func(paths ...string) (text string, err error) {
		for i, path := range paths {
			if !filepath.IsAbs(path) {
				path = filepath.Join(dir, path)
			}
			if filepath.Ext(path) == "" {
				path = path + ext
			}
			var part string
			part, err = bundle(path)
			if err != nil {
				return
			}
			if i > 0 {
				text += "\n"
			}
			text += part
		}
		return
	}
}
