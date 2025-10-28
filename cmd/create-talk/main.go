package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/*
var templatesFS embed.FS

type TalkData struct {
	Title       string
	Date        string
	Event       string
	Description string
	Slug        string
}

func main() {
	dateFlag := flag.String("date", "", "Talk date in YYYY-MM-DD format")
	slugFlag := flag.String("slug", "", "Talk slug for directory name")
	titleFlag := flag.String("title", "", "Talk title (optional, auto-generated from slug)")

	flag.Parse()

	if *dateFlag == "" || *slugFlag == "" {
		fmt.Println("❌ Error: DATE and SLUG are required")
		fmt.Println("Usage: create-talk -date 2025-11-15 -slug kubernetes-scaling")
		os.Exit(1)
	}

	if err := createTalk(*dateFlag, *slugFlag, *titleFlag); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func createTalk(dateStr, slug, title string) error {
	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format '%s'. Use YYYY-MM-DD", dateStr)
	}

	// Generate title from slug if not provided
	if title == "" {
		title = strings.Title(strings.ReplaceAll(strings.ReplaceAll(slug, "-", " "), "_", " "))
	}

	// Create directory name
	year := date.Format("2006")
	monthDay := formatMonthDay(date)
	dirName := fmt.Sprintf("%s-%s", monthDay, slug)
	fullPath := filepath.Join(year, dirName)

	// Check if directory exists
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("directory already exists: %s", fullPath)
	}

	// Create directory
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("📁 Creating new talk directory: %s\n", fullPath)

	// Prepare template data
	data := TalkData{
		Title:       title,
		Date:        dateStr,
		Event:       "Conference/Meetup Name",
		Description: "Add a brief description of your talk here",
		Slug:        slug,
	}

	// Create files from templates
	files := []string{"metadata.yaml.tmpl", "README.md.tmpl", "README-es.md.tmpl"}
	outputs := []string{"metadata.yaml", "README.md", "README-es.md"}

	for i, tmplFile := range files {
		if err := renderTemplate(tmplFile, filepath.Join(fullPath, outputs[i]), data); err != nil {
			return fmt.Errorf("failed to create %s: %w", outputs[i], err)
		}
	}

	fmt.Println("✅ Talk directory created successfully!")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Printf("  1. Edit %s/metadata.yaml with your talk details\n", fullPath)
	fmt.Printf("  2. Update %s/README.md with your content\n", fullPath)
	fmt.Printf("  3. Update %s/README-es.md with Spanish content\n", fullPath)
	fmt.Println("  4. Run 'make update-index' to regenerate the talks index")
	fmt.Println("")
	fmt.Println("📝 Files created:")
	fmt.Printf("  - %s/metadata.yaml\n", fullPath)
	fmt.Printf("  - %s/README.md\n", fullPath)
	fmt.Printf("  - %s/README-es.md\n", fullPath)

	return nil
}

func formatMonthDay(date time.Time) string {
	month := strings.ToLower(date.Format("Jan"))
	day := date.Day()

	suffix := "th"
	if day%100 < 11 || day%100 > 13 {
		switch day % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}

	return fmt.Sprintf("%s-%d%s", month, day, suffix)
}

func renderTemplate(tmplFile, outputFile string, data TalkData) error {
	tmplContent, err := templatesFS.ReadFile("templates/" + tmplFile)
	if err != nil {
		return err
	}

	tmpl, err := template.New(tmplFile).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
