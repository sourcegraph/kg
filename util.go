package kg

import (
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
