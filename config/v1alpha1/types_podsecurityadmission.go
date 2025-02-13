package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PSAEnforcementMode indicates the actual enforcement state of Pod Security Admission
// in the cluster. Unlike PSATargetMode, which reflects the user’s desired or “target”
// setting, PSAEnforcementMode describes the effective mode currently active.
//
// The modes define a progression from no enforcement, to label-based enforcement,
// to label-based plus global config enforcement. enforcement mode for Pod Security Admission rollout.
type PSAEnforcementMode string

const (
	// EnforcementModePrivileged indicates that no Pod Security restrictions
	// are effectively applied.
	// This aligns with a pre-rollout or fully "privileged" cluster state,
	// where neither enforce labels are set nor the global config enforces "Restricted".
	EnforcementModePrivileged PSAEnforcementMode = "Privileged"

	// EnfrocementModeLabel indicates that the cluster is enforcing Pod Security
	// labels at the Namespace level (via the PodSecurityAdmissionLabelSynchronizationController),
	// but the global kube-apiserver configuration is still "Privileged."
	EnfrocementModeLabel PSAEnforcementMode = "LabelEnforcement"

	// EnforcmentModeFull indicates that the cluster is enforcing
	// labels at the Namespace level, and the global configuration has been set
	// to "Restricted" on the kube-apiserver.
	// This represents full enforcement, where both Namespace labels and the global config
	// enforce Pod Security Admission restrictions.
	EnforcmentModeFull PSAEnforcementMode = "FullEnforcement"
)

// PSATargetMode reflects the user’s chosen (“target”) enforcement level.
type PSATargetMode string

const (
	// TargetModePrivileged indicates that the user wants no Pod Security
	// restrictions applied. The desired outcome is that the cluster remains
	// in a fully privileged (pre-rollout) state, ignoring any label enforcement
	// or global config changes.
	TargetModePrivileged PSATargetMode = "Privileged"

	// TargetModeConditional indicates that the user is willing to let the cluster
	// automatically enforce a stricter enforcement once there are no violating Namespaces.
	// If violations exist, the cluster stays in "Privileged" until those are resolved.
	// This allows a gradual move towards label and global config enforcement without
	// immediately breaking workloads that are not yet compliant.
	TargetModeConditional PSATargetMode = "Conditional"

	// TargetModeRestricted indicates that the user wants the strictest possible
	// enforcement, causing the cluster to ignore any existing violations and
	// enforce "Restricted" anyway. This reflects a final, fully enforced state.
	TargetModeRestricted PSATargetMode = "Restricted"
)

// PSAEnforcementConfig is the config for the PSA enforcement.
type PSAEnforcementConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user-settable values for configuring Pod Security Admission
	// enforcement
	Spec PSAEnforcementConfigSpec `json:"spec"`

	// status communicates the targeted enforcement mode, including any discovered
	// issues in Namespaces.
	Status PSAEnforcementConfigStatus `json:"status"`
}

// PSAEnforcementConfigSpec defines the desired configuration for Pod Security
// Admission enforcement.
type PSAEnforcementConfigSpec struct {
	// targetMode is the user-selected Pod Security Admission enforcement level.
	// Valid values are:
	//   - "Privileged": ensures the cluster runs with no restrictions
	//   - "Conditional": defers the decision to cluster-based evaluation
	//   - "Restricted": enforces the strictest Pod Security admission
	//
	// If this field is not set, it defaults to "Conditional".
	//
	// +kubebuilder:default=Conditional
	TargetMode PSATargetMode `json:"targetMode"`
}

// PSAEnforcementConfigStatus defines the observed state of Pod Security
// Admission enforcement.
type PSAEnforcementConfigStatus struct {
	EnforcementMode PSAEnforcementMode `json:"enforcmentMode"`

	// violatingNamespaces is a list of namespaces that can initially block the
	// cluster from fully enforcing a "Restricted" mode. Administrators should
	// review each listed Namespace to fix any issues to enable strict enforcement.
	//
	// If a cluster is already in "Restricted" mode and new violations emerge,
	// it remains in "Restricted" until the user explicitly switches to
	// "spec.mode = Privileged".
	//
	// To revert "Restricted" mode the Administrators need to set the
	// PSAEnfocementMode to "Privileged".
	//
	// +optional
	ViolatingNamespaces []ViolatingNamespace `json:"violatingNamespaces,omitempty"`
}

// ViolatingNamespace provides information about a namespace that cannot comply
// with the chosen enforcement mode.
type ViolatingNamespace struct {
	// name is the Namespace that has been flagged as potentially violating if
	// enforced.
	Name string `json:"name"`

	// reason is a textual description explaining why the Namespace is incompatible
	// with the requested Pod Security mode and highlights which mode is affected.
	//
	// Possible values are:
	// - PSAConfig: Misconfigured OpenShift Namespace
	// - PSAConfig: PSA label syncer disabled
	// - PSALabel: ServiceAccount with insufficient SCCs
	//
	// +optional
	Reason string `json:"reason,omitempty"`
}
