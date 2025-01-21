package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Add other necessary imports
)

// PSAEnforcementMode is the enforcement mode for Pod Security Admission.
type PSAEnforcementMode string

const (
	// EnforcementModePrivileged indicates that no Pod Security restrictions apply.
	// This effectively doesn't change the cluster state.
	EnforcementModePrivileged PSAEnforcementMode = "Privileged"

	// EnforcementModeNoOpinion defers the enforcement decision to cluster logic.
	// In practice:
	//   - If any violating namespaces exist, the cluster remains at "Privileged".
	//   - If no violating namespaces exist, the cluster enforces "Restricted".
	// State changes:
	//   - If the state changes from "any violation" to "no violations", the cluster
	//     will start switching to enforcing "Restricted". For a controlled switch,
	//     set EnforcementMode to "Privileged" and change to "NoOpinion" back once ready.
	//   - If the state changes from "no violations" to "any violation", and the cluster
	//     settled on enforcing "Restricted", the state won't change back, except the
	//     EnforcementMode is set to "Privileged".
	EnforcementModeNoOpinion PSAEnforcementMode = "NoOpinion"

	// EnforcementModeRestricted indicates that the strictest Pod Security restrictions apply.
	// This effectively moves the cluster into the "Restricted" state, despite violations.
	EnforcementModeRestricted PSAEnforcementMode = "Restricted"
)

// PSAEnforcementConfig is the config for the PSA enforcement.
type PSAEnforcementConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user-settable values for configuring Pod Security Admission
	// enforcement
	Spec PSAEnforcementConfigSpec `json:"spec"`

	// status communicates the targeted enforcement mode, including any discovered
	// issues in namespaces.
	Status PSAEnforcementConfigStatus `json:"status"`
}

// PSAEnforcementConfigSpec defines the desired configuration for Pod Security
// Admission enforcement.
type PSAEnforcementConfigSpec struct {
	// mode is the user-selected Pod Security Admission enforcement level.
	// Valid values are:
	//   - "Privileged": ensures the cluster runs with no restrictions
	//   - "NoOpinion": defers the decision to cluster-based evaluation
	//   - "Restricted": enforces strict Pod Security admission
	//
	// If this field is not set, it defaults to "NoOpinion".
	//
	// +kubebuilder:default=NoOpinion
	Mode PSAEnforcementMode `json:"mode"`
}

// PSAEnforcementConfigStatus defines the observed state of Pod Security
// Admission enforcement.
type PSAEnforcementConfigStatus struct {
	// effectiveMode is the actual Pod Security Admission mode being enforced by
	// the cluster. This will differ from spec.mode if spec.mode is NoOpinion,
	// then the initial decision will be based on violating Namespaces.
	EffectiveMode PSAEnforcementMode `json:"effectiveMode,omitempty"`

	// violatingNamespaces is a list of namespaces that can initially block the
	// cluster from fully enforcing a "Restricted" mode. Administrators should
	// review each listed namespace and fix any issues to enable strict
	// enforcement.
	//
	// Violations after "Restricted" mode is being applied, have no effect. To
	// revert "Restricted" mode the Administrators need to set the
	// PSAEnfocementMode
	//
	// +optional
	ViolatingNamespaces []ViolatingNamespace `json:"violatingNamespaces,omitempty"`
}

// ViolatingNamespace provides information about a namespace that cannot comply
// with the chosen enforcement mode.
type ViolatingNamespace struct {
	// name is the namespace that has been flagged as potentially violating if
	// enfocred.
	Name string `json:"name"`

	// reason is an optional textual description explaining why the namespace is
	// incompatible with the requested Pod Security mode.
	// +optional
	Reason string `json:"reason,omitempty"`
}
