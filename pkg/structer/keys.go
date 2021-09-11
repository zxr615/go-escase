package structer

import (
	"reflect"
)

func Keys(out interface{}) []string {
	v := reflect.ValueOf(out).Elem()

	ks := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Tag.Get("json")
		ks[i] = key
	}

	return ks
}
