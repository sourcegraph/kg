package kg

import (
	"sort"

	kube "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ContainerOpt func(pod *kube.PodSpec, container *kube.Container)

func Command(command ...string) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Command = command
	}
}

func Args(args ...string) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Args = args
	}
}

func Env(vars map[string]string, addIfNotExist bool) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		for name, value := range vars {
			exists := false
			for i := range container.Env {
				if container.Env[i].Name == name {
					container.Env[i].Value = value
					exists = true
				}
			}

			if !exists && addIfNotExist {
				container.Env = append(container.Env, kube.EnvVar{Name: name, Value: value})
			}
		}
		sort.Sort(byName(container.Env))
	}
}

func EnvAllowEmpty(vars map[string]string) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		for name, value := range vars {
			container.Env = append(container.Env, kube.EnvVar{Name: name, Value: value})
		}
		sort.Sort(byName(container.Env))
	}
}

func EnvVarFrom(name string, source *kube.EnvVarSource) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Env = append(container.Env, kube.EnvVar{Name: name, ValueFrom: source})
		sort.Sort(byName(container.Env))
	}
}

func EnvVarFromSecret(varName string, secretName string, secretKey string, optional bool) ContainerOpt {
	return EnvVarFrom(varName, &kube.EnvVarSource{
		SecretKeyRef: &kube.SecretKeySelector{
			LocalObjectReference: kube.LocalObjectReference{Name: secretName},
			Key:                  secretKey,
			Optional:             BoolPtr(optional),
		},
	})
}

func EnvVarFromFieldSelector(name string, fieldPath string) ContainerOpt {
	return EnvVarFrom(name, &kube.EnvVarSource{
		FieldRef: &kube.ObjectFieldSelector{FieldPath: fieldPath},
	})
}

type byName []kube.EnvVar

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func ResourceRequests(cpu, memory string) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Resources.Requests = kube.ResourceList{
			kube.ResourceCPU:    resource.MustParse(cpu),
			kube.ResourceMemory: resource.MustParse(memory),
		}
	}
}

func ResourceLimits(cpu, memory string) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Resources.Limits = kube.ResourceList{
			kube.ResourceCPU:    resource.MustParse(cpu),
			kube.ResourceMemory: resource.MustParse(memory),
		}
	}
}

func ContainerPort(name string, port int32) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Ports = append(container.Ports, kube.ContainerPort{
			Name:          name,
			ContainerPort: port,
		})
	}
}

func VolumeMount(name string, mountPath string, opts ...VolumeMountOpt) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		mount := kube.VolumeMount{
			Name:      name,
			MountPath: mountPath,
		}
		for _, opt := range opts {
			opt(&mount)
		}
		container.VolumeMounts = append(container.VolumeMounts, mount)
	}
}

type VolumeMountOpt func(mount *kube.VolumeMount)

func ReadOnly() VolumeMountOpt {
	return func(mount *kube.VolumeMount) {
		mount.ReadOnly = true
	}
}

func ReadinessProbe(p *kube.Probe) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.ReadinessProbe = p
	}
}

func LivenessProbe(p *kube.Probe) ContainerOpt {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.LivenessProbe = p
	}
}
