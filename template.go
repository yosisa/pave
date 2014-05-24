package main

import (
	"os"
	"path/filepath"
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
		"env":  Getenv,
		"ipv4": IPv4,
		"ipv6": IPv6,
	}
	name := filepath.Base(t.Src)
	tmpl := template.Must(template.New(name).Funcs(funcMap).ParseFiles(t.Src))

	dst, err := os.Create(t.Dst)
	if err != nil {
		return err
	}
	defer dst.Close()

	tmpl.Execute(dst, "")

	return nil
}

func Getenv(name string, defaults ...string) string {
	s := os.Getenv(name)
	if s == "" && len(defaults) > 0 {
		s = defaults[0]
	}
	return s
}
