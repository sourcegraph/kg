package kg

import (
	kubeext "k8s.io/api/apps/v1"
)

type StatefulSetOp func(sset *kubeext.StatefulSet)

func StatefulSetReplicas(count int32) StatefulSetOp {
	return func(sset *kubeext.StatefulSet) {
		sset.Spec.Replicas = IntPtr(count)
	}
}

func StatefulSetPod(podOps ...PodSpecOp) StatefulSetOp {
	return func(sset *kubeext.StatefulSet) {
		for _, op := range podOps {
			op(&sset.Spec.Template.Spec)
		}
	}
}
