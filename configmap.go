package kg

import (
	core "k8s.io/api/core/v1"
)

type ConfigMapOp func(cm *core.ConfigMap)

func ConfigMapData(data map[string]string) ConfigMapOp {
	return func(cm *core.ConfigMap) {
		if cm.Data == nil {
			cm.Data = make(map[string]string)
		}
		for k, v := range data {
			cm.Data[k] = v
		}
	}
}
