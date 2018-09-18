package kg

import (
	"sort"

	kube "k8s.io/api/core/v1"
)

func EnvVarFrom(name string, source *kube.EnvVarSource) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Env = append(container.Env, kube.EnvVar{Name: name, ValueFrom: source})
		sort.Sort(byName(container.Env))
	}
}

func EnvVarFromSecret(varName string, secretName string, secretKey string, optional bool) ContainerOp {
	return EnvVarFrom(varName, &kube.EnvVarSource{
		SecretKeyRef: &kube.SecretKeySelector{
			LocalObjectReference: kube.LocalObjectReference{Name: secretName},
			Key:                  secretKey,
			Optional:             BoolPtr(optional),
		},
	})
}

func EnvVarFromFieldSelector(name string, fieldPath string) ContainerOp {
	return EnvVarFrom(name, &kube.EnvVarSource{
		FieldRef: &kube.ObjectFieldSelector{FieldPath: fieldPath},
	})
}

func ContainerPort(name string, port int32) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Ports = append(container.Ports, kube.ContainerPort{
			Name:          name,
			ContainerPort: port,
		})
	}
}
