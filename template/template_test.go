package template

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
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
