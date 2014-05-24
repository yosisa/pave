package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const suffix = ".tpl"

type Template struct {
	Src  string
	Dst  string
	tmpl *template.Template
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

	b, err := ioutil.ReadFile(t.Src)
	if err != nil {
		return err
	}
	t.parse(string(b))

	dst, err := os.Create(t.Dst)
	if err != nil {
		return err
	}
	defer dst.Close()

	t.tmpl.Execute(dst, "")

	return nil
}

func (t *Template) parse(text string) {
	funcMap := template.FuncMap{
		"env":  Getenv,
		"ipv4": IPv4,
		"ipv6": IPv6,
	}

	name := filepath.Base(t.Src)
	t.tmpl = template.Must(template.New(name).Funcs(funcMap).Parse(text))
}

func Getenv(name string, defaults ...string) string {
	s := os.Getenv(name)
	if s == "" && len(defaults) > 0 {
		s = defaults[0]
	}
	return s
}
