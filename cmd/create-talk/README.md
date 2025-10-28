# 📝 Template System Guide

This document explains the template system for generating talk files.

## 🎯 Overview

The talks repository uses templates to generate consistent, well-structured talk files. Templates are stored in `cmd/create-talk/templates/` and are rendered when you create a new talk.

## 📁 Template Files

### Available Templates

```
cmd/create-talk/templates/
├── metadata.yaml.tmpl    # Talk metadata template
├── README.md.tmpl        # English README template
└── README-es.md.tmpl     # Spanish README template
```

## 🔧 How Templates Work

### 1. Template Variables

Templates use these variables (automatically filled by `bin/create-talk`):

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.Title}}` | Talk title (auto-generated from slug) | "Kubernetes Scaling" |
| `{{.Date}}` | Talk date | "2025-11-15" |
| `{{.Event}}` | Event name | "Conference/Meetup Name" |
| `{{.Description}}` | Brief description | "Add a brief description..." |
| `{{.Slug}}` | Talk slug | "kubernetes-scaling" |

### 2. Creating a Talk

When you run:
```bash
make create-talk DATE=2025-11-15 SLUG=kubernetes-scaling
```

The system:
1. Converts slug to title: `kubernetes-scaling` → `Kubernetes Scaling`
2. Loads templates from `cmd/create-talk/templates/`
3. Renders each template with the variables
4. Writes files to the talk directory

## 📝 Template Examples

### metadata.yaml.tmpl

```yaml
title: "{{.Title}}"
date: "{{.Date}}"
event: "{{.Event}}"
topics:
  - Topic1
  - Topic2
description: "{{.Description}}"
```

### README.md.tmpl

```markdown
# {{.Title}}

## 📅 Talk Information

- **Date**: {{.Date}}
- **Event**: {{.Event}}

## 📝 Description

{{.Description}}
```

## 🎨 Customizing Templates

### Editing Existing Templates

1. Edit the template file in `cmd/create-talk/templates/`
2. Use template syntax for variables: `{{.Variable}}`
3. Rebuild: `make build`
4. Test by creating a new talk: `make create-talk DATE=2025-12-01 SLUG=test`

### Adding New Sections

Want to add a new section to all talks? Just edit the template:

```markdown
## 🆕 New Section

This will appear in all future talks!

## 🔗 Related Talks

- [Previous talk about {{.Slug}}](../previous-talk/)
```

### Conditional Content

Use template conditionals:

```markdown
{{if .VideoURL}}
## 🎥 Recording

Watch the recording: [{{.VideoURL}}]({{.VideoURL}})
{{else}}
## 🎥 Recording

Recording will be available after the talk.
{{end}}
```

### Loops

Add dynamic lists:

```markdown
## 📚 Topics Covered

{{range .Topics}}
- {{.}}
{{end}}
```

## 🚀 Advanced Usage

### Adding More Variables

Edit `cmd/create-talk/main.go` to add more context variables:

```go
type TalkData struct {
    Title       string
    Date        string
    Event       string
    Description string
    Slug        string
    Author      string  // New variable
    Twitter     string  // New variable
}

data := TalkData{
    Title:       title,
    Date:        dateStr,
    Event:       "Conference/Meetup Name",
    Description: "Add a brief description...",
    Slug:        slug,
    Author:      "Shanky",
    Twitter:     "@shankyjs",
}
```

Then use in templates:

```markdown
**Author**: {{.Author}} ({{.Twitter}})
```

### Custom Functions

Templates support custom functions:

```go
funcMap := template.FuncMap{
    "upper": strings.ToUpper,
    "lower": strings.ToLower,
}

tmpl, err := template.New("readme").Funcs(funcMap).ParseFiles(templatePath)
```

Use in template:

```markdown
# {{upper .Title}}
```

## 🎯 Template Best Practices

1. **Keep It Consistent** - All talks should have the same structure
2. **Use Placeholders** - Make it obvious what needs to be filled in
3. **Include Examples** - Show users what good content looks like
4. **Stay DRY** - Use templates instead of copy-paste
5. **Version Control** - Templates are in git, so changes are tracked
6. **Rebuild After Changes** - Run `make build` after editing templates

## 💡 Example Workflow

```bash
# 1. Customize the template
vim cmd/create-talk/templates/README.md.tmpl

# 2. Rebuild tools
make build

# 3. Create a talk to test
make create-talk DATE=2025-12-01 SLUG=test-template

# 4. Check the generated file
cat 2025/dec-1st-test-template/README.md

# 5. If good, keep it. If not, edit and repeat
rm -rf 2025/dec-1st-test-template
```
