package kubegen

import (
	kube "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	rbac "k8s.io/api/rbac/v1"
)

func ServiceAccount(name string) *kube.ServiceAccount {
	return &kube.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func Role(name string, rules ...rbac.PolicyRule) *rbac.Role {
	return &rbac.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: rules,
	}
}

func ClusterRole(name string, rules ...rbac.PolicyRule) *rbac.ClusterRole {
	return &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: rules,
	}
}

func RoleBindingForServiceAccount(name string, role string) *rbac.RoleBinding {
	return &rbac.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Subjects: []rbac.Subject{{
			Kind: rbac.ServiceAccountKind,
			Name: name,
		}},
		RoleRef: rbac.RoleRef{
			Kind: "Role",
			Name: role,
		},
	}
}

func ClusterRoleBindingForServiceAccount(namespace string, name string, role string) *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace + "-" + name,
		},
		Subjects: []rbac.Subject{{
			Kind:      rbac.ServiceAccountKind,
			Namespace: namespace,
			Name:      name,
		}},
		RoleRef: rbac.RoleRef{
			Kind: "ClusterRole",
			Name: role,
		},
	}
}
