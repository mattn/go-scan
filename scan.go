package scan

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile("^([^0-9\\s\\[][^\\s\\[]*)?(\\[[0-9]+\\])?$")

func Scan(v interface{}, t interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(string); ok {
				err = errors.New(e)
			} else if e, ok := r.(error); ok {
				err = e
			}
			err = errors.New("Unknown error")
		}
	}()
	rt := reflect.ValueOf(t).Elem()
	rv := reflect.ValueOf(v)
	tv := rv.Type().Kind()

	if tv == reflect.Slice || tv == reflect.Array {
		ia := rv.Interface().([]interface{})
		rt.Set(reflect.MakeSlice(rt.Type(), len(ia), len(ia)))
		for n, _ := range ia {
			rt.Index(n).Set(rv.Index(n).Elem())
		}
	} else {
		rt.Set(rv)
	}
	return nil
}

func ScanTree(v interface{}, p string, t interface{}) (err error) {
	if p == "" {
		return errors.New("invalid path")
	}
	var ok bool
	for _, token := range strings.Split(p, "/") {
		sl := re.FindAllStringSubmatch(token, -1)
		if len(sl) == 0 {
			return errors.New("invalid path")
		}
		ss := sl[0]
		if ss[1] != "" {
			if v, ok = v.(map[string]interface{})[ss[1]]; !ok {
				return errors.New("invalid path: " + ss[1])
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
					for _, vv := range vm {
						v = vv
						break
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

func ScanJSON(r io.Reader, p string, t interface{}) (err error) {
	var a interface{}
	if err = json.NewDecoder(r).Decode(&a); err != nil {
		return
	}
	return ScanTree(a, p, t)
}
