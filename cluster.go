package kg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/api/apps/v1"
	kube "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func ModifyCluster(rootDir, newFilesDir string, apply func(*Cluster)) error {
	var yamlFiles []string
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ext != ".yaml" {
			return nil
		}
		yamlFiles = append(yamlFiles, path)
		return nil
	})

	c, err := NewCluster(yamlFiles, newFilesDir)
	if err != nil {
		return err
	}

	apply(c)

	return c.Write()
}

func NewCluster(files []string, newFilesDir string) (*Cluster, error) {
	c := &Cluster{files: make(map[string]runtime.Object), newFilesDir: newFilesDir}
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(b, nil, nil)
		if err != nil {
			return nil, err
		}

		c.files[file] = obj
	}

	return c, nil
}

type Cluster struct {
	files map[string]runtime.Object

	// newFilesDir is the directory to which to add new files created by modifications
	newFilesDir string
}

func (c *Cluster) Write() error {
	for file, obj := range c.files {
		e := json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil)

		var buf bytes.Buffer
		err := e.Encode(obj, &buf)
		if err != nil {
			return err
		}

		var untyped map[string]interface{}
		yaml.Unmarshal(buf.Bytes(), &untyped)
		sanitize(untyped)

		sanitized, err := yaml.Marshal(untyped)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(file, sanitized, 0666); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cluster) Deployment(name string) *v1.Deployment {
	for _, obj := range c.files {
		if deploy, ok := obj.(*v1.Deployment); ok {
			if deploy.ObjectMeta.Name == name {
				return deploy
			}
		}
	}
	return nil
}

func (c *Cluster) Deployments(names ...string) (selected []*v1.Deployment) {
	selectAll := false
	nameSet := make(map[string]struct{})
	for _, name := range names {
		if name == "*" {
			selectAll = true
		} else {
			nameSet[name] = struct{}{}
		}
	}
	for _, obj := range c.files {
		if deploy, ok := obj.(*v1.Deployment); ok {
			_, exists := nameSet[deploy.ObjectMeta.Name]
			if selectAll || exists {
				selected = append(selected, deploy)
			}
		}
	}
	return selected
}

func (c *Cluster) StatefulSets(names ...string) (selected []*v1.StatefulSet) {
	selectAll := false
	nameSet := make(map[string]struct{})
	for _, name := range names {
		if name == "*" {
			selectAll = true
		} else {
			nameSet[name] = struct{}{}
		}
	}
	for _, obj := range c.files {
		if sset, ok := obj.(*v1.StatefulSet); ok {
			_, exists := nameSet[sset.ObjectMeta.Name]
			if selectAll || exists {
				selected = append(selected, sset)
			}
		}
	}
	return selected
}

func (c *Cluster) StatefulSet(name string) *v1.StatefulSet {
	for _, obj := range c.files {
		if sset, ok := obj.(*v1.StatefulSet); ok {
			if sset.ObjectMeta.Name == name {
				return sset
			}
		}
	}
	return nil
}

func (c *Cluster) PVC(names ...string) (selected []*kube.PersistentVolumeClaim) {
	selectAll := false
	nameSet := make(map[string]struct{})
	for _, name := range names {
		if name == "*" {
			selectAll = true
		} else {
			nameSet[name] = struct{}{}
		}
	}
	for _, obj := range c.files {
		if pvc, ok := obj.(*kube.PersistentVolumeClaim); ok {
			_, exists := nameSet[pvc.ObjectMeta.Name]
			if selectAll || exists {
				selected = append(selected, pvc)
			}
		}
	}
	return selected
}

func (c *Cluster) Secrets(createIfNotExists bool, names ...string) (selected []*kube.Secret) {
	selectAll := false
	nameSet := make(map[string]bool)
	for _, name := range names {
		if name == "*" {
			selectAll = true
		} else {
			nameSet[name] = false
		}
	}
	for _, obj := range c.files {
		if secret, ok := obj.(*kube.Secret); ok {
			_, exists := nameSet[secret.ObjectMeta.Name]
			if exists {
				nameSet[secret.ObjectMeta.Name] = true
				if selectAll {
					selected = append(selected, secret)
				}
			}
		}
	}
	if createIfNotExists {
		for name, found := range nameSet {
			if found {
				continue
			}
			newFile := filepath.Join(c.newFilesDir, fmt.Sprintf("%s.Secret.yaml", name))
			if _, exists := c.files[newFile]; exists {
				log.Fatalf("new file %s would conflict with existing file", newFile)
			}
			newSecret := Secret(name)
			c.files[newFile] = newSecret
			selected = append(selected, newSecret)
		}
	}
	return selected
}

func (c *Cluster) ModifyPVC(names []string, opts ...PersistentVolumeClaimOpt) {
	for _, pvc := range c.PVC(names...) {
		for _, opt := range opts {
			opt(pvc)
		}
	}
}

func (c *Cluster) ModifyDeployments(names []string, opts ...DeploymentOpt) {
	for _, deploy := range c.Deployments(names...) {
		for _, opt := range opts {
			opt(deploy)
		}
	}
}

func (c *Cluster) ModifyStatefulSets(names []string, opts ...StatefulSetOpt) {
	for _, sset := range c.StatefulSets(names...) {
		for _, opt := range opts {
			opt(sset)
		}
	}
}

func (c *Cluster) ModifySecrets(names []string, opts ...SecretOpt) {
	for _, s := range c.Secrets(true, names...) {
		for _, opt := range opts {
			opt(s)
		}
	}
}

// sanitize removes fields that shouldn't be present in the persisted YAML files but are emitted by
// the k8s config serializer.
func sanitize(m interface{}) {
	switch m := m.(type) {
	case map[interface{}]interface{}:
		delete(m, "status")
		delete(m, "creationTimestamp")
		delete(m, "dataSource")
		if r, ok := m["resources"].(map[interface{}]interface{}); ok && len(r) == 0 {
			delete(m, "resources")
		}
		for _, v := range m {
			sanitize(v)
		}
	case map[string]interface{}:
		delete(m, "status")
		delete(m, "creationTimestamp")
		delete(m, "dataSource")
		if r, ok := m["resources"].(map[interface{}]interface{}); ok && len(r) == 0 {
			delete(m, "resources")
		}
		for _, v := range m {
			sanitize(v)
		}
	case []interface{}:
		for _, v := range m {
			sanitize(v)
		}
	}
}
