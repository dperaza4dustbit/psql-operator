package controllers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dustbitcomv1alpha1 "github.com/dperaza4dustbit/psql-operator/api/v1alpha1"
)

func labels(v *dustbitcomv1alpha1.PSQLInstance, tier string) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app":   v.Spec.DatabaseName,
		"tier":  tier,
		"phase": "test",
	}
}

// ensureDeployment ensures Deployment resource presence in given namespace.
func (r *PSQLInstanceReconciler) ensureDeployment(request reconcile.Request,
	instance *dustbitcomv1alpha1.PSQLInstance,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendDeployment is a code for Creating Deployment
func (r *PSQLInstanceReconciler) backendDeployment(v *dustbitcomv1alpha1.PSQLInstance) *appsv1.Deployment {

	labels := labels(v, "backend")
	runAsNonRoot := true
	size := int32(1)
	deploymentName := fmt.Sprintf("%s-bee", v.Spec.DatabaseName)
	metaAnnotation := make(map[string]string)
	metaAnnotation["app.kubernetes.io/part-of"] = "ssm"
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deploymentName,
			Namespace:   v.Namespace,
			Annotations: metaAnnotation,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           "registry.access.redhat.com/rhscl/postgresql-10-rhel7:latest",
						ImagePullPolicy: corev1.PullAlways,
						Name:            deploymentName,
						Env: []corev1.EnvVar{
							{
								Name: "POSTGRESQL_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("io.servicebinding.%s", v.Spec.DatabaseName),
										},
										Key: "password",
									},
								},
							},
							{
								Name: "POSTGRESQL_USER",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: fmt.Sprintf("io.servicebinding.%s", v.Spec.DatabaseName),
										},
										Key: "username",
									},
								},
							},
							{
								Name:  "POSTGRESQL_DATABASE",
								Value: v.Spec.DatabaseName,
							}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 5432,
						}},
						SecurityContext: &corev1.SecurityContext{
							RunAsNonRoot: &runAsNonRoot,
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("100m"),
								corev1.ResourceMemory: resource.MustParse("128Mi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("100m"),
								corev1.ResourceMemory: resource.MustParse("128Mi"),
							},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								TCPSocket: &corev1.TCPSocketAction{
									Port: intstr.FromInt(5432),
								},
							},
							InitialDelaySeconds: 15,
							PeriodSeconds:       20,
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}
