package scan_test

import (
	"."
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
	err = scan.Scan(a, &s)
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
	err = scan.Scan(a, &b)
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
	err = scan.Scan(a, &f)
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
	err = scan.ScanTree(a, "/foo/bar", &s)
	if err != nil {
		t.Fatal(err)
	}
	if s != "baz" {
		t.Fatalf("Expected %v for Scan, but %v:", "baz", s)
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
	err = scan.ScanTree(a, "/foo/bar", &s)
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
	err = scan.ScanTree(a, "/foo/bar", &f)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]float64{3, 2, 1}, f) {
		t.Fatalf("Expected %v for Scan, but %v:", `[3, 2, 1]`, f)
	}
}

func TestScanJSON(t *testing.T) {
	s := `{"foo":{"bar": [3, 2, 1]}}`
	var f []float64
	err := scan.ScanJSON(strings.NewReader(s), "/foo/bar", &f)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]float64{3, 2, 1}, f) {
		t.Fatalf("Expected %v for Scan, but %v:", `[3, 2, 1]`, f)
	}
}

func TestScanPanic(t *testing.T) {
	var b bool
	err := scan.Scan(nil, &b)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
}

func TestScanJSONError(t *testing.T) {
	err := scan.ScanJSON(nil, "", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	sr := strings.NewReader("")
	err = scan.ScanJSON(sr, "/", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
}

func TestInvalidPath(t *testing.T) {
	err := scan.Scan(nil, nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanTree(nil, "", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(nil, "", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanTree(nil, "a", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(nil, "a", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanTree(nil, "/[a]", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(nil, "/[a]", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanTree(nil, "/a", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(nil, "/a", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}

	s := `{"foo":{"bar": [3, 2, 1]}}`
	var f []float64
	err = scan.ScanJSON(strings.NewReader(s), "/fooo/bar", &f)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(strings.NewReader(s), "/foo[999999999999999999999999999999999999999]", &f)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(strings.NewReader(s), "/foo/bar[0]/[0]", &f)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(strings.NewReader(s), "/foo/bar[20]", &f)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(strings.NewReader(s), "/foo[9]", &f)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
	err = scan.ScanJSON(strings.NewReader(s), "[9]", nil)
	if err == nil {
		t.Fatalf("Expected error but not")
	}
}
