package kg

import "reflect"

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
