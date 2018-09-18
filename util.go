package kg

import (
	"io/ioutil"
	"log"

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

func ReadFile(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file %s: %v", filename, err)
	}
	return string(b)
}
