# Servicios OTEL, Jaeger y Go

Este proyecto configura un entorno de desarrollo local con OpenTelemetry, Jaeger y servicios Golang (Frontend y Backend) utilizando Kind y Helmfile.

## Requisitos previos

- Docker
- `kind` (instalado en el host, o asegúrate de que `make cluster-up` pueda ejecutarlo)
- `make`

## Uso

Todo el flujo de trabajo se gestiona a través del `Makefile`.

### 1. Crear Clúster y Desplegar

```bash
make all
```

Este comando:
1.  Creará un clúster Kind (si no existe).
2.  Desplegará los servicios Jaeger, OTEL Collector, Backend y Frontend usando Helmfile (ejecutándose en Docker).

### 2. Acceder a los Servicios

Puedes usar los siguientes objetivos de make para acceder fácilmente a los servicios:

- **Interfaz de usuario de Jaeger**:
  ```bash
  make port-forward-jaeger
  ```
  Abre [http://localhost:16686](http://localhost:16686).

- **Aplicación Frontend**:
  ```bash
  make port-forward-app
  ```
  Abre [http://localhost:8080](http://localhost:8080).

### 3. Kubeconfig

Por defecto, este proyecto utiliza un archivo `.kube/config` local para mantener limpio tu entorno host y asegurar que las herramientas Dockerizadas funcionen correctamente.

Si deseas usar `kubectl` desde tu host sin especificar el archivo de configuración, ejecuta:

```bash
make kubeconfig-export
```

Esto fusionará la configuración del clúster en tu `~/.kube/config` predeterminado.

### 4. Limpieza

```bash
make clean
```
Esto destruye el clúster Kind.

## Arquitectura

- **Jaeger**: Backend de rastreo todo en uno (almacenamiento en memoria).
- **OTEL Collector**: Recibe trazas de las aplicaciones y las exporta a Jaeger.
- **Backend**: Servicio Golang (simulado con Nginx por ahora).
- **Frontend**: Servicio Golang (simulado con Nginx por ahora).

## Configuración

- `conf/helmfile.yaml`: Descriptor de despliegue principal.
- `conf/values/`: Archivos de valores para los charts.
- `charts/`: Charts de Helm locales para la aplicación.
