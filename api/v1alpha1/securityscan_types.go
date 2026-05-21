package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecurityScanSpec define el escaneo de seguridad deseado
type SecurityScanSpec struct {
	// Target es el repositorio URL o imagen Docker a escanear
	// +kubebuilder:validation:Required
	Target string `json:"target"`

	// ScanType lista los tipos de escaneo: sast, sca, secrets, container, iac, sbom
	// +kubebuilder:validation:Items=Enum=sast;sca;secrets;container;iac;sbom
	ScanType []string `json:"scanType,omitempty"`

	// PolicyRef es el nombre del OSDOPolicy a aplicar (en el mismo namespace)
	PolicyRef string `json:"policyRef,omitempty"`

	// ScanImage permite usar una imagen personalizada en lugar de la oficial
	// +kubebuilder:default="ghcr.io/opensecdevops/osdo-scanner:latest"
	ScanImage string `json:"scanImage,omitempty"`

	// TokenSecret es el nombre del Secret de Kubernetes que contiene tokens
	// (GITHUB_TOKEN, GITLAB_TOKEN) para autenticación durante el escaneo
	TokenSecret string `json:"tokenSecret,omitempty"`

	// Schedule en formato cron para escaneos periódicos (opcional)
	// Ejemplo: "0 2 * * *" para escanear cada día a las 2am
	Schedule string `json:"schedule,omitempty"`
}

// SecurityScanStatus define el estado observado del SecurityScan
type SecurityScanStatus struct {
	// Phase indica el estado actual: Pending, Running, Completed, Failed
	// +kubebuilder:validation:Enum=Pending;Running;Completed;Failed
	Phase string `json:"phase,omitempty"`

	// Score es la puntuación de seguridad (0-100)
	Score int `json:"score,omitempty"`

	// Level es el nivel de madurez DevSecOps
	// Valores: Sin DevSecOps, Inicial, Básico, Establecido, Avanzado, Completo
	Level string `json:"level,omitempty"`

	// ReportRef es el nombre del SecurityReport generado
	ReportRef string `json:"reportRef,omitempty"`

	// JobRef es el nombre del Job de Kubernetes que ejecuta el escaneo
	JobRef string `json:"jobRef,omitempty"`

	// Message contiene información adicional o el motivo del fallo
	Message string `json:"message,omitempty"`

	// StartTime cuando comenzó el último escaneo
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime cuando terminó el último escaneo
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Conditions lista de condiciones estándar de Kubernetes
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ss
// +kubebuilder:printcolumn:name="Target",type=string,JSONPath=`.spec.target`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Score",type=integer,JSONPath=`.status.score`
// +kubebuilder:printcolumn:name="Level",type=string,JSONPath=`.status.level`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// SecurityScan es el CRD para solicitar un escaneo de seguridad de OSDO
type SecurityScan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecurityScanSpec   `json:"spec,omitempty"`
	Status SecurityScanStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SecurityScanList contiene una lista de SecurityScan
type SecurityScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityScan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecurityScan{}, &SecurityScanList{})
}
