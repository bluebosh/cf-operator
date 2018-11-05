package cfdeployment_test

import (
	"testing"

	"code.cloudfoundry.org/cf-operator/pkg/apis"
	fissile "code.cloudfoundry.org/cf-operator/pkg/apis/fissile/v1alpha1"
	cfd "code.cloudfoundry.org/cf-operator/pkg/controller/cfdeployment"
	cfakes "code.cloudfoundry.org/cf-operator/pkg/controller/cfdeployment/fakes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func init() {
	// Add types to scheme: https://github.com/kubernetes-sigs/controller-runtime/issues/137
	apis.AddToScheme(scheme.Scheme)
}

func TestReconcile(t *testing.T) {
	assert := assert.New(t)

	c := fake.NewFakeClient(
		&fissile.CFDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			Spec: fissile.CFDeploymentSpec{ManifestRef: "config"},
		},
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "config",
				Namespace: "default",
			},
			Data: map[string]string{"manifest": "---"},
		},
	)

	m := &cfakes.FakeManager{}
	m.GetClientReturns(c)
	r := cfd.NewReconciler(m)

	request := reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}
	_, err := r.Reconcile(request)

	assert.Error(err)
	assert.Contains(err.Error(), "manifest is missing instance groups")
}
