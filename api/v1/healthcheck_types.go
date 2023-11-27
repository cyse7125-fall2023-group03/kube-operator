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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HealthcheckSpec defines the desired state of Healthcheck
type HealthcheckSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Healthcheck. Edit healthcheck_types.go to remove/update
	Foo                    string `json:"foo,omitempty"`
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	URI                    string `json:"uri,omitempty"`
	IsPaused               bool   `json:"is_paused,omitempty"`
	NumRetries             int    `json:"num_retries,omitempty"`
	UptimeSLA              int    `json:"uptime_sla,omitempty"`
	ResponseTimeSLA        int    `json:"response_time_sla,omitempty"`
	UseSSL                 bool   `json:"use_ssl,omitempty"`
	ResponseStatusCode     int    `json:"response_status_code,omitempty"`
	CheckIntervalInSeconds int    `json:"check_interval_in_seconds,omitempty"`
	CheckCreated           string `json:"check_created,omitempty"`
	CheckUpdated           string `json:"check_updated,omitempty"`
}

// HealthcheckStatus defines the observed state of Healthcheck
type HealthcheckStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// LastExecutionTime time.Time `json:"inner"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// A list of pointers to currently running jobs.
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	// Information when was the last time the job was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Healthcheck is the Schema for the healthchecks API
type Healthcheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HealthcheckSpec   `json:"spec,omitempty"`
	Status HealthcheckStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HealthcheckList contains a list of Healthcheck
type HealthcheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Healthcheck `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Healthcheck{}, &HealthcheckList{})
}
