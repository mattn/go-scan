package scan

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile("^([^0-9\\s\\[][^\\s\\[]*)?(\\[[0-9]+\\])?$")

func Scan(v interface{}, t interface{}) (err error) {
	rt := reflect.ValueOf(t).Elem()
	rv := reflect.ValueOf(v)
	tv := rv.Type().Kind()

	defer func() {
		if errstr := recover(); errstr != nil {
			err = errors.New(errstr.(string))
		}
	}()
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

func ScanTree(v interface{}, jp string, t interface{}) (err error) {
	if jp == "" {
		return errors.New("invalid path")
	}
	var ok bool
	for _, token := range strings.Split(jp, "/") {
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
