package structer

import (
	"reflect"
	"testing"
)

func Test_keys(t *testing.T) {
	type info struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	ks := Keys(new(info))
	want := []string{"name", "age"}

	if !reflect.DeepEqual(ks, want) {
		t.Errorf("keys() = %v, want %v", ks, want)
	}
}
