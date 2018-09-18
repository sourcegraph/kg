package kg

import (
	"strconv"

	kube "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Service(name string, ops ...ServiceOp) *kube.Service {
	return ServiceForApp(name, name, ops...)
}

func ServiceForApp(name string, app string, ops ...ServiceOp) *kube.Service {
	svc := &kube.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": app,
			},
			Annotations: make(map[string]string),
		},
		Spec: kube.ServiceSpec{
			Type: kube.ServiceTypeClusterIP,
			Selector: map[string]string{
				"app": app,
			},
		},
	}

	for _, op := range ops {
		op(svc)
	}

	return svc
}

type ServiceOp func(*kube.Service)

func ServicePort(name string, port int32) ServiceOp {
	return func(svc *kube.Service) {
		svc.Spec.Ports = append(svc.Spec.Ports, kube.ServicePort{
			Name:       name,
			Port:       port,
			TargetPort: intstr.FromString(name),
		})
	}
}

func NodePort(name string, port int32) ServiceOp {
	return func(svc *kube.Service) {
		svc.Spec.Type = kube.ServiceTypeNodePort
		svc.Spec.Ports = append(svc.Spec.Ports, kube.ServicePort{
			Name:       name,
			Port:       port,
			NodePort:   port,
			TargetPort: intstr.FromString(name),
		})
	}
}

func MetricsPort(port int32) ServiceOp {
	return func(svc *kube.Service) {
		svc.ObjectMeta.Annotations["prometheus.io/scrape"] = "true"
		svc.ObjectMeta.Annotations["prometheus.io/port"] = strconv.Itoa(int(port))
	}
}

func MetricsPortWithPath(port int32, path string) ServiceOp {
	return func(svc *kube.Service) {
		svc.ObjectMeta.Annotations["prometheus.io/scrape"] = "true"
		svc.ObjectMeta.Annotations["prometheus.io/path"] = path
		svc.ObjectMeta.Annotations["prometheus.io/port"] = strconv.Itoa(int(port))
	}
}

func PublicIP(ip string) ServiceOp {
	return func(svc *kube.Service) {
		svc.Spec.Type = kube.ServiceTypeLoadBalancer
		svc.Spec.LoadBalancerIP = ip
	}
}

func Headless() ServiceOp {
	return func(svc *kube.Service) {
		svc.Spec.ClusterIP = "None"
		svc.Spec.Ports = append(svc.Spec.Ports, kube.ServicePort{
			Name:       "unused",
			Port:       10811,
			TargetPort: intstr.FromInt(10811),
		})
	}
}

func Selector(labels map[string]string) ServiceOp {
	return func(svc *kube.Service) {
		svc.ObjectMeta.Labels = labels
		svc.Spec.Selector = labels
	}
}

// GroupedService creates a headless service on pods with the label
// "group"=group. This is used to hint at the scheduler to not place them on
// the same nodes.
func GroupedService(group string) *kube.Service {
	return Service(group, Headless(), Selector(map[string]string{"group": group}))
}
