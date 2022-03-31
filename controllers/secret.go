package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	dustbitcomv1alpha1 "github.com/dperaza4dustbit/psql-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func generatePassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var builder strings.Builder
	for i := 0; i < length; i++ {
		builder.WriteRune(chars[rand.Intn(len(chars))])
	}
	return builder.String()
}

// ensureService ensures Service is Running in a namespace.
func (r *PSQLInstanceReconciler) ensureSecret(request reconcile.Request,
	instance *dustbitcomv1alpha1.PSQLInstance,
	secret *corev1.Secret,
) (*reconcile.Result, error) {

	// See if service already exists and create if it doesn't
	found := &corev1.Secret{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      secret.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the secret
		err = r.Create(context.TODO(), secret)

		if err != nil {
			// Secret creation failed
			return &reconcile.Result{}, err
		} else {
			// Secret creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendSecret is a function for creating a Secret
func (r *PSQLInstanceReconciler) backendSecret(v *dustbitcomv1alpha1.PSQLInstance) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("io.servicebinding.%s", v.Spec.DatabaseName),
			Namespace: v.Namespace,
		},
		Type: "servicebinding.io/postgresql",
		StringData: map[string]string{
			"type":     "postgresql",
			"provider": "redhat",
			"host":     v.Spec.DatabaseName,
			"port":     "5432",
			"username": v.Spec.UserName,
			"password": generatePassword(16),
			"database": v.Spec.DatabaseName,
		},
	}

	controllerutil.SetControllerReference(v, secret, r.Scheme)
	return secret
}
