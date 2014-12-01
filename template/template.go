package template

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

	name := filepath.Base(t.Src)
	text := Render(name, string(b))

	err = ioutil.WriteFile(t.Dst, []byte(text), 0644)
	return err
}

func Render(name, text string) string {
	funcMap := template.FuncMap{
		"env":   Getenv,
		"ipv4":  IPv4,
		"ipv6":  IPv6,
		"split": Split,
	}

	tmpl := template.Must(template.New(name).Funcs(funcMap).Parse(text))

	w := new(bytes.Buffer)
	tmpl.Execute(w, "")
	return w.String()
}

func Getenv(name string, defaults ...string) string {
	s := os.Getenv(name)
	if s == "" && len(defaults) > 0 {
		s = defaults[0]
	}
	return s
}

func Split(sep, s string) []string {
	items := strings.Split(s, sep)
	result := make([]string, 0, len(items))
	for _, v := range items {
		if val := strings.Trim(v, " "); val != "" {
			result = append(result, val)
		}
	}
	return result
}
