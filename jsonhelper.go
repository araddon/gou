package gou

import (
	"strconv"
	"strings"
)

// A simple wrapper tohelp json config files be more easily used
type JsonHelper map[string]interface{}

func (j JsonHelper) Helper(n string) JsonHelper {
	if v, ok := j[n]; ok {
		switch v.(type) {
		case map[string]interface{}:
			cn := JsonHelper{}
			for n, val := range v.(map[string]interface{}) {
				cn[n] = val
			}
			return cn
		case map[string]string:
			cn := JsonHelper{}
			for n, val := range v.(map[string]string) {
				cn[n] = val
			}
			return cn
		default:
			Debug("no map? ", v)
		}
	}
	return nil
}
func jsonList(v interface{}) []interface{} {
	switch v.(type) {
	case []interface{}:
		return v.([]interface{})
	}
	return nil
}
func jsonEntry(name string, v interface{}) (interface{}, bool) {
	switch v.(type) {
	case map[string]interface{}:
		if root, ok := v.(map[string]interface{})[name]; ok {
			return root, true
		} else {
			return nil, false
		}
	case []interface{}:
		return v, true
	default:
		Debug("no type? ", name, " ", v)
		return nil, false
	}
	return nil, false
}
func (j JsonHelper) Get(n string) interface{} {
	parts := strings.Split(n, ".")
	var root interface{}
	var err error
	var ok, isList, listEntry bool
	var ln, st, idx int
	for ict, name := range parts {
		isList = strings.HasSuffix(name, "[]")
		listEntry = strings.HasSuffix(name, "]") && !isList
		ln, idx = len(name), -1
		if isList || listEntry {
			st = strings.Index(name, "[")
			idx, err = strconv.Atoi(name[st+1 : ln-1])
			name = name[:st]
		}
		if ict == 0 {
			root, ok = j[name]
		} else {
			root, ok = jsonEntry(name, root)
		}
		//Debug(isList, listEntry, " ", name, " ", root, " ", ok, err)
		if !ok {
			return nil
		}
		if isList {
			return jsonList(root)
		} else if listEntry && err == nil {
			if lst := jsonList(root); lst != nil && len(lst) > idx {
				root = lst[idx]
			} else {
				return nil
			}
		}

	}
	return root
}

func (j JsonHelper) Int64(n string) int64 {
	v := j.Get(n)
	if v != nil {
		switch v.(type) {
		case int:
			return int64(v.(int))
		case int64:
			return int64(v.(int64))
		case uint32:
			return int64(v.(uint32))
		case uint64:
			return int64(v.(uint64))
		case float32:
			f := float64(v.(float32))
			return int64(f)
		case float64:
			f := v.(float64)
			return int64(f)
		default:
			Debug("no type? ", n, " ", v)
		}
	}
	return -1
}
func (j JsonHelper) String(n string) string {
	v := j.Get(n)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
func (j JsonHelper) StringSafe(n string) (string, bool) {
	v := j.Get(n)
	if v != nil {
		if s, ok := v.(string); ok {
			return s, ok
		}
	}
	return "", false
}
func (j JsonHelper) Int(n string) int {
	v := j.Get(n)
	if v != nil {
		switch v.(type) {
		case int:
			return v.(int)
		case int64:
			return int(v.(int64))
		case uint32:
			return int(v.(uint32))
		case uint64:
			return int(v.(uint64))
		case float32:
			f := float64(v.(float32))
			return int(f)
		case float64:
			f := v.(float64)
			return int(f)
		default:
			Debug("no type? ", n, " ", v)
		}
	}
	return -1
}
func (j JsonHelper) Uint64(n string) uint64 {
	v := j.Get(n)
	if v != nil {
		switch v.(type) {
		case int:
			return uint64(v.(int))
		case int64:
			return uint64(v.(int64))
		case uint32:
			return uint64(v.(uint32))
		case uint64:
			return uint64(v.(uint64))
		case float32:
			f := float64(v.(float32))
			return uint64(f)
		case float64:
			f := v.(float64)
			return uint64(f)
		default:
			Debug("no type? ", n, " ", v)
		}
	}
	return 0
}
func (j JsonHelper) Bool(n string) bool {
	v := j.Get(n)
	if v != nil {
		switch v.(type) {
		case bool:
			return v.(bool)
		case string:
			if s := v.(string); len(s) > 0 {
				if b, err := strconv.ParseBool(s); err == nil {
					return b
				}
			}
		default:
			Debug("no type? ", n, " ", v)
		}
	}
	return false
}
