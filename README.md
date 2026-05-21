# OSDO Operator

[![CI](https://github.com/opensecdevops/osdo-operator/actions/workflows/ci.yml/badge.svg)](https://github.com/opensecdevops/osdo-operator/actions/workflows/ci.yml)
[![GHCR](https://img.shields.io/badge/GHCR-osdo--operator-blue)](https://ghcr.io/opensecdevops/osdo-operator)

Kubernetes Operator para DevSecOps nativo de cluster. Gestiona 3 CRDs:

| CRD | Descripción |
|-----|-------------|
| `SecurityScan` | Solicita un escaneo de seguridad de un repo o imagen |
| `SecurityReport` | Resultado generado automáticamente por el operator |
| `OSDOPolicy` | Define umbrales y reglas de seguridad para el cluster |

## Instalación (Helm)

```bash
helm repo add osdo https://opensecdevops.github.io/osdo-operator
helm install osdo-operator osdo/osdo-operator -n osdo-system --create-namespace
```

## Uso rápido

```yaml
# 1. Definir una política
apiVersion: osdo.dev/v1alpha1
kind: OSDOPolicy
metadata:
  name: mi-politica
spec:
  criticalThreshold: 0
  minScore: 60
  requiredControls: [secrets, sast, sca]
  failAction: warn

---
# 2. Solicitar un escaneo
apiVersion: osdo.dev/v1alpha1
kind: SecurityScan
metadata:
  name: mi-proyecto
spec:
  target: https://github.com/mi-org/mi-proyecto
  scanType: [sast, sca, secrets]
  policyRef: mi-politica
```

```bash
# Ver resultados
kubectl get securityscans
kubectl get securityreports
kubectl describe securityreport mi-proyecto-report
```

## Licencia

Apache 2.0 — ver [LICENSE](./LICENSE)
