package v1alpha1

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SingletonClusterName                      = "cluster"
	InternetReachableFromMaster ConditionType = "InternetReachableFromMaster"
	InternetReachableFromWorker ConditionType = "InternetReachableFromWorker"
	MachineValid                ConditionType = "MachineValid"
	ServicePrincipalValid       ConditionType = "ServicePrincipalValid"
)

// AllConditionTypes is a operator conditions currently in use, any condition not in this list is not
// added to the operator.status.conditions list
func AllConditionTypes() []ConditionType {
	return []ConditionType{InternetReachableFromMaster, InternetReachableFromWorker, MachineValid, ServicePrincipalValid}
}

// ClusterChecksTypes represents checks performed on the cluster to verify basic functionality
func ClusterChecksTypes() []ConditionType {
	return []ConditionType{InternetReachableFromMaster, InternetReachableFromWorker, MachineValid, ServicePrincipalValid}
}

type GenevaLoggingSpec struct {
	// +kubebuilder:validation:Pattern:=`[0-9]+.[0-9]+`
	ConfigVersion string `json:"configVersion,omitempty"`
	// +kubebuilder:validation:Enum=DiagnosticsProd;Test
	MonitoringGCSEnvironment string `json:"monitoringGCSEnvironment,omitempty"`
}

type InternetCheckerSpec struct {
	URLs []string `json:"urls,omitempty"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// ResourceID is the Azure resourceId of the cluster
	ResourceID          string              `json:"resourceId,omitempty"`
	Domain              string              `json:"domain,omitempty"`
	ACRDomain           string              `json:"acrDomain,omitempty"`
	AZEnvironment       string              `json:"azEnvironment,omitempty"`
	Location            string              `json:"location,omitempty"`
	InfraID             string              `json:"infraId,omitempty"`
	ArchitectureVersion int                 `json:"architectureVersion,omitempty"`
	GenevaLogging       GenevaLoggingSpec   `json:"genevaLogging,omitempty"`
	InternetChecker     InternetCheckerSpec `json:"internetChecker,omitempty"`
	VnetID              string              `json:"vnetId,omitempty"`
	APIIntIP            string              `json:"apiIntIP,omitempty"`
	IngressIP           string              `json:"ingressIP,omitempty"`

	Features FeaturesSpec `json:"features,omitempty"`
}

// FeaturesSpec defines ARO operator feature gates
type FeaturesSpec struct {
	ReconcileNSGs                  bool `json:"reconcileNSGs,omitempty"`
	ReconcileAlertWebhook          bool `json:"reconcileAlertWebhook,omitempty"`
	ReconcileDNSMasq               bool `json:"reconcileDNSMasq,omitempty"`
	ReconcileGenevaLogging         bool `json:"reconcileGenevaLogging,omitempty"`
	ReconcileMonitoringConfig      bool `json:"reconcileMonitoringConfig,omitempty"`
	ReconcileNodeDrainer           bool `json:"reconcileNodeDrainer,omitempty"`
	ReconcilePullSecret            bool `json:"reconcilePullSecret,omitempty"`
	ReconcileRouteFix              bool `json:"reconcileRouteFix,omitempty"`
	ReconcileWorkaroundsController bool `json:"reconcileWorkaroundsController,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	OperatorVersion   string      `json:"operatorVersion,omitempty"`
	Conditions        []Condition `json:"conditions,omitempty"`
	RedHatKeysPresent []string    `json:"redHatKeysPresent,omitempty"`
}

// +kubebuilder:object:root=true
// +genclient
// +genclient:nonNamespaced
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterList contains a list of Cluster
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

// ConditionStatus indicates the status of a condition (true, false, or unknown).
type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition;
// "ConditionFalse" means a resource is not in the condition; "ConditionUnknown" means kubernetes
// can't decide if a resource is in the condition or not. In the future, we could add other
// intermediate conditions, e.g. ConditionDegraded.
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// ConditionType is a valid value for APIServiceCondition.Type
type ConditionType string

const (
	// Available indicates that the service exists and is reachable
	Available ConditionType = "Available"
)

// Condition represents an observation of an object's
// state. Conditions are an extension mechanism intended to be used
// when the details of an observation are not a priori known or would
// not apply to all instances of a given Kind. \n Conditions should
// be added to explicitly convey properties that users and components
// care about rather than requiring those properties to be inferred
// from other observations. Once defined, the meaning of a Condition
// can not be changed arbitrarily - it becomes part of the API, and
// has the same backwards- and forwards-compatibility concerns of
// any other part of the API.
type Condition struct {
	// ConditionType is the type of the condition and
	// is typically a CamelCased word or short phrase. \n Condition
	// types should indicate state in the \"abnormal-true\" polarity.
	// For example, if the condition indicates when a policy is invalid,
	// the \"is valid\" case is probably the norm, so the condition
	// should be called \"Invalid\".
	Type ConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// ConditionReason is intended to be a one-word, CamelCase representation of
	// the category of cause of the current status.
	// It is intended to be used in concise output, such as one-line
	// kubectl get output, and in summarizing occurrences of causes.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
