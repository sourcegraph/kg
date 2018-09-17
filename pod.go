package kg

import kube "k8s.io/api/core/v1"

func PodSpec(opts ...PodSpecOpt) *kube.PodSpec {
	pod := &kube.PodSpec{}
	for _, opt := range opts {
		opt(pod)
	}
	return pod
}

type PodSpecOpt func(pod *kube.PodSpec)

func SecurityContext(podSecurityContext *kube.PodSecurityContext) PodSpecOpt {
	return func(pod *kube.PodSpec) {
		if podSecurityContext != nil {
			pod.SecurityContext = podSecurityContext
		}
	}
}

func NodeSelector(nodeSelector map[string]string) PodSpecOpt {
	return func(pod *kube.PodSpec) {
		if pod.NodeSelector == nil {
			pod.NodeSelector = make(map[string]string)
		}
		for k, v := range nodeSelector {
			pod.NodeSelector[k] = v
		}
	}
}

func Container(name string, opts ...ContainerOpt) PodSpecOpt {
	return func(pod *kube.PodSpec) {
		var containers []*kube.Container
		if name == "" {
			for i := range pod.Containers {
				containers = append(containers, &pod.Containers[i])
			}
		} else {
			for i, c := range pod.Containers {
				if c.Name == name {
					containers = []*kube.Container{&pod.Containers[i]}
					break
				}
			}
			if len(containers) == 0 {
				pod.Containers = append(pod.Containers, kube.Container{
					Name: name,
				})
				containers = []*kube.Container{&pod.Containers[len(pod.Containers)-1]}
			}
		}

		for _, container := range containers {
			for _, opt := range opts {
				opt(pod, container)
			}
		}
	}
}

func Volume(name string, source kube.VolumeSource) PodSpecOpt {
	return func(pod *kube.PodSpec) {
		pod.Volumes = append(pod.Volumes, kube.Volume{
			Name:         name,
			VolumeSource: source,
		})
	}
}

func ClaimedVolume(name string, claimName string) PodSpecOpt {
	return Volume(name, kube.VolumeSource{
		PersistentVolumeClaim: &kube.PersistentVolumeClaimVolumeSource{
			ClaimName: claimName,
		},
	})
}

func SecretVolume(name string, secretName string) PodSpecOpt {
	return Volume(name, kube.VolumeSource{
		Secret: &kube.SecretVolumeSource{
			SecretName: secretName,
		},
	})
}

func TerminationGracePeriod(seconds int64) PodSpecOpt {
	return func(pod *kube.PodSpec) {
		pod.TerminationGracePeriodSeconds = Int64Ptr(seconds)
	}
}

func UseServiceAccount(name string) PodSpecOpt {
	return func(pod *kube.PodSpec) {
		pod.ServiceAccountName = name
	}
}
