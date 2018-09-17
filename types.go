package kubegen

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Object interface {
	GetObjectKind() schema.ObjectKind
	metav1.Object
}
