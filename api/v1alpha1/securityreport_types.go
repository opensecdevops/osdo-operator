package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ControlResult representa el resultado de un control de seguridad OSDO
type ControlResult struct {
	// ID del control (p.ej. "secrets", "sast", "sca")
	ID string `json:"id"`

	// Name nombre descriptivo del control
	Name string `json:"name"`

	// Category categoría: esencial, recomendado, completo
	Category string `json:"category"`

	// Points puntos posibles para este control
	Points int `json:"points"`

	// Score puntos obtenidos
	Score int `json:"score"`

	// Covered indica si el control está cubierto
	Covered bool `json:"covered"`

	// Evidence evidencia que demuestra el resultado
	Evidence string `json:"evidence,omitempty"`
}

// Finding representa un hallazgo de seguridad
type Finding struct {
	// RuleID identificador de la regla que generó el hallazgo
	RuleID string `json:"ruleId,omitempty"`

	// Message descripción del hallazgo
	Message string `json:"message"`

	// Severity: critical, high, medium, low, info
	// +kubebuilder:validation:Enum=critical;high;medium;low;info
	Severity string `json:"severity"`

	// Location donde se encontró el hallazgo
	Location string `json:"location,omitempty"`

	// Tool herramienta que detectó el hallazgo
	Tool string `json:"tool,omitempty"`
}

// SecurityReportSpec define el contenido del reporte de seguridad
type SecurityReportSpec struct {
	// ScanRef nombre del SecurityScan que generó este reporte
	ScanRef string `json:"scanRef"`

	// Score puntuación total de seguridad (0-100)
	Score int `json:"score"`

	// Level nivel de madurez DevSecOps
	Level string `json:"level"`

	// Controls lista de controles evaluados
	Controls []ControlResult `json:"controls,omitempty"`

	// Findings lista de hallazgos encontrados (puede estar vacía)
	Findings []Finding `json:"findings,omitempty"`

	// GeneratedAt timestamp de generación del reporte
	GeneratedAt string `json:"generatedAt"`

	// PolicyStatus resultado de la evaluación contra OSDOPolicy
	// Valores: Compliant, NonCompliant, NoPolicyDefined
	PolicyStatus string `json:"policyStatus,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=sr
// +kubebuilder:printcolumn:name="ScanRef",type=string,JSONPath=`.spec.scanRef`
// +kubebuilder:printcolumn:name="Score",type=integer,JSONPath=`.spec.score`
// +kubebuilder:printcolumn:name="Level",type=string,JSONPath=`.spec.level`
// +kubebuilder:printcolumn:name="Policy",type=string,JSONPath=`.spec.policyStatus`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// SecurityReport es generado automáticamente por el operator después de un scan
type SecurityReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SecurityReportSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// SecurityReportList contiene una lista de SecurityReport
type SecurityReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityReport `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecurityReport{}, &SecurityReportList{})
}
