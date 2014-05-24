package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

var testText = `{{env "USER"}}
{{env "EMPTY_VALUE_FOR_TEST" "default"}}
{{ipv4 "127.0"}}
`

func TestExecute(t *testing.T) {
	// Create a temporary file and write test data
	f, err := ioutil.TempFile("", "pave-test-")
	assert.Nil(t, err)
	defer os.Remove(f.Name())

	f.Write([]byte(testText))
	f.Close()

	tmpl := NewTemplate(f.Name())
	tmpl.Execute()
	defer os.Remove(tmpl.Src)

	b, err := ioutil.ReadFile(tmpl.Src)
	assert.Nil(t, err)
	assert.Equal(t, testText, string(b))

	expect := fmt.Sprintf("%s\ndefault\n127.0.0.1\n", os.Getenv("USER"))
	b, err = ioutil.ReadFile(f.Name())
	assert.Nil(t, err)
	assert.Equal(t, expect, string(b))
}
