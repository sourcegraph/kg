package kg

import (
	kube "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretOpt func(s *kube.Secret)

func Secret(name string) *kube.Secret {
	return &kube.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func SecretData(d map[string]string) SecretOpt {
	return func(s *kube.Secret) {
		if s.StringData != nil {
			s.StringData = map[string]string{}
		}
		for k, v := range d {
			s.StringData[k] = v
		}
	}
}

func SecretType(t kube.SecretType) SecretOpt {
	return func(s *kube.Secret) {
		s.Type = t
	}
}

func SecretMeta(meta metav1.ObjectMeta) SecretOpt {
	return func(s *kube.Secret) {
		s.ObjectMeta = meta
	}
}
