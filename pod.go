package kg

import kube "k8s.io/api/core/v1"

func PodSpec(ops ...PodSpecOp) *kube.PodSpec {
	pod := &kube.PodSpec{}
	for _, op := range ops {
		op(pod)
	}
	return pod
}

type PodSpecOp func(pod *kube.PodSpec)

func SecurityContext(podSecurityContext *kube.PodSecurityContext) PodSpecOp {
	return func(pod *kube.PodSpec) {
		pod.SecurityContext = podSecurityContext
	}
}

func NodeSelector(nodeSelector map[string]string) PodSpecOp {
	return func(pod *kube.PodSpec) {
		if pod.NodeSelector == nil {
			pod.NodeSelector = make(map[string]string)
		}
		for k, v := range nodeSelector {
			pod.NodeSelector[k] = v
		}
	}
}

func Container(name string, ops ...ContainerOp) PodSpecOp {
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
			for _, op := range ops {
				op(pod, container)
			}
		}
	}
}

func Volume(name string, source kube.VolumeSource) PodSpecOp {
	return func(pod *kube.PodSpec) {
		found := false
		for i := range pod.Volumes {
			if pod.Volumes[i].Name == name {
				pod.Volumes[i].VolumeSource = source
				found = true
			}
		}

		if !found {
			pod.Volumes = append(pod.Volumes, kube.Volume{
				Name:         name,
				VolumeSource: source,
			})
		}
	}
}

func ClaimedVolume(name string, claimName string) PodSpecOp {
	return Volume(name, kube.VolumeSource{
		PersistentVolumeClaim: &kube.PersistentVolumeClaimVolumeSource{
			ClaimName: claimName,
		},
	})
}

func SecretVolume(name string, secretName string) PodSpecOp {
	return Volume(name, kube.VolumeSource{
		Secret: &kube.SecretVolumeSource{
			SecretName: secretName,
		},
	})
}

func TerminationGracePeriod(seconds int64) PodSpecOp {
	return func(pod *kube.PodSpec) {
		pod.TerminationGracePeriodSeconds = Int64Ptr(seconds)
	}
}

func UseServiceAccount(name string) PodSpecOp {
	return func(pod *kube.PodSpec) {
		pod.ServiceAccountName = name
	}
}
