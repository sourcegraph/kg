package kg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	yaml "gopkg.in/yaml.v2"

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
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file %s: %v", filename, err)
	}
	return string(b)
}

func ReadYAML(filename string) (v map[string]interface{}) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open file %s: %v", filename, err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&v); err != nil {
		log.Fatalf("Could not unmarshal YAML from file %s: %v", filename, err)
	}
	return v
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

// JSON marshals JSON using json.Marshal, but can hand internal objects of type
// map[interface{}]interface{}. Panics if an error occurs.
func JSON(m map[string]interface{}) string {
	normalize(m)
	s, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	return string(s)
}

// normalize converts all map[interface{}]interface{} to map[string]interface{}. Its parameter can
// be of any type *except* map[interface{}]interface{}. It mutates m.
func normalize(m interface{}) interface{} {
	switch m := m.(type) {
	case map[string]interface{}:
		for k, v := range m {
			switch v := v.(type) {
			case map[interface{}]interface{}:
				m[k] = normalize(convertMap(v))
			default:
				m[k] = normalize(v)
			}
		}
	case []interface{}:
		for i, v := range m {
			switch v := v.(type) {
			case map[interface{}]interface{}:
				m[i] = normalize(convertMap(v))
			default:
				m[i] = normalize(v)
			}
		}
	case map[interface{}]interface{}:
		log.Fatal("unreachable")
	}
	return m
}

func convertMap(m map[interface{}]interface{}) map[string]interface{} {
	m2 := make(map[string]interface{})
	for k, v := range m {
		m2[fmt.Sprintf("%s", k)] = v
	}
	return m2
}
