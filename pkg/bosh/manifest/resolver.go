package manifest

import (
	"context"
	"fmt"

	fissile "code.cloudfoundry.org/cf-operator/pkg/apis/fissile/v1alpha1"
	ipl "code.cloudfoundry.org/cf-operator/pkg/bosh/manifest/interpolator"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Resolver resolves references from CRD to a BOSH manifest
type Resolver interface {
	ResolveCRD(fissile.BOSHDeploymentSpec, string) (*Manifest, error)
}

// ResolverImpl implements Resolver interface
type ResolverImpl struct {
	client       client.Client
	interpolator ipl.Interpolator
}

// NewResolver constructs a resolver
func NewResolver(client client.Client, interpolator ipl.Interpolator) *ResolverImpl {
	return &ResolverImpl{client: client, interpolator: interpolator}
}

// ResolveCRD returns manifest referenced by our CRD
func (r *ResolverImpl) ResolveCRD(spec fissile.BOSHDeploymentSpec, namespace string) (*Manifest, error) {
	manifest := &Manifest{}

	// TODO for now we only support config map ref
	ref := spec.ManifestRef

	config := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: ref, Namespace: namespace}, config)
	if err != nil {
		return manifest, errors.Wrapf(err, "Failed to retrieve configmap '%s/%s' via client.Get", namespace, ref)
	}

	// unmarshal manifest.data into bosh deployment manifest...
	// TODO re-use LoadManifest() from fissile
	m, ok := config.Data["manifest"]
	if !ok {
		return manifest, fmt.Errorf("configmap doesn't contain manifest key")
	}

	// unmarshal ops.data into bosh ops if exist
	//opsConfig := &corev1.ConfigMap{}
	opsSecret := &corev1.Secret{}
	opsRef := spec.OpsRef
	if len(opsRef) == 0 {
		err = yaml.Unmarshal([]byte(m), manifest)
		return manifest, err
	}

	//err = r.client.Get(context.TODO(), types.NamespacedName{Name: opsRef, Namespace: namespace}, opsConfig)
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: opsRef[1]["secretRef"], Namespace: namespace}, opsSecret)
	if err != nil {
		return manifest, errors.Wrapf(err, "Failed to retrieve configmap '%s/%s' via client.Get", namespace, opsRef)
	}

	//opsData, ok := opsConfig.Data["ops"]
	encodedData, ok := opsSecret.Data["ops"]
	if !ok {
		return manifest, fmt.Errorf("configmap doesn't contain ops key")
	}
	opsData := fmt.Sprintf("%s", encodedData)
	err = r.interpolator.BuildOps([]byte(opsData))
	if err != nil {
		return manifest, errors.Wrapf(err, "Failed to build ops: %#v", opsData)
	}

	bytes, err := r.interpolator.Interpolate([]byte(m))
	if err != nil {
		return manifest, errors.Wrapf(err, "Failed to interpolate %#v by %#v", m, opsData)
	}

	err = yaml.Unmarshal(bytes, manifest)

	return manifest, err
}
