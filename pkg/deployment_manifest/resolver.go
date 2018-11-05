package deployment_manifest

import (
	"context"
	"fmt"

	fissile "code.cloudfoundry.org/cf-operator/pkg/apis/fissile/v1alpha1"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ResolveCRD(spec fissile.CFDeploymentSpec, client client.Client, namespace string) (*Manifest, error) {
	manifest := &Manifest{}

	// TODO for now we only support config map ref
	ref := spec.ManifestRef

	config := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: ref, Namespace: namespace}, config)
	if err != nil {
		return manifest, errors.Wrap(err, "failed to retrieve via client.Get")
	}

	// unmarshal manifest.data into bosh deployment manifest...
	// TODO re-use LoadManifest() from fisisle
	m, ok := config.Data["manifest"]
	if !ok {
		return manifest, fmt.Errorf("configmap doesn't contain manifest key")
	}
	err = yaml.Unmarshal([]byte(m), manifest)

	return manifest, err
}
