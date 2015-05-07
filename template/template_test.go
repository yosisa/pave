package template

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestIP(t *testing.T) {
	nics, err := net.Interfaces()
	assert.Nil(t, err)
	assert.True(t, len(nics) > 0)

	var ipv4, ipv6 []string
	addrs, err := nics[0].Addrs()
	assert.Nil(t, err)
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		assert.Nil(t, err)
		if ip.To4() != nil {
			ipv4 = append(ipv4, ip.String())
		} else {
			ipv6 = append(ipv6, ip.String())
		}
	}
	ipv4, ipv6 = append(ipv4, ""), append(ipv6, "")

	ifname := nics[0].Name
	r := Render("", fmt.Sprintf(`{{ipv4 "%s"}}, {{ipv6 "%s"}}`, ifname, ifname))
	assert.Equal(t, fmt.Sprintf("%s, %s", ipv4[0], ipv6[0]), r)

	r = Render("", `{{ipv4 "ethX" "lo" "lo0" "127"}}`)
	assert.Equal(t, "127.0.0.1", r)
}

func TestSplit(t *testing.T) {
	tmpl := `{{range $i, $v := env "PATH" | split ":"}}{{if $i}}, {{end}}{{.}}{{end}}`
	expect := strings.Replace(os.Getenv("PATH"), ":", ", ", -1)
	assert.Equal(t, expect, Render("", tmpl))

	tmpl = `{{range $v := "a,b,,c" | split ","}}{{.}}{{end}}`
	assert.Equal(t, "abc", Render("", tmpl))
}

func TestDefault(t *testing.T) {
	tmpl := `{{"no" | default "yes"}}`
	assert.Equal(t, "no", Render("", tmpl))

	tmpl = `{{"" | default "yes"}}`
	assert.Equal(t, "yes", Render("", tmpl))
}

func TestMath(t *testing.T) {
	assert.Equal(t, "3", Render("", "{{add 2 1}}"))
	assert.Equal(t, "1", Render("", "{{sub 2 1}}"))
	assert.Equal(t, "2", Render("", "{{mul 2 1}}"))
	assert.Equal(t, "2", Render("", "{{div 4 2}}"))
	assert.Equal(t, "3", Render("", "{{major 4}}"))
	assert.Equal(t, "3", Render("", "{{major 5}}"))
	assert.Equal(t, "2", Render("", "{{div 4 3 | ceil}}"))
	assert.Equal(t, "1", Render("", "{{div 4 3 | floor}}"))
}
