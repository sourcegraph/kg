package kg

import kubeext "k8s.io/api/apps/v1"

type DeploymentOp func(depl *kubeext.Deployment)

func Replicas(count int32) DeploymentOp {
	return func(depl *kubeext.Deployment) {
		depl.Spec.Replicas = IntPtr(count)
	}
}

func Pod(podOps ...PodSpecOp) DeploymentOp {
	return func(depl *kubeext.Deployment) {
		for _, podOp := range podOps {
			podOp(&depl.Spec.Template.Spec)
		}
	}
}
