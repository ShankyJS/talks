package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Event       string   `yaml:"event"`
	Topics      []string `yaml:"topics"`
	Description string   `yaml:"description"`
}

func main() {
	errors := []string{}
	warnings := []string{}

	// Find all year directories
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if directory name is a year
		if len(entry.Name()) == 4 {
			if _, err := time.Parse("2006", entry.Name()); err == nil {
				checkYear(entry.Name(), &errors, &warnings)
			}
		}
	}

	// Print results
	if len(warnings) > 0 {
		fmt.Println("\n⚠️  Warnings:")
		for _, w := range warnings {
			fmt.Printf("  %s\n", w)
		}
	}

	if len(errors) > 0 {
		fmt.Println("\n❌ Errors:")
		for _, e := range errors {
			fmt.Printf("  %s\n", e)
		}
		fmt.Println("\nPlease fix the errors above.")
		os.Exit(1)
	}

	if len(warnings) == 0 && len(errors) == 0 {
		fmt.Println("✅ All talk directories have valid metadata!")
	}
}

func checkYear(year string, errors, warnings *[]string) {
	talks, err := os.ReadDir(year)
	if err != nil {
		return
	}

	for _, talk := range talks {
		if !talk.IsDir() {
			continue
		}

		if talk.Name()[0] == '.' {
			continue
		}

		talkPath := filepath.Join(year, talk.Name())
		metadataPath := filepath.Join(talkPath, "metadata.yaml")

		// Check if metadata exists
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			*errors = append(*errors, fmt.Sprintf("❌ Missing metadata.yaml: %s", talkPath))
			continue
		}

		// Check if READMEs exist
		readmeEN := filepath.Join(talkPath, "README.md")
		readmeES := filepath.Join(talkPath, "README-es.md")

		if _, err := os.Stat(readmeEN); os.IsNotExist(err) {
			*warnings = append(*warnings, fmt.Sprintf("⚠️  Missing README.md: %s", talkPath))
		}

		if _, err := os.Stat(readmeES); os.IsNotExist(err) {
			*warnings = append(*warnings, fmt.Sprintf("⚠️  Missing README-es.md: %s", talkPath))
		}

		// Validate metadata content
		data, err := os.ReadFile(metadataPath)
		if err != nil {
			*errors = append(*errors, fmt.Sprintf("❌ Error reading %s: %v", metadataPath, err))
			continue
		}

		var meta Metadata
		if err := yaml.Unmarshal(data, &meta); err != nil {
			*errors = append(*errors, fmt.Sprintf("❌ Error parsing %s: %v", metadataPath, err))
			continue
		}

		// Check required fields
		if meta.Title == "" {
			*errors = append(*errors, fmt.Sprintf("❌ Missing required field 'title' in %s", metadataPath))
		}
		if meta.Date == "" {
			*errors = append(*errors, fmt.Sprintf("❌ Missing required field 'date' in %s", metadataPath))
		}
		if len(meta.Topics) == 0 {
			*errors = append(*errors, fmt.Sprintf("❌ Missing required field 'topics' in %s", metadataPath))
		}
	}
}
