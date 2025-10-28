# Intro To Flux With EKS

[![en](https://img.shields.io/badge/lang-en-red.svg)](./README.md)
[![es](https://img.shields.io/badge/lang-es-yellow.svg)](./README-es.md)

## 📅 Información de la Charla

- **Fecha**: 2025-10-30
- **Evento**: GitOps & Kubernetes Meetup
- **Temas**: GitOps, Flux, EKS, Terramate, Infraestructura como Código

## 📝 Descripción

Esta demo muestra un flujo de trabajo completo de GitOps usando Flux en Amazon EKS, orquestado con Terramate para la gestión de infraestructura. Aprende cómo configurar despliegues automatizados de Kubernetes donde Git se convierte en la única fuente de verdad para tu infraestructura y aplicaciones.

## 🎯 Lo Que Aprenderás

- Cómo aprovisionar clusters EKS con Terraform y gestionarlos con stacks de Terramate
- Configurar Flux para entrega continua GitOps en Kubernetes
- Gestionar secretos de forma segura (ejemplo de integración con 1Password - adaptable a otros proveedores)
- Organizar infraestructura multi-stack con dependencias apropiadas
- Mejores prácticas de Infraestructura como Código con gestión de estado

## 🏗️ Resumen de Arquitectura

La demo está organizada en 3 stacks de Terramate que se construyen uno sobre otro:

```
┌─────────────────────────────────────────────────────────────┐
│ Stack 01: Backend S3                                        │
│ ├── Bucket S3 (para estado de Terraform)                   │
│ └── Tabla DynamoDB (para bloqueo de estado)                │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Stack 02: Cluster EKS                                       │
│ ├── VPC con subredes públicas/privadas                     │
│ ├── Plano de Control EKS (v1.31)                           │
│ ├── Grupo de Nodos Administrado (t3.medium)                │
│ ├── Controlador EBS CSI                                    │
│ └── Entrada de Acceso SSO (AdministratorAccess)            │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Stack 03: Flux GitOps                                       │
│ ├── Controladores Flux (source, kustomize, helm, notif.)   │
│ ├── Integración con Repositorio Git (GitHub)               │
│ └── Gestión de Secretos (proveedor 1Password - customizable)│
└─────────────────────────────────────────────────────────────┘
```

### Cómo se Conectan los Componentes

1. **Terramate** orquesta el orden de despliegue y comparte configuración mediante globales
2. **Stack 01** crea el backend S3 que los Stacks 02 y 03 usan para estado remoto
3. **Stack 02** aprovisiona el cluster EKS y otorga acceso admin a tu rol SSO
4. **Stack 03** instala Flux en el cluster, conectándolo a tu repositorio Git
5. **Flux** observa tu repositorio Git y aplica cambios automáticamente al cluster

## 🚀 Demo

### Prerequisitos

- [Terraform](https://www.terraform.io/) v1.12+ (ARM64 recomendado para Apple Silicon)
- [Terramate](https://terramate.io/) CLI
- [AWS CLI](https://aws.amazon.com/cli/) configurado con SSO
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- Cuenta AWS con permisos EKS

**Gestión de Secretos** (elige uno):
- [1Password CLI](https://developer.1password.com/docs/cli) (por defecto en esta demo)
- O modifica `03-flux/providers.tf` para usar tu proveedor de secretos preferido (AWS Secrets Manager, Vault, etc.)

### Configuración

1. **Actualiza los globales** en `live/terramate.tm.hcl`:
   ```hcl
   aws = {
     region = "us-west-2"  # Tu región AWS
     sso_admin_role_name = "AWSReservedSSO_AdministratorAccess_xxxxx"  # Tu rol SSO
   }

   github = {
     user       = "TuUsuarioGitHub"
     repository = "tu-repo"
   }
   ```

2. **Guarda tu PAT de GitHub** (si usas 1Password):
   - Guarda tu Personal Access Token de GitHub en 1Password
   - Actualiza la referencia en `03-flux/main.tf` si es necesario

### Ejecutando la Demo

```bash
cd live/

terramate generate # Genera tu configuración
# Inicializa y aplica todos los stacks en orden
terramate run -- sh -c 'terraform init && terraform apply'

# O aplica cada stack individualmente:
terramate run --tags s3-backend -- terraform apply
terramate run --tags eks-cluster -- terraform apply
terramate run --tags flux -- terraform apply

./scripts/update-kubeconfigs.sh
# Verifica la instalación de Flux
kubectl get pods -n flux-system

# Limpieza (destruye en orden inverso)
terramate run --reverse -- terraform destroy --auto-approve
```

### Solución de Problemas

**Usuarios de Apple Silicon**: Si encuentras errores de timeout de plugins, asegúrate de usar Terraform ARM64:
```bash
# Usa la función auxiliar de dotfiles
tfenv_switch_arch  # Alterna entre ARM64 e Intel

# Verifica la arquitectura
terraform version  # Debería mostrar "darwin_arm64"
```

## 📁 Estructura del Repositorio

```
intro-to-flux-with-eks/
├── README.md              # Versión en inglés
├── README-es.md           # Este archivo
├── metadata.yaml          # Metadata de la charla
└── live/                  # Infraestructura de la demo en vivo
    ├── terramate.tm.hcl   # Configuración global
    ├── 01-s3-backend/     # S3 + DynamoDB para estado
    ├── 02-eks-cluster/    # Cluster EKS + VPC
    ├── 03-flux/           # Configuración de Flux GitOps
    └── scripts/           # Scripts auxiliares
└── /clusters/flux-demo/   # Flux Cluster spec
```

## 📚 Recursos

- [Documentación de Flux](https://fluxcd.io/docs/)
- [Documentación de Terramate](https://terramate.io/docs/)
- [Mejores Prácticas de EKS](https://aws.github.io/aws-eks-best-practices/)
- [Proveedor AWS de Terraform](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [Principios de GitOps](https://opengitops.dev/)

## 🎥 Grabación

La grabación estará disponible después de la charla.

## 📧 Retroalimentación

¿Preguntas o comentarios? No dudes en [abrir un issue](https://github.com/shankyjs/talks/issues) o contactarme!

---

**Presentado por**: [@shankyjs](https://github.com/shankyjs)
