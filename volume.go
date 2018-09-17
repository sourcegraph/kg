package kubegen

import (
	kube "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func FixedVolume(name string, size string) *kube.PersistentVolume {
	return &kube.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: kube.PersistentVolumeSpec{
			PersistentVolumeSource: kube.PersistentVolumeSource{GCEPersistentDisk: &kube.GCEPersistentDiskVolumeSource{
				PDName: name,
				FSType: "ext4",
			}},
			AccessModes: []kube.PersistentVolumeAccessMode{kube.ReadWriteOnce},
			Capacity: kube.ResourceList{
				kube.ResourceStorage: resource.MustParse(size),
			},
		},
	}
}

func VolumeClaim(name string, size string, volumeName string) *kube.PersistentVolumeClaim {
	return &kube.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: kube.PersistentVolumeClaimSpec{
			AccessModes: []kube.PersistentVolumeAccessMode{kube.ReadWriteOnce},
			Resources: kube.ResourceRequirements{
				Requests: map[kube.ResourceName]resource.Quantity{
					kube.ResourceStorage: resource.MustParse(size),
				},
			},
			VolumeName: volumeName,
		},
	}
}
