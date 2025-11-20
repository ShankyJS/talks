# 🎤 Repositorio de Charlas y Demos

[![en](https://img.shields.io/badge/lang-en-red.svg)](../README.md)
[![es](https://img.shields.io/badge/lang-es-yellow.svg)](./README-es.md)

## 👋 Sobre Mí

¡Hola! Soy **Shanky** ([@shankyjs](https://github.com/shankyjs)), un Sr. Platform Engineer apasionado por las tecnologías cloud-native, DevOps y software de código abierto. Me encanta compartir conocimiento a través de charlas, talleres y demostraciones.

<div align="center">
  <a href="https://github.com/shankyjs">
    <img src="https://github.com/shankyjs.png" width="150" alt="Shanky"/>
  </a>
</div>

**Participación en la Comunidad:**
- 🇨🇦 Organizador de [Cloud Native Vancouver](https://community.cncf.io/cloud-native-vancouver/)
- 🇸🇻 Organizador de [Cloud Native San Salvador](https://community.cncf.io/cloud-native-san-salvador/)

## 📚 Sobre Este Repositorio

¡Bienvenido a mi repositorio de charlas y demos! Aquí es donde recopilo y comparto todas las presentaciones, demostraciones y ejemplos de código que he creado a lo largo de los años. Ya sea de conferencias, meetups, talleres o eventos comunitarios, encontrarás los recursos aquí.

Cada charla incluye:
- 📝 Materiales de presentación y diapositivas
- 💻 Código y ejemplos de demos
- 📖 Instrucciones paso a paso
- 🔗 Recursos y referencias adicionales

## 📊 Estadísticas

- 🎤 **Total de Charlas**: 2
- ✅ **Pasadas**: 2
- 🔜 **Próximas**: 0
- 🏷️ **Temas Principales**: AWS (1), GitOps (1), Go (1)

## 📑 Índice de Charlas

Explora todas las charlas por año, tema y evento. Haz clic en cualquier charla para acceder a la demo completa, código y materiales.

### 2025

| Fecha | Título de la Charla | Temas | Evento/Ubicación | Materiales |
|-------|---------------------|-------|------------------|------------|
| 2025-11-19 | [**Otel Jaeger Go Services**](./2025/nov-19th-otel-jaeger-go-services) | Otel, Jaeger, Go | Cloud Native Vancouver: Nov 2025 | [EN](./2025/nov-19th-otel-jaeger-go-services/README.md) / [ES](./2025/nov-19th-otel-jaeger-go-services/README-es.md) |
| 2025-10-30 | [**Intro To Flux With EKS**](./2025/oct-30th-intro-to-flux-with-eks) | GitOps, AWS, Kubernetes | October 30th Cloud Native Vancouver event | [EN](./2025/oct-30th-intro-to-flux-with-eks/README.md) / [ES](./2025/oct-30th-intro-to-flux-with-eks/README-es.md) |


### Próximamente 🚀

¡Más charlas y demos se agregarán aquí a medida que sucedan!

---

## 🏷️ Buscar por Tema

- **AWS**: [Intro To Flux With EKS (2025)](./2025/oct-30th-intro-to-flux-with-eks)
- **GitOps**: [Intro To Flux With EKS (2025)](./2025/oct-30th-intro-to-flux-with-eks)
- **Go**: [Otel Jaeger Go Services (2025)](./2025/nov-19th-otel-jaeger-go-services)
- **Jaeger**: [Otel Jaeger Go Services (2025)](./2025/nov-19th-otel-jaeger-go-services)
- **Kubernetes**: [Intro To Flux With EKS (2025)](./2025/oct-30th-intro-to-flux-with-eks)
- **Otel**: [Otel Jaeger Go Services (2025)](./2025/nov-19th-otel-jaeger-go-services)


## 🤝 Contribuir

¿Encontraste un error o quieres mejorar algo? ¡No dudes en abrir un issue o enviar un pull request!

## 📫 Contacto

- GitHub: [@shankyjs](https://github.com/shankyjs)
- ¡No dudes en contactarme si tienes preguntas sobre cualquiera de las demos o charlas!

## 📄 Licencia

A menos que se especifique lo contrario, todo el contenido de este repositorio está disponible con fines educativos. Por favor referencia este repositorio si utilizas alguno de los materiales.

---

⭐ ¡Si encuentras estos recursos útiles, considera darle una estrella a este repositorio!

## Contribuye

```bash
# 1. Compilar herramientas de automatización
make build

# Esto compila todas las herramientas de automatización:
# - create-talk (crear nuevos directorios de charlas)
# - generate-index (actualizar índice de charlas)
# - check-metadata (validar archivos de metadata)
# - generate-stats (generar estadísticas)

# 2. Instalar hooks de pre-commit (opcional pero recomendado)
pip install pre-commit  # o brew install pre-commit
pre-commit install
```

> **Nota**: Toda la automatización está construida usando Go. Ejecuta `make build` para compilar los binarios.

### Crear una Nueva Charla

```bash
# Usar el comando Makefile
make create-talk DATE=2025-11-15 SLUG=mi-charla-increible

# O usar el alias más corto
make new DATE=2025-11-15 SLUG=mi-charla-increible

# Esto crea:
# - 2025/nov-15th-mi-charla-increible/
# - metadata.yaml (¡edita esto!)
# - README.md
# - README-es.md
```

### Actualizar el Índice

```bash
# Después de crear o editar charlas
make update-index

# O simplemente
make regen
```

### Hooks de Pre-commit

Una vez instalados, los hooks de pre-commit:
- ✅ Auto-generan el índice al hacer commit
- ✅ Validan los archivos de metadata
- ✅ Verifican archivos faltantes
- ✅ Corrigen espacios en blanco al final

```bash
# Ejecutar manualmente
pre-commit run --all-files
```

### Comandos Rápidos

```bash
make help           # Mostrar todos los comandos
make build          # Compilar herramientas de automatización
make install        # Alias para build
make create-talk    # Crear nueva charla (requiere DATE y SLUG)
make new            # Alias para create-talk
make update-index   # Regenerar índice
make generate-stats # Generar estadísticas
make check          # Validar metadata
make clean          # Limpiar
```

### Flujo de Trabajo de Ejemplo

```bash
# 1. Crear charla
make create-talk DATE=2025-12-10 SLUG=secretos-kubernetes

# 2. Editar metadata
vim 2025/dec-10th-secretos-kubernetes/metadata.yaml

# 3. Agregar contenido
vim 2025/dec-10th-secretos-kubernetes/README.md

# 4. Actualizar índice
make update-index

# 5. Hacer commit (¡pre-commit hace el resto!)
git add .
git commit -m "feat: Agregar charla de secretos en Kubernetes"
```
