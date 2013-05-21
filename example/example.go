package main

import (
	"github.com/mattn/go-scan"
	"strings"
)

var js = strings.NewReader(`
{
	"foo": {
		"bar": [
			{
				"baz": "bbb",
				"noo": 3 
			},
			{
				"maz": true,
				"moo": ["foo", "bar"]
			}
		],
		"boo": {
			"bag": "ddd",
			"bug": "ccc"
		}
	}
}
`)

func main() {
	var s []string
	if err := scan.ScanJSON(js, "/foo/bar[1]/moo", &s); err != nil {
		println(err.Error())
	}
	println(s[0]) // should be "foo"
	println(s[1]) // should be "bar"
}
