package kg

import (
	"sort"

	kube "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ContainerOp func(pod *kube.PodSpec, container *kube.Container)

func Command(command ...string) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Command = command
	}
}

func Args(args ...string) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Args = args
	}
}

func Env(vars map[string]string, addIfNotExist bool) ContainerOp {
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

type byName []kube.EnvVar

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func ResourceRequests(cpu, memory string) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Resources.Requests = kube.ResourceList{
			kube.ResourceCPU:    resource.MustParse(cpu),
			kube.ResourceMemory: resource.MustParse(memory),
		}
	}
}

func ResourceLimits(cpu, memory string) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.Resources.Limits = kube.ResourceList{
			kube.ResourceCPU:    resource.MustParse(cpu),
			kube.ResourceMemory: resource.MustParse(memory),
		}
	}
}

func VolumeMount(name string, mountPath string, ops ...VolumeMountOp) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		var mount *kube.VolumeMount
		for i := range container.VolumeMounts {
			if container.VolumeMounts[i].Name == name {
				mount = &container.VolumeMounts[i]
				break
			}
		}
		if mount == nil {
			container.VolumeMounts = append(container.VolumeMounts, kube.VolumeMount{
				Name:      name,
				MountPath: mountPath,
			})
			mount = &container.VolumeMounts[len(container.VolumeMounts)-1]
		}

		for _, op := range ops {
			op(mount)
		}
	}
}

func ReadinessProbe(p *kube.Probe) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.ReadinessProbe = p
	}
}

func LivenessProbe(p *kube.Probe) ContainerOp {
	return func(pod *kube.PodSpec, container *kube.Container) {
		container.LivenessProbe = p
	}
}

type VolumeMountOp func(mount *kube.VolumeMount)

func ReadOnly() VolumeMountOp {
	return func(mount *kube.VolumeMount) {
		mount.ReadOnly = true
	}
}
