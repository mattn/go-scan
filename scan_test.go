package scan

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var ScanTests = []struct {
	j string
	p string
	v interface{}
}{
	{
		"foo",
		"bar",
		"baz",
	},
}

func TestScanString(t *testing.T) {
	var a interface{}
	j :=
		`
"foo"
`
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		t.Fatal(err)
	}
	var s string
	err = Scan(a, &s)
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatalf("Expected %v for Scan, but %v:", "foo", s)
	}
}

func TestScanBool(t *testing.T) {
	var a interface{}
	j :=
		`
true
`
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		t.Fatal(err)
	}
	var b bool
	err = Scan(a, &b)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatalf("Expected %v for Scan, but %v:", true, b)
	}
}
func TestScanFloat64(t *testing.T) {
	var a interface{}
	j :=
		`
123
`
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		t.Fatal(err)
	}
	var f float64
	err = Scan(a, &f)
	if err != nil {
		t.Fatal(err)
	}
	if f != 123.0 {
		t.Fatalf("Expected %f for Scan, but %f:", 123.0, f)
	}
}

func TestScanTreeMap(t *testing.T) {
	var a interface{}
	j :=
		`
{"foo":{"bar": "baz"}}
`
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		t.Fatal(err)
	}
	var s string
	err = ScanTree(a, "/foo/bar", &s)
	if err != nil {
		t.Fatal(err)
	}
	if s != "baz" {
		t.Fatalf("Expected %v for Scan, but %v:", "baz", s)
	}
}

func TestScanTreeStringKeyMap(t *testing.T) {
	var a map[string]interface{}
	j :=
		`
{"foo":{"bar": "baz"}}
`
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		t.Fatal(err)
	}
	var s string
	err = ScanTree(a, "/foo/bar", &s)
	if err != nil {
		t.Fatal(err)
	}
	if s != "baz" {
		t.Fatalf("Expected %v for Scan, but %v:", "baz", s)
	}
}

func TestScanTreeInterfaceKeyMap(t *testing.T) {
	a := map[interface{}]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
		3: "baba",
	}
	var s string
	err := ScanTree(a, "/foo/bar", &s)
	if err != nil {
		t.Fatal(err)
	}
	if s != "baz" {
		t.Fatalf("Expected %v for Scan, but %v:", "baz", s)
	}
}

func TestScanTreeInvalidKeyMap(t *testing.T) {
	a := map[interface{}]interface{}{
		"foo": map[string]interface{}{
			"bar": "bar",
		},
		3: "baba",
	}
	var s string
	err := ScanTree(a, "/foo/baz", &s)
	if err == nil {
		t.Fatal("Expected error but not")
	}
}

func TestScanTreeInvalidMap(t *testing.T) {
	a := map[interface{}]interface{}{
		"foo": map[interface{}]interface{}{
			"bar":    func() {},
			"barbar": func() {},
		},
		3: "baba",
	}
	var s string
	err := ScanTree(a, "/foo[0]", &s)
	if err == nil {
		t.Fatal("Expected error but not")
	}
}

func TestScanTreeSliceOfString(t *testing.T) {
	var a interface{}
	j :=
		`
{"foo":{"bar": ["baz", "baba"]}}
`
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		t.Fatal(err)
	}
	var s []string
	err = ScanTree(a, "/foo/bar", &s)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{"baz", "baba"}, s) {
		t.Fatalf("Expected %v for Scan, but %v:", `["bar": "baba"]`, s)
	}
}

func TestScanTreeSliceOfFloat64(t *testing.T) {
	var a interface{}
	j :=
		`
{"foo":{"bar": [3, 2, 1]}}
`
	err := json.Unmarshal([]byte(j), &a)
	if err != nil {
		t.Fatal(err)
	}
	var f []float64
	err = ScanTree(a, "/foo/bar", &f)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]float64{3, 2, 1}, f) {
		t.Fatalf("Expected %v for Scan, but %v:", `[3, 2, 1]`, f)
	}
}

func TestScanAny(t *testing.T) {
	s := `{"foo":{"bar": [3, 2, 1]}}`
	var a Any
	err := ScanJSON(strings.NewReader(s), "/foo/bar", &a)
	if err != nil {
		t.Fatal(err)
	}
	var v interface{}
	v = interface{}([]interface{}{3.0, 2.0, 1.0})
	if !reflect.DeepEqual(v, a) {
		t.Fatalf("Expected %v for Scan, but %v:", `[3, 2, 1]`, a)
	}
}

func TestScanJSON(t *testing.T) {
	s := `{"foo":{"bar": [3, 2, 1]}}`
	var f []float64
	err := ScanJSON(strings.NewReader(s), "/foo/bar", &f)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]float64{3, 2, 1}, f) {
		t.Fatalf("Expected %v for Scan, but %v:", `[3, 2, 1]`, f)
	}
	var i int
	err = ScanJSON(strings.NewReader(s), "/foo/bar[2]", &i)
	if err != nil {
		t.Fatal(err)
	}
	if i != 1 {
		t.Fatalf("Expected %v for Scan, but %v:", 1, i)
	}
}

func TestScanPanic(t *testing.T) {
	var b bool
	err := Scan(nil, &b)
	if err == nil {
		t.Fatal("Expected error but not")
	}
}

func TestScanJSONError(t *testing.T) {
	err := ScanJSON(nil, "", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	sr := strings.NewReader("")
	err = ScanJSON(sr, "/", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
}

func TestIndexWithMap(t *testing.T) {
	var js = `
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
`
	var a interface{}
	err := json.Unmarshal([]byte(js), &a)
	if err != nil {
		t.Fatal(err)
	}
	var s string
	if err := ScanTree(a, "/foo/boo[0]", &s); err != nil {
		if err := ScanTree(a, "/foo/boo[1]", &s); err != nil {
			t.Fatal(err)
		}
	}
}

func TestScan(t *testing.T) {
	var f float32
	err := Scan(nil, &f)
	if err == nil {
		t.Fatal("Expected error but not")
	}
}

func TestInvalidPath(t *testing.T) {
	err := Scan(nil, nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanTree(nil, "", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(nil, "", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanTree(nil, "a", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(nil, "a", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanTree(nil, "/[a]", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(nil, "/[a]", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanTree(nil, "/a", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(nil, "/a", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}

	s := `{"foo":{"bar": [3, 2, 1]}}`
	var f []float64
	err = ScanJSON(strings.NewReader(s), "/fooo/bar", &f)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(strings.NewReader(s), "/foo[999999999999999999999999999999999999999]", &f)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(strings.NewReader(s), "/foo/bar[0]/[0]", &f)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(strings.NewReader(s), "/foo/bar[20]", &f)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(strings.NewReader(s), "/foo[9]", &f)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	err = ScanJSON(strings.NewReader(s), "[9]", nil)
	if err == nil {
		t.Fatal("Expected error but not")
	}
}

func TestToError(t *testing.T) {
	err := toError(1)
	if err == nil {
		t.Fatal("Expected error but not")
	}
	if err.Error() != "Unknown error" {
		t.Fatal("Expected unknown error but not")
	}
	if toError(nil) != nil {
		t.Fatal("Expected nil error but not")
	}
}

func TestSlash(t *testing.T) {
	s := `{"foo bar":{"bar/baz": [3, 2, 1]}}`
	var f []float64
	err := ScanJSON(strings.NewReader(s), `/foo bar/bar\/baz`, &f)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]float64{3, 2, 1}, f) {
		t.Fatalf("Expected %v for Scan, but %v:", `[3, 2, 1]`, f)
	}
	var i int
	err = ScanJSON(strings.NewReader(s), `/foo bar/bar\/baz[2]`, &i)
	if err != nil {
		t.Fatal(err)
	}
	if i != 1 {
		t.Fatalf("Expected %v for Scan, but %v:", 1, i)
	}
}
