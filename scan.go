package scan

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile("^([^\\[]+)?(\\[[0-9]+\\])?$")

var t1 = reflect.TypeOf((map[string]interface{})(nil))
var t2 = reflect.TypeOf((map[interface{}]interface{})(nil))

// Any provide interface to scan any types.
type Any interface{}

func toError(v interface{}) error {
	if v != nil {
		if e, ok := v.(error); ok {
			return e
		}
		if e, ok := v.(string); ok {
			return errors.New(e)
		}
		return errors.New("Unknown error")
	}
	return nil
}

// Scan work to scan any type to specified type
func Scan(v interface{}, t interface{}) (err error) {
	defer func() {
		if err == nil {
			err = toError(recover())
		}
	}()
	rt := reflect.ValueOf(t).Elem()
	rv := reflect.ValueOf(v)
	tv := rv.Type().Kind()

	if tv == reflect.Slice || tv == reflect.Array {
		if _, ok := t.(*Any); ok {
			rt.Set(rv)
			return nil
		}
		ia := rv.Interface().([]interface{})
		rt.Set(reflect.MakeSlice(rt.Type(), len(ia), len(ia)))
		for n := range ia {
			rt.Index(n).Set(rv.Index(n).Elem())
		}
	} else {
		if rv.Type().ConvertibleTo(rt.Type()) {
			rt.Set(rv.Convert(rt.Type()))
		} else {
			rt.Set(rv)
		}
	}
	return nil
}

func split(s string) []string {
	i := 0
	a := []string{}
	t := ""
	rs := []rune(s)
	l := len(rs)
	for i < l {
		r := rs[i]
		switch r {
		case '\\':
			i++
			if i < l {
				t += string(rs[i])
			}
		case '/':
			if t != "" {
				a = append(a, t)
				t = ""
			}
		default:
			t += string(r)
		}
		i++
	}
	if t != "" {
		a = append(a, t)
	}
	return a
}

// ScanTree work to scan value to specified value with the path
func ScanTree(v interface{}, p string, t interface{}) (err error) {
	defer func() {
		if err == nil {
			err = toError(recover())
		}
	}()
	if p == "" {
		return errors.New("invalid path")
	}
	var ok bool
	for _, token := range split(p) {
		sl := re.FindAllStringSubmatch(token, -1)
		if len(sl) == 0 {
			return errors.New("invalid path")
		}
		ss := sl[0]
		if ss[1] != "" {
			rv := reflect.ValueOf(v)
			rt := rv.Type()
			if rt != t1 && rv.Type().ConvertibleTo(t1) {
				v = rv.Convert(t1).Interface()
			}
			if vm, ok := v.(map[string]interface{}); ok {
				if v, ok = vm[ss[1]]; !ok {
					return errors.New("invalid path: " + ss[1])
				}
			} else {
				if rt != t2 && rv.Type().ConvertibleTo(t2) {
					v = rv.Convert(t2).Interface()
				}
				if vm, ok := v.(map[interface{}]interface{}); ok {
					if v, ok = vm[ss[1]]; !ok {
						return errors.New("invalid path: " + ss[1])
					}
				} else {
					return errors.New("invalid path: " + ss[1])
				}
			}
		}
		if ss[2] != "" {
			i, err := strconv.Atoi(ss[2][1 : len(ss[2])-1])
			if err != nil {
				return errors.New("invalid path: " + ss[2])
			}
			var vl []interface{}
			if vl, ok = v.([]interface{}); !ok {
				if vm, ok := v.(map[string]interface{}); ok {
					n, found := 0, false
					for _, vv := range vm {
						if n == i {
							found = true
							v = vv
							break
						}
						n++
					}
					if !found {
						return errors.New("invalid path: " + ss[2])
					}
				} else if vm, ok := v.(map[interface{}]interface{}); ok {
					n, found := 0, false
					for _, vv := range vm {
						if n == i {
							found = true
							v = vv
							break
						}
						n++
					}
					if !found {
						return errors.New("invalid path: " + ss[2])
					}
				} else {
					return errors.New("invalid path: " + ss[2])
				}
			} else {
				if i < 0 || i > len(vl)-1 {
					return errors.New("invalid path: " + ss[2])
				}
				v = vl[i]
			}
		}
	}
	return Scan(v, t)
}

// ScanJSON work as same sa ScanTree. it allow to give Reader.
func ScanJSON(r io.Reader, p string, t interface{}) (err error) {
	defer func() {
		if err == nil {
			err = toError(recover())
		}
	}()
	var a interface{}
	if err = json.NewDecoder(r).Decode(&a); err != nil {
		return err
	}
	return ScanTree(a, p, t)
}
