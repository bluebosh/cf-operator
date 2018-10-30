package cfdeployment_test

import (
	"log"
	"testing"

	cfd "code.cloudfoundry.org/cf-operator/pkg/controller/cfdeployment"
	"k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/stretchr/testify/assert"
)

func TestReconcile(t *testing.T) {
	assert := assert.New(t)

	// TODO newReconciler?
	r := &cfd.ReconcileCFDeployment{client: fake.NewSimpleClientset(), scheme: mgr.GetScheme()}

	log.Printf("Reconciling CFDeployment %s/%s\n", request.Namespace, request.Name)
	request := reconcile.Request{}
	result, err := r.Reconcile(request)

	assert.NoError(err)
}
