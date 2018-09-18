package kg

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeext "k8s.io/api/apps/v1"
	kube "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Deployment creates a deployment for 1 replica. If the podSpec has no
// volumes (ie stateless), the rollout strategy will prevent
// downtime. Otherwise minimal downtime will occur during rollout to allow k8s
// to remount the volume on the new node.
func Deployment(name string, description string, podSpec *kube.PodSpec, ops ...DeploymentOp) *kubeext.Deployment {
	maxUnavailable := 1
	hasRealVolume := false
	for _, vol := range podSpec.Volumes {
		if vol.VolumeSource.ConfigMap == nil {
			hasRealVolume = true
			break
		}
	}
	if !hasRealVolume {
		// This pod is stateless, so we can use a zero downtime
		// rollout strategy
		maxUnavailable = 0
	}

	depl := &kubeext.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"description": description,
			},
		},
		Spec: kubeext.DeploymentSpec{
			Replicas:             IntPtr(1),
			RevisionHistoryLimit: IntPtr(10),
			MinReadySeconds:      10,
			Strategy: kubeext.DeploymentStrategy{
				Type: kubeext.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &kubeext.RollingUpdateDeployment{
					MaxUnavailable: IntstrPtr(intstr.FromInt(maxUnavailable)),
					MaxSurge:       IntstrPtr(intstr.FromInt(1)),
				},
			},
			Template: kube.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: *podSpec,
			},
		},
	}

	for _, op := range ops {
		op(depl)
	}

	return depl
}

// GroupedDeployment is a Deployment with an extra pod label called
// group. The intended use of this is to hint to the scheduler not to schedule
// pods with the same group on the same node. This is done by adding a
// headless service per group.
func GroupedDeployment(group, name string, description string, podSpec *kube.PodSpec) *kubeext.Deployment {
	d := Deployment(name, description, podSpec)
	d.Spec.Template.Labels["group"] = group
	return d
}
