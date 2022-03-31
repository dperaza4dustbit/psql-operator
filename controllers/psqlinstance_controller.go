/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dustbitcomv1alpha1 "github.com/dperaza4dustbit/psql-operator/api/v1alpha1"
)

// PSQLInstanceReconciler reconciles a PSQLInstance object
type PSQLInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dustbit.com,resources=psqlinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dustbit.com,resources=psqlinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dustbit.com,resources=psqlinstances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PSQLInstance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *PSQLInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("PSQLInstance", req.NamespacedName)

	// Fetch the PSQLInstance instance
	instance := &dustbitcomv1alpha1.PSQLInstance{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if this Secret already exists
	secret := r.backendSecret(instance)
	result, err := r.ensureSecret(req, instance, secret)
	if result != nil {
		log.Error(err, "Secret Not ready")
		return *result, err
	}

	// Update Instance with status
	instance.Status.Binding.Name = secret.Name
	err = r.Status().Update(context.Background(), instance)
	if err != nil {
		log.Error(err, "Binding Status Failed")
		return *result, err
	}

	// Check if this Deployment already exists
	result, err = r.ensureDeployment(req, instance, r.backendDeployment(instance))
	if result != nil {
		log.Error(err, "Deployment Not ready")
		return *result, err
	}

	// Check if this Service already exists
	result, err = r.ensureService(req, instance, r.backendService(instance))
	if result != nil {
		log.Error(err, "Service Not ready")
		return *result, err
	}

	// Secret, Deployment and Service already exists - don't requeue
	log.Info("Skip reconcile: Secret, Deployment and service already exists")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PSQLInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dustbitcomv1alpha1.PSQLInstance{}).
		Complete(r)
}
