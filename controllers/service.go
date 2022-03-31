package controllers

import (
	"context"
	dustbitcomv1alpha1 "github.com/dperaza4dustbit/psql-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ensureService ensures Service is Running in a namespace.
func (r *PSQLInstanceReconciler) ensureService(request reconcile.Request,
	instance *dustbitcomv1alpha1.PSQLInstance,
	service *corev1.Service,
) (*reconcile.Result, error) {

	// See if service already exists and create if it doesn't
	found := &corev1.Service{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      service.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		err = r.Create(context.TODO(), service)

		if err != nil {
			// Service creation failed
			return &reconcile.Result{}, err
		} else {
			// Service creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendService is a code for creating a Service
func (r *PSQLInstanceReconciler) backendService(v *dustbitcomv1alpha1.PSQLInstance) *corev1.Service {
	labels := labels(v, "backend")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Spec.DatabaseName,
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Name:       "psql",
				Protocol:   corev1.ProtocolTCP,
				Port:       5432,
				TargetPort: intstr.FromInt(5432),
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	controllerutil.SetControllerReference(v, service, r.Scheme)
	return service
}
