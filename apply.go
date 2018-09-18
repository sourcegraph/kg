package kg

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

type Deployments []*apps.Deployment

func (s Deployments) Apply(ops ...DeploymentOpt) {
	for _, c := range s {
		for _, op := range ops {
			op(c)
		}
	}
}

type PersistentVolumeClaims []*core.PersistentVolumeClaim

func (s PersistentVolumeClaims) Apply(ops ...PersistentVolumeClaimOpt) {
	for _, c := range s {
		for _, op := range ops {
			op(c)
		}
	}
}

type StatefulSets []*apps.StatefulSet

func (s StatefulSets) Apply(ops ...StatefulSetOpt) {
	for _, c := range s {
		for _, op := range ops {
			op(c)
		}
	}
}

type Secrets []*core.Secret

func (s Secrets) Apply(ops ...SecretOpt) {
	for _, c := range s {
		for _, op := range ops {
			op(c)
		}
	}
}
