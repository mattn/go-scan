go-scan
=======

[![Build Status](https://travis-ci.org/mattn/go-scan.png?branch=master)](https://travis-ci.org/mattn/go-scan)
[![Coverage Status](https://coveralls.io/repos/mattn/go-scan/badge.png?branch=HEAD)](https://coveralls.io/r/mattn/go-scan?branch=HEAD)

Easily way to get the elements via xpath like string

Usage
-----

```go
var js = strings.NewReader(`
{
	"foo": {
		"bar": [
			{
				"faz": true,
				"moo": ["goo", "mar"]
			},
			{
				"maz": true,
				"moo": ["foo", "bar"]
			}
		]
	}
}
`)
var s []string
scan.ScanJSON(js, "/foo/bar[1]/moo", &s) // s should be ["foo", "bar"]
```

Install
-------

```
go get github.com/mattn/go-scan
```

License
-------

MIT: http://mattn.mit-license.org/2013

Author
------

Yasuhiro Matsumoto (mattn.jp@gmail.com)
