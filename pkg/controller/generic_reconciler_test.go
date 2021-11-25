package controller

import (
	"context"
	"testing"

	"github.com/app-sre/deployment-validation-operator/pkg/testutils"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestGenericReconciler(t *testing.T) {
	deployment, err := testutils.CreateDeploymentFromTemplate(
		testutils.NewTemplateArgs(), "")
	if err != nil {
		t.Errorf("Error creating deployment from template %v", err)
	}

	// Objects to track in the fake client.
	objs := []client.Object{&deployment}

	// Create a fake client to mock API calls.
	//	client := fake.NewFakeClient(objs...)
	client := fake.NewClientBuilder().WithObjects(objs...)

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{Name: "foo", Namespace: "bar"},
	}

	gr := &GenericReconciler{
		client:         client.Build(),
		reconciledKind: testutils.ObjectKind(&deployment),
		reconciledObj:  &deployment}

	_, err = gr.Reconcile(context.TODO(), request)
	if err != nil {
		t.Errorf("Error reconciling %v", err)
	}

	// since we're not modifying k8s objects, not much to see here except for
	// checking metrics registered, but that is done in the validation tests
}

func TestGenericReconcilerNamespaceDeletion(t *testing.T) {
	namespace, err := testutils.CreateNamespaceFromTemplate(
		testutils.NewTemplateArgs(), "default",
	)
	if err != nil {
		t.Errorf("Error creating namespace from template %v", err)
	}

	// deletedNamespace, err := testutils.CreateNamespaceFromTemplate(
	// 	testutils.NewTemplateArgs(), "deleted",
	// )
	// if err != nil {
	// 	t.Errorf("Error creating namespace from template %v", err)
	// }

	// deployment, err := testutils.CreateDeploymentFromTemplate(
	// 	testutils.NewTemplateArgs(), "default",
	// )
	// if err != nil {
	// 	t.Errorf("Error creating deployment from template %v", err)
	// }

	// Objects to track in the fake client.
	objs := []client.Object{&namespace}

	// Create a fake client to mock API calls.
	//	client := fake.NewFakeClient(objs...)
	client := fake.NewClientBuilder().WithObjects(objs...)

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"},
	}

	gr := &GenericReconciler{
		client:         client.Build(),
		reconciledKind: testutils.ObjectKind(&namespace),
		reconciledObj:  &namespace}

	out, err := gr.Reconcile(context.TODO(), request)
	t.Log(out)
	if err != nil {
		t.Errorf("Error reconciling %v", err)
	}

	// // since we're not modifying k8s objects, not much to see here except for
	// // checking metrics registered, but that is done in the validation tests
}
