# GitOps en 30 minutos: de cero a flujo real con FluxCD

[![en](https://img.shields.io/badge/lang-en-red.svg)](./README.md)
[![es](https://img.shields.io/badge/lang-es-yellow.svg)](./README-es.md)

## Informacion de la Charla

- **Fecha**: 2026-02-25
- **Evento**: Cloud Native Community Meetup
- **Temas**: GitOps, FluxCD, Kubernetes, EKS, Terraform
- **Slides**: [slides.com](https://slides.com/shankyjs_/2026-gitops-en-30-min-con-flux)

## Descripcion

Esta charla demuestra como implementar GitOps con FluxCD en AWS EKS. Vamos de cero a un flujo completo de PR-a-Produccion con un solo `terraform apply`.

## Lo Que Aprenderas

- Fundamentos de GitOps y el modelo pull-based
- Como funciona la reconciliacion de FluxCD
- Configuracion de EKS + Flux con Terraform
- El flujo PR-a-Produccion sin pipelines de CI
- Deteccion de drift y auto-correccion en accion

## Demo

### Prerequisitos

- AWS CLI configurado con credenciales
- Terraform >= 1.0
- kubectl
- 1Password CLI (para el GitHub PAT)
- GitHub Personal Access Token con permisos de repo

### Configuracion

Toda la configuracion esta centralizada en `terraform/locals.tf`. Edita este archivo para personalizar:

```hcl
locals {
  # General
  project_name = "gitops-flux-demo"
  region       = "us-west-2"

  # EKS
  cluster_name    = "gitops-demo"
  cluster_version = "1.31"

  # GitHub / Flux
  github_owner      = "shankyjs"
  github_repository = "talks"
  flux_path         = "2026/feb-25th-gitops-flux-demo/clusters/gitops-demo"

  # ...
}
```

### Inicio Rapido

```bash
# Navegar al directorio de terraform
cd 2026/feb-25th-gitops-flux-demo/terraform

# Inicializar Terraform
terraform init

# Revisar el plan
terraform plan

# Crear todo (VPC + EKS + Flux bootstrap)
terraform apply

# Configurar kubectl
aws eks update-kubeconfig --region us-west-2 --name gitops-demo

# Verificar estado de Flux
flux get all
```

### Demo PR-a-Produccion

1. **Verificar que la app esta corriendo:**
   ```bash
   kubectl get pods -n demo-app
   kubectl get svc -n demo-app
   ```

2. **Obtener la URL del LoadBalancer:**
   ```bash
   kubectl get svc demo-app -n demo-app -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
   ```

3. **Crear un PR para actualizar la app:**
   - Editar `apps/demo-app/configmap.yaml`
   - Cambiar `Version: v1.0.0` a `Version: v2.0.0`
   - Abrir PR, revisar, hacer merge

4. **Ver la reconciliacion de Flux:**
   ```bash
   flux get kustomizations --watch
   ```

5. **Refrescar el navegador** - la app ahora muestra v2.0.0

### Demo de Deteccion de Drift

```bash
# Cambiar algo manualmente
kubectl scale deployment demo-app -n demo-app --replicas=5

# Ver como Flux lo revierte (en 5 minutos o forzar reconciliacion)
flux reconcile kustomization apps

# Verificar que las replicas volvieron a 2
kubectl get pods -n demo-app
```

### Limpieza

```bash
terraform destroy
```

## Estructura del Repositorio

```
gitops-flux-demo/
├── terraform/              # Config TF unica para VPC + EKS + Flux
│   ├── main.tf            # VPC, EKS, Flux bootstrap
│   ├── locals.tf          # Toda la configuracion en un solo lugar
│   ├── providers.tf       # Providers de AWS, Kubernetes, Flux
│   ├── data.tf            # Data sources (EKS auth, 1Password)
│   └── outputs.tf         # Outputs utiles
├── clusters/
│   └── gitops-demo/       # Objetivo de sincronizacion de Flux
│       ├── kustomization.yaml
│       ├── flux-system/   # Auto-generado por Flux bootstrap
│       └── apps/          # Kustomizations de apps
└── apps/
    └── demo-app/          # Aplicacion de ejemplo
        ├── namespace.yaml
        ├── deployment.yaml
        ├── service.yaml
        ├── configmap.yaml
        └── kustomization.yaml
```

## Recursos

- [Documentacion de FluxCD](https://fluxcd.io/docs/)
- [Flux Terraform Provider](https://registry.terraform.io/providers/fluxcd/flux/latest/docs)
- [Modulo EKS de Terraform](https://registry.terraform.io/modules/terraform-aws-modules/eks/aws/latest)
- [Principios de GitOps](https://opengitops.dev/)

## Grabacion

La grabacion estara disponible despues de la charla.

## Retroalimentacion

Preguntas o comentarios? No dudes en [abrir un issue](https://github.com/shankyjs/talks/issues) o contactarme!

---

**Presentado por**: [@shankyjs](https://github.com/shankyjs)
