package main

import (
	"encoding/json"
	"github.com/mattn/go-scan"
)

func main() {
	var a interface{}
	json.Unmarshal([]byte(`
{
	"foo": ["baz"]
}
`), &a)

	var s []string
	if err := scan.ScanTree(a, "/foo", &s); err != nil {
		println(err.Error())
	}
	println(s[0])
}
