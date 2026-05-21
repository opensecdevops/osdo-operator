package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OSDOPolicySpec define los umbrales y reglas de seguridad
type OSDOPolicySpec struct {
	// CriticalThreshold máximo de vulnerabilidades críticas permitidas (default: 0)
	// +kubebuilder:default=0
	// +kubebuilder:validation:Minimum=0
	CriticalThreshold int `json:"criticalThreshold"`

	// HighThreshold máximo de vulnerabilidades altas permitidas (default: 5)
	// +kubebuilder:default=5
	// +kubebuilder:validation:Minimum=0
	HighThreshold int `json:"highThreshold"`

	// SecretsThreshold máximo de secretos permitidos (default: 0)
	// +kubebuilder:default=0
	// +kubebuilder:validation:Minimum=0
	SecretsThreshold int `json:"secretsThreshold"`

	// MinScore puntuación mínima de seguridad requerida (0-100)
	// +kubebuilder:default=0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	MinScore int `json:"minScore"`

	// RequiredControls lista de controles que deben estar cubiertos
	// Valores válidos: secrets, sast, sca, container-scan, sbom, iac, signing, slsa, dast, policy-gate, license, scorecard
	RequiredControls []string `json:"requiredControls,omitempty"`

	// FailAction qué hacer cuando el scan no cumple la política
	// block: falla el pipeline | warn: advertencia | report: solo registra
	// +kubebuilder:default="warn"
	// +kubebuilder:validation:Enum=block;warn;report
	FailAction string `json:"failAction"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=op
// +kubebuilder:printcolumn:name="MinScore",type=integer,JSONPath=`.spec.minScore`
// +kubebuilder:printcolumn:name="Critical",type=integer,JSONPath=`.spec.criticalThreshold`
// +kubebuilder:printcolumn:name="High",type=integer,JSONPath=`.spec.highThreshold`
// +kubebuilder:printcolumn:name="FailAction",type=string,JSONPath=`.spec.failAction`

// OSDOPolicy define los umbrales de seguridad para los SecurityScans
type OSDOPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OSDOPolicySpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// OSDOPolicyList contiene una lista de OSDOPolicy
type OSDOPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OSDOPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OSDOPolicy{}, &OSDOPolicyList{})
}
