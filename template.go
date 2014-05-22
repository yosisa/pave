package main

import (
	"os"
	"text/template"
)

const suffix = ".tpl"

type Template struct {
	Src string
	Dst string
}

func NewTemplate(path string) *Template {
	return &Template{
		Src: path + suffix,
		Dst: path,
	}
}

func (t *Template) Execute() error {
	if _, err := os.Stat(t.Src); os.IsNotExist(err) {
		if err = os.Rename(t.Dst, t.Src); err != nil {
			return err
		}
	}

	funcMap := template.FuncMap{
		"env": os.Getenv,
	}
	tmpl := template.Must(template.New(t.Src).Funcs(funcMap).ParseFiles(t.Src))

	dst, err := os.Create(t.Dst)
	if err != nil {
		return err
	}
	defer dst.Close()

	tmpl.Execute(dst, "")

	return nil
}
