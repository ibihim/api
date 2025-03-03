package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PSAEnforcementMode string

const (
	// PSAEnforcementModeLegacy is the previous behavior, without PSA enforcement.
	PSAEnforcementModeLegacy PSAEnforcementMode = "Legacy"

	// PSAEnforcementModeRestricted is the secure default, with PSA enforcement.
	PSAEnforcementModeRestricted PSAEnforcementMode = "Restricted"
)

type PSAEnforcementConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec gives the user the opportunity to stay in Legacy mode.
	// This field is required.
	// +required
	Spec PSAEnforcementConfigSpec `json:"spec"`

	// status is the current state of the PSA enforcement.
	// +required
	Status PSAEnforcementConfigStatus `json:"status"`
}

type PSAEnforcementConfigSpec struct {
	// enforcementMode is the mode of PSA enforcement that shows the intended behavior:
	// - Legacy: no PSA enforcement
	// - Restricted: PSA enforcement
	// - "" (empty string): the evaluated status.enforcementMode is used.
	// The default is an empty string, which means that the evaluated status.enforcementMode is used.
	// If violations are detected, the spec.enforcementMode is set to Legacy.
	// +required
	EnforcementMode PSAEnforcementMode `json:"enforcementMode"`
}

type PSAEnforcementConfigStatus struct {
	// enforcementMode is the mode of PSA enforcement that will be applied to the cluster:
	// - Legacy: no PSA enforcement
	// - Restricted: PSA enforcement
	// If unset, the cluster is Upgradeable=false.
	// +required
	EnforcementMode PSAEnforcementMode `json:"enforcementMode"`

	// violatingNamespaces is a list of namespaces that are violating the PSA enforcement policy.
	// +optional
	ViolatingNamespaces []ViolatingNamespace `json:"violatingNamespaces,omitempty"`
}

type ViolatingNamespace struct {
	// name is the name of the violating namespace.
	// +required
	Name string `json:"name"`

	// reason is the reason why the namespace is violating the PSA enforcement policy.
	// It might indicate the reason why the namespace is violating the PSA enforcement policy.
	// +optional
	Reason string `json:"reason,omitempty"`

	// lastTimeViolating is the last time when the namespace was violating the PSA enforcement policy.
	// The evaluation happens every couple of hours.
	// +required
	LastTimeViolating metav1.Time `json:"lastTimeViolating"`
}
