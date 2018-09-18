package kg

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

type Deployments []*apps.Deployment

func (s Deployments) Apply(ops ...DeploymentOp) {
	for _, c := range s {
		for _, op := range ops {
			op(c)
		}
	}
}

type PersistentVolumeClaims []*core.PersistentVolumeClaim

func (s PersistentVolumeClaims) Apply(ops ...PersistentVolumeClaimOp) {
	for _, c := range s {
		for _, op := range ops {
			op(c)
		}
	}
}

type StatefulSets []*apps.StatefulSet

func (s StatefulSets) Apply(ops ...StatefulSetOp) {
	for _, c := range s {
		for _, op := range ops {
			op(c)
		}
	}
}

type Secrets []*core.Secret

func (s Secrets) Apply(ops ...SecretOp) {
	for _, c := range s {
		for _, op := range ops {
			op(c)
		}
	}
}
