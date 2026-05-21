package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	osdov1alpha1 "github.com/opensecdevops/osdo-operator/api/v1alpha1"
)

// SecurityScanReconciler reconcilia objetos SecurityScan
type SecurityScanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// OsdoAuditResult es el JSON que produce 'osdo audit --format json'
type OsdoAuditResult struct {
	Score    int                          `json:"score"`
	Level    string                       `json:"level"`
	Controls []osdov1alpha1.ControlResult `json:"controls"`
	Findings []osdov1alpha1.Finding       `json:"findings,omitempty"`
}

// +kubebuilder:rbac:groups=osdo.dev,resources=securityscans,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=osdo.dev,resources=securityscans/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=osdo.dev,resources=securityreports,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=osdo.dev,resources=osdopolicies,verbs=get;list;watch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods/log,verbs=get;list

func (r *SecurityScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. Obtener el SecurityScan
	var scan osdov1alpha1.SecurityScan
	if err := r.Get(ctx, req.NamespacedName, &scan); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// 2. Si ya está completado o fallido, no hacer nada
	if scan.Status.Phase == "Completed" || scan.Status.Phase == "Failed" {
		return ctrl.Result{}, nil
	}

	// 3. Si está Running, verificar el Job
	if scan.Status.Phase == "Running" && scan.Status.JobRef != "" {
		return r.checkJob(ctx, &scan)
	}

	// 4. Fase inicial: crear el Job de escaneo
	logger.Info("Iniciando escaneo de seguridad", "target", scan.Spec.Target)

	jobName := fmt.Sprintf("osdo-scan-%s", scan.Name)
	scanImage := scan.Spec.ScanImage
	if scanImage == "" {
		scanImage = "ghcr.io/opensecdevops/osdo-scanner:latest"
	}

	// Construir argumentos del comando osdo
	args := r.buildScanArgs(&scan)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: scan.Namespace,
			Labels: map[string]string{
				"osdo.dev/scan":                scan.Name,
				"app.kubernetes.io/managed-by": "osdo-operator",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: scan.APIVersion,
					Kind:       scan.Kind,
					Name:       scan.Name,
					UID:        scan.UID,
					Controller: boolPtr(true),
				},
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: int32Ptr(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"osdo.dev/scan": scan.Name},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "osdo-scanner",
							Image: scanImage,
							Args:  args,
							Env:   r.buildEnvVars(&scan),
						},
					},
				},
			},
		},
	}

	// Crear el Job
	if err := r.Create(ctx, job); err != nil && !errors.IsAlreadyExists(err) {
		logger.Error(err, "Error creando Job de escaneo")
		return ctrl.Result{}, err
	}

	// Actualizar estado a Running
	scan.Status.Phase = "Running"
	scan.Status.JobRef = jobName
	scan.Status.StartTime = &metav1.Time{Time: time.Now()}
	if err := r.Status().Update(ctx, &scan); err != nil {
		return ctrl.Result{}, err
	}

	// Volver a verificar en 30 segundos
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *SecurityScanReconciler) checkJob(ctx context.Context, scan *osdov1alpha1.SecurityScan) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var job batchv1.Job
	if err := r.Get(ctx, types.NamespacedName{Name: scan.Status.JobRef, Namespace: scan.Namespace}, &job); err != nil {
		if errors.IsNotFound(err) {
			scan.Status.Phase = "Failed"
			scan.Status.Message = "Job de escaneo no encontrado"
			_ = r.Status().Update(ctx, scan)
		}
		return ctrl.Result{}, err
	}

	// Job aún en progreso
	if job.Status.CompletionTime == nil && job.Status.Failed == 0 {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Job falló
	if job.Status.Failed > 0 {
		logger.Info("Job de escaneo falló", "jobRef", scan.Status.JobRef)
		scan.Status.Phase = "Failed"
		scan.Status.Message = "El Job de escaneo terminó con errores"
		scan.Status.CompletionTime = &metav1.Time{Time: time.Now()}
		return ctrl.Result{}, r.Status().Update(ctx, scan)
	}

	// Job completado exitosamente — leer logs y crear SecurityReport
	logger.Info("Job de escaneo completado", "jobRef", scan.Status.JobRef)
	return r.processJobResult(ctx, scan, &job)
}

func (r *SecurityScanReconciler) processJobResult(ctx context.Context, scan *osdov1alpha1.SecurityScan, _ *batchv1.Job) (ctrl.Result, error) {
	// En una implementación completa: leer logs del pod via Pods().GetLogs()
	// Para esta versión, usamos resultados placeholder y la integración real
	// se completa cuando el Job escribe a un ConfigMap o PVC compartido.

	// Resultado de ejemplo (en producción: leer del log del pod)
	result := OsdoAuditResult{
		Score:    0,
		Level:    "Inicial",
		Controls: []osdov1alpha1.ControlResult{},
	}

	// Crear SecurityReport
	reportName := fmt.Sprintf("%s-report", scan.Name)
	policyStatus := r.evaluatePolicy(ctx, scan, &result)

	report := &osdov1alpha1.SecurityReport{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reportName,
			Namespace: scan.Namespace,
			Labels:    map[string]string{"osdo.dev/scan": scan.Name},
			OwnerReferences: []metav1.OwnerReference{
				{APIVersion: scan.APIVersion, Kind: scan.Kind, Name: scan.Name, UID: scan.UID, Controller: boolPtr(true)},
			},
		},
		Spec: osdov1alpha1.SecurityReportSpec{
			ScanRef:      scan.Name,
			Score:        result.Score,
			Level:        result.Level,
			Controls:     result.Controls,
			Findings:     result.Findings,
			GeneratedAt:  time.Now().Format(time.RFC3339),
			PolicyStatus: policyStatus,
		},
	}

	if err := r.Create(ctx, report); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	// Actualizar estado del SecurityScan
	scan.Status.Phase = "Completed"
	scan.Status.Score = result.Score
	scan.Status.Level = result.Level
	scan.Status.ReportRef = reportName
	scan.Status.CompletionTime = &metav1.Time{Time: time.Now()}
	_ = json.Marshal(result) // evitar import no utilizado
	return ctrl.Result{}, r.Status().Update(ctx, scan)
}

func (r *SecurityScanReconciler) evaluatePolicy(ctx context.Context, scan *osdov1alpha1.SecurityScan, result *OsdoAuditResult) string {
	if scan.Spec.PolicyRef == "" {
		return "NoPolicyDefined"
	}

	var policy osdov1alpha1.OSDOPolicy
	if err := r.Get(ctx, types.NamespacedName{Name: scan.Spec.PolicyRef, Namespace: scan.Namespace}, &policy); err != nil {
		return "NoPolicyDefined"
	}

	if result.Score < policy.Spec.MinScore {
		return "NonCompliant"
	}
	return "Compliant"
}

func (r *SecurityScanReconciler) buildScanArgs(scan *osdov1alpha1.SecurityScan) []string {
	args := []string{"audit", "--format", "json"}
	if len(scan.Spec.ScanType) > 0 {
		for _, t := range scan.Spec.ScanType {
			args = append(args, "--type", t)
		}
	}
	if scan.Spec.Target != "" {
		args = append(args, "--path", "/workspace")
	}
	return args
}

func (r *SecurityScanReconciler) buildEnvVars(scan *osdov1alpha1.SecurityScan) []corev1.EnvVar {
	if scan.Spec.TokenSecret == "" {
		return nil
	}
	return []corev1.EnvVar{
		{
			Name: "GITHUB_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: scan.Spec.TokenSecret},
					Key:                  "GITHUB_TOKEN",
					Optional:             boolPtr(true),
				},
			},
		},
	}
}

func (r *SecurityScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&osdov1alpha1.SecurityScan{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func boolPtr(b bool) *bool    { return &b }
func int32Ptr(i int32) *int32 { return &i }
