/*
Copyright 2023.

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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	batchv1 "k8s.io/api/batch/v1"
	// batchv1beta1 "k8s.io/api/batch/v1beta1"

	webappv1 "uday-kiran-k/operator/api/v1"
)

// HealthcheckReconciler reconciles a Healthcheck object
type HealthcheckReconciler struct {
	client.Client
	// Log logr.Logger
	Scheme *runtime.Scheme
	// Clock
}

// type realClock struct{}

// func (_ realClock) Now() time.Time { return time.Now() }

// // clock knows how to get the current time.
// // It can be used to fake out timing for testing.
// type Clock interface {
// 	Now() time.Time
// }

// func ignoreNotFound(err error) error {
// 	if apierrs.IsNotFound(err) {
// 		return nil
// 	}
// 	return err
// }

// var (
// 	scheduledTimeAnnotation = "batch.tutorial.kubebuilder.io/scheduled-at"
// )

//+kubebuilder:rbac:groups=webapp.udaykk.me,resources=healthchecks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.udaykk.me,resources=healthchecks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.udaykk.me,resources=healthchecks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Healthcheck object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *HealthcheckReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	// Fetch the HealthCheck object
	healthCheck := &webappv1.Healthcheck{}
	err := r.Get(ctx, req.NamespacedName, healthCheck)
	if err != nil {
		if errors.IsNotFound(err) {
			// HealthCheck resource not found, handle deletion if required
			log.Info("HealthCheck resource not found or deleted")
			return ctrl.Result{}, nil
		}
		// Error fetching HealthCheck
		log.Error(err, "Failed to get HealthCheck resource")
		return ctrl.Result{}, err
	}
	finalizerName := "batch.tutorial.kubebuilder.io/finalizer"

	log.Info("Reconciling HealthCheck", "Name", healthCheck.Name)
	log.Info("specs", "spec", healthCheck.Spec)

	payload := map[string]interface{}{
		"name":                      healthCheck.Spec.Name,
		"uri":                       healthCheck.Spec.URI + "/api/v1/kafka/publish",
		"is_paused":                 healthCheck.Spec.IsPaused,
		"num_retries":               healthCheck.Spec.NumRetries,
		"uptime_sla":                healthCheck.Spec.UptimeSLA,
		"response_time_sla":         healthCheck.Spec.ResponseTimeSLA,
		"use_ssl":                   healthCheck.Spec.UseSSL,
		"response_status_code":      healthCheck.Spec.ResponseStatusCode,
		"check_interval_in_seconds": healthCheck.Spec.CheckIntervalInSeconds,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error creating JSON:", err)
		return ctrl.Result{}, err
	}

	// Convert JSON data to string
	jsonString := string(jsonData)

	// Build the CronJob based on HealthCheck Spec
	// Build the CronJob based on HealthCheck Spec
	cronJob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "health-check-" + healthCheck.Name,
			Namespace: req.Namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: "*/" + strconv.Itoa(int(healthCheck.Spec.CheckIntervalInSeconds)) + " * * * *", // Schedule based on check interval
			JobTemplate: batchv1.JobTemplateSpec{ // Create a JobTemplateSpec
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "health-check-container",
									Image: "curlimages/curl:latest", // Use curl image
									Args: []string{
										"curl",
										"--location",
										"http://" + healthCheck.Spec.URI + ".consumer.svc.cluster.local/api/v1/kafka/publish",
										"--header",
										"Content-Type: application/json",
										"--data",
										jsonString,
										// `{
										// 	"name": "` + healthCheck.Spec.Name + `",
										// 	"uri": "http://" + healthCheck.Spec.URI + ".consumer.svc.cluster.local/api/v1/kafka/publish"
										// 	"is_paused": ` + strconv.FormatBool(healthCheck.Spec.IsPaused) + `,
										// 	"num_retries": ` + strconv.Itoa(int(healthCheck.Spec.NumRetries)) + `,
										// 	"uptime_sla": ` + strconv.Itoa(int(healthCheck.Spec.UptimeSLA)) + `,
										// 	"response_time_sla": ` + strconv.Itoa(int(healthCheck.Spec.ResponseTimeSLA)) + `,
										// 	"use_ssl": ` + strconv.FormatBool(healthCheck.Spec.UseSSL) + `,
										// 	"response_status_code": ` + strconv.Itoa(int(healthCheck.Spec.ResponseStatusCode)) + `,
										// 	"check_interval_in_seconds": ` + strconv.Itoa(int(healthCheck.Spec.CheckIntervalInSeconds)) + `
										// }`,
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}

	// Set HealthCheck instance as the owner and controller
	if err := controllerutil.SetControllerReference(healthCheck, cronJob, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference for CronJob")
		return ctrl.Result{}, err
	}

	// Check if CronJob already exists, if not create it
	existingCronJob := &batchv1.CronJob{}
	err = r.Get(ctx, types.NamespacedName{Name: cronJob.Name, Namespace: cronJob.Namespace}, existingCronJob)
	if err != nil && errors.IsNotFound(err) {
		// Create the CronJob
		err = r.Create(ctx, cronJob)
		if err != nil {
			log.Error(err, "Failed to create CronJob")
			return ctrl.Result{}, err
		}
		log.Info("Created CronJob", "Name", cronJob.Name)
		// CronJob created successfully
		return ctrl.Result{}, nil
	} else if err != nil {
		// Error while checking CronJob existence
		log.Error(err, "Failed to check CronJob existence")
		return ctrl.Result{}, err
	}

	// CronJob already exists, update it if needed
	if !controllerutil.ContainsFinalizer(existingCronJob, finalizerName) {
		controllerutil.AddFinalizer(existingCronJob, finalizerName)
		cronJob.Spec = constructCronJobSpecFromHealthCheck(healthCheck)
		existingCronJob.Spec.Schedule = "*/" + strconv.Itoa(int(healthCheck.Spec.CheckIntervalInSeconds)) + " * * * *" // Update schedule based on check interval
		if err := r.Update(ctx, existingCronJob); err != nil {
			log.Error(err, "Failed to update CronJob")
			return ctrl.Result{}, err
		}
		log.Info("Updated CronJob", "Name", cronJob.Name)
	}

	// Check if the HealthCheck resource is being deleted
	if !healthCheck.ObjectMeta.DeletionTimestamp.IsZero() {
		// The HealthCheck resource is being deleted, handle cleanup here
		// Delete associated CronJob
		cronJob := &batchv1.CronJob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "health-check-" + healthCheck.Name,
				Namespace: req.Namespace,
			},
		}
		if err := r.Delete(ctx, cronJob); err != nil && !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		// Check if CronJob has been deleted and remove finalizer from HealthCheck
		if errors.IsNotFound(err) {
			controllerutil.RemoveFinalizer(healthCheck, finalizerName)
			if err := r.Update(ctx, healthCheck); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

func constructCronJobSpecFromHealthCheck(healthCheck *webappv1.Healthcheck) batchv1.CronJobSpec {
	// Construct and return the CronJobSpec based on the HealthCheck
	return batchv1.CronJobSpec{
		Schedule: "*/" + strconv.Itoa(int(healthCheck.Spec.CheckIntervalInSeconds)) + " * * * *", // Schedule based on check interval
		JobTemplate: batchv1.JobTemplateSpec{ // Create a JobTemplateSpec
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "health-check-container",
								Image: "curlimages/curl:latest", // Use curl image
								Args: []string{
									"curl",
									"--location",
									"http://localhost:8082/api/v1/kafka/publish", // Update with your API endpoint
									"--header",
									"Content-Type: application/json",
									"--data",
									`{
									"name": "` + healthCheck.Spec.Name + `",
									"uri": "` + healthCheck.Spec.URI + `",
									"is_paused": ` + strconv.FormatBool(healthCheck.Spec.IsPaused) + `,
									"num_retries": ` + strconv.Itoa(int(healthCheck.Spec.NumRetries)) + `,
									"uptime_sla": ` + strconv.Itoa(int(healthCheck.Spec.UptimeSLA)) + `,
									"response_time_sla": ` + strconv.Itoa(int(healthCheck.Spec.ResponseTimeSLA)) + `,
									"use_ssl": ` + strconv.FormatBool(healthCheck.Spec.UseSSL) + `,
									"response_status_code": ` + strconv.Itoa(int(healthCheck.Spec.ResponseStatusCode)) + `,
									"check_interval_in_seconds": ` + strconv.Itoa(int(healthCheck.Spec.CheckIntervalInSeconds)) + `
									}`,
								},
							},
						},
						RestartPolicy: corev1.RestartPolicyNever,
					},
				},
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *HealthcheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Healthcheck{}).
		Complete(r)
}
