package kg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	// files is a map from filename, relative to the cluster root directory, to the Kubernetes
	// object deserialized from that file
	files map[string]runtime.Object

	// newFilesDir is the directory to which to add new files created by modifications
	newFilesDir string
}

func (c *Cluster) Write() error {
	for file, _ := range c.files {
		if strings.HasPrefix(file, strings.TrimSuffix(c.newFilesDir, string(filepath.Separator))+string(filepath.Separator)) {
			if err := os.MkdirAll(c.newFilesDir, 0777); err != nil {
				return err
			}
			break
		}
	}

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

func (c *Cluster) Deployments(names ...string) (selected Deployments) {
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

func (c *Cluster) StatefulSets(names ...string) (selected StatefulSets) {
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

func (c *Cluster) PersistentVolumeClaims(names ...string) (selected PersistentVolumeClaims) {
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

func (c *Cluster) Secrets(names ...string) (selected Secrets) {
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

	// Create non-existent secrets
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
	return selected
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
