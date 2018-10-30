package deployment_manifest

import (
	"context"

	fissile "code.cloudfoundry.org/cf-operator/pkg/apis/fissile/v1alpha1"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ResolveCRD(spec fissile.CFDeploymentSpec, client client.Client, namespace string) (*Manifest, error) {
	manifest := &Manifest{}

	// for now we only support config map ref
	ref := spec.ManifestRef

	config := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: ref, Namespace: namespace}, config)
	if err != nil {
		return manifest, err
	}
	// unmarshal manifet.data into bosh deployment manifest...
	// TODO LoadManifest() from fisisle
	m, ok := config.Data["manifest"]
	if !ok {
		return manifest, nil
	}
	err = yaml.Unmarshal([]byte(m), manifest)
	return manifest, nil
}
