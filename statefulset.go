package kg

import (
	kubeext "k8s.io/api/apps/v1"
)

type StatefulSetOpt func(sset *kubeext.StatefulSet)

func StatefulSetReplicas(count int32) StatefulSetOpt {
	return func(sset *kubeext.StatefulSet) {
		sset.Spec.Replicas = IntPtr(count)
	}
}

func StatefulSetPod(podOpts ...PodSpecOpt) StatefulSetOpt {
	return func(sset *kubeext.StatefulSet) {
		for _, opt := range podOpts {
			opt(&sset.Spec.Template.Spec)
		}
	}
}
