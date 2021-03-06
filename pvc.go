package kg

import (
	kube "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type PersistentVolumeClaimOp func(pvc *kube.PersistentVolumeClaim)

func DiskSize(size string) PersistentVolumeClaimOp {
	return func(pvc *kube.PersistentVolumeClaim) {
		pvc.Spec.Resources.Requests[kube.ResourceStorage] = resource.MustParse(size)
	}
}
