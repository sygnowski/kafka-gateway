package util

import (
	"fmt"
	"os"
)

//not valid for windows os
func UUID() (uuid string) {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

func MatchNestedMapPath(c interface{}, paths []string) bool {
	var context interface{}
	context = c.(map[string]interface{})

	for _, path := range paths {
		if m, isMap := context.(map[string]interface{}); isMap {
			if val, ok := m[path]; ok {
				context = val
			} else {
				return false
			}
		}
	}
	return true
}

func GetMapPath(c interface{}, paths []string) interface{} {
	var context interface{}
	context = c.(map[string]interface{})

	for _, path := range paths {
		if m, isMap := context.(map[string]interface{}); isMap {
			if val, ok := m[path]; ok {
				context = val
			} else {
				return nil
			}
		}
	}
	return context
}

func GetLastMapPath(c interface{}, paths []string) (fullpath bool, ctx interface{}) {
	var context interface{}
	context = c.(map[string]interface{})
	deepEnd := len(paths) - 1

	fullpath = false
	for d, path := range paths {
		if m, isMap := context.(map[string]interface{}); isMap {
			if val, ok := m[path]; ok {
				context = val
				if d == deepEnd {
					fullpath = true
				}
			} else {
				break
			}
		}
	}
	ctx = context
	return
}
