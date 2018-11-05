package deployment_manifest_test

import (
	"testing"

	fissile "code.cloudfoundry.org/cf-operator/pkg/apis/fissile/v1alpha1"
	dm "code.cloudfoundry.org/cf-operator/pkg/deployment_manifest"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	fake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestResolveCRD(t *testing.T) {
	assert := assert.New(t)

	config := &corev1.ConfigMap{}
	client := fake.NewFakeClient(config)
	spec := fissile.CFDeploymentSpec{}
	manifest, err := dm.ResolveCRD(spec, client, "default")

	assert.NoError(err)
	assert.NotNil(manifest)
}
