package kg

import (
	"io/ioutil"
	"log"
	"reflect"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func IntPtr(v int32) *int32 {
	return &v
}

func Int32Ptr(v int32) *int32 {
	return &v
}

func Int64Ptr(v int64) *int64 {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}

func IntstrPtr(v intstr.IntOrString) *intstr.IntOrString {
	return &v
}

func ReadString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file %s: %v", filename, err)
	}
	return string(b)
}

func ReadYAML(filename string) map[string]interface{} {

}

// Merge merges the from map into into
func Merge(into map[string]interface{}, from map[string]interface{}) {
	mergeReflect(reflect.ValueOf(into), reflect.ValueOf(from))
}

func mergeReflect(into reflect.Value, from reflect.Value) {
	if into.Type().Kind() != reflect.Map || from.Type().Kind() != reflect.Map {
		return
	}

	if !into.IsValid() || !from.IsValid() || into.IsNil() || from.IsNil() {
		return
	}

	for _, k := range from.MapKeys() {
		fromV := from.MapIndex(k)
		intoV := into.MapIndex(k)
		if !intoV.IsValid() || intoV.IsNil() || intoV.Type().Kind() != reflect.Map || fromV.Type().Kind() != reflect.Map {
			// Overwrite if not map-to-map
			into.SetMapIndex(k, fromV)
			continue
		}
		// Map-to-map merge, so recurse
		mergeReflect(intoV, fromV)
	}
}
