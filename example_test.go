package scan_test

import (
	"fmt"
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

func Example() {
	var s []string
	if err := scan.ScanJSON(js, "/foo/bar[1]/moo", &s); err != nil {
		println(err.Error())
	}
	fmt.Println(s[0]) // should be "foo"
	fmt.Println(s[1]) // should be "bar"
}
