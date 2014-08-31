# pave
[![Build Status](https://drone.io/github.com/yosisa/pave/status.png)](https://drone.io/github.com/yosisa/pave/latest) [![Coverage Status](https://coveralls.io/repos/yosisa/pave/badge.png?branch=master)](https://coveralls.io/r/yosisa/pave?branch=master)

pave is a tiny program which provides process management and template rendering
for configuration files before running processes. In addition to configuration
files, each command to run a process can be used as a template. It is intended
to use inside a Docker container.

pave uses [Golang template language] and extends following functions that can be
used in a template.

* env {{env KEY DEFAULT}}
    * Retrieves the value of the environment variable named by the KEY. If such
      environment variable not defined or its value is empty, DEFAULT is used.
* ipv4 {{ipv4 KEY...}}
    * Resolve the IPv4 address suitable for given KEY. KEY is a interface name
      or prefix of a IP address. KEY can be specified multiple times. In this
      case, the first matching non-empty IP address is used.
* ipv6 {{ipv6 KEY...}}
    * Same as `ipv4` but resolves IPv6 address.

[Golang template language]: http://golang.org/pkg/text/template/

## Quick start
```
pave -c 'echo {{env "USER"}} {{ipv4 "eth0" "en0"}}'
```

## Usage
```
Usage:
  pave [OPTIONS]

Application Options:
  -f, --file=         Files to be rendered
  -c, --command=      Commands to be executed
  -r, --restart=      Restart strategy (none|always|error)
  -w, --restart-wait= Duration for restarting

Help Options:
  -h, --help          Show this help message
```

## License
The MIT License (MIT)

Copyright (c) 2014 Yoshihisa Tanaka

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
