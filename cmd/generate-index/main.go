package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Event       string   `yaml:"event"`
	Topics      []string `yaml:"topics"`
	Description string   `yaml:"description"`
	SlidesURL   string   `yaml:"slides_url"`
	VideoURL    string   `yaml:"video_url"`
}

type Talk struct {
	Metadata
	Path string
	Year string
}

func main() {
	fmt.Println("🔍 Scanning for talks...")

	talks, err := findAllTalks()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📚 Found %d talks\n", len(talks))

	// Update English README
	if err := updateReadme("README.md", talks, "en"); err != nil {
		fmt.Printf("❌ Error updating README.md: %v\n", err)
		os.Exit(1)
	}

	// Update Spanish README
	if err := updateReadme("docs/README-es.md", talks, "es"); err != nil {
		fmt.Printf("❌ Error updating README-es.md: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✨ Index generation complete!")
}

func findAllTalks() ([]Talk, error) {
	var talks []Talk

	// Find all year directories
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if directory name is a year (4 digits)
		if len(entry.Name()) == 4 {
			if _, err := time.Parse("2006", entry.Name()); err == nil {
				yearTalks, err := scanYearDirectory(entry.Name())
				if err != nil {
					return nil, err
				}
				talks = append(talks, yearTalks...)
			}
		}
	}

	// Sort talks by date (newest first)
	sort.Slice(talks, func(i, j int) bool {
		return talks[i].Date > talks[j].Date
	})

	return talks, nil
}

func scanYearDirectory(year string) ([]Talk, error) {
	var talks []Talk

	entries, err := os.ReadDir(year)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metadataPath := filepath.Join(year, entry.Name(), "metadata.yaml")
		if _, err := os.Stat(metadataPath); err == nil {
			data, err := os.ReadFile(metadataPath)
			if err != nil {
				continue
			}

			var meta Metadata
			if err := yaml.Unmarshal(data, &meta); err != nil {
				continue
			}

			talks = append(talks, Talk{
				Metadata: meta,
				Path:     filepath.Join(year, entry.Name()),
				Year:     year,
			})
		}
	}

	return talks, nil
}

func updateReadme(readmePath string, talks []Talk, lang string) error {
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Find index section
	var indexMarker string
	if lang == "es" {
		indexMarker = "## 📑 Índice de Charlas"
	} else {
		indexMarker = "## 📑 Talks Index"
	}

	indexStart := strings.Index(contentStr, indexMarker)
	if indexStart == -1 {
		return fmt.Errorf("could not find index section")
	}

	// Find end section
	var endMarker string
	if lang == "es" {
		endMarker = "## 🤝 Contribuir"
	} else {
		endMarker = "## 🤝 Contributing"
	}

	endStart := strings.Index(contentStr[indexStart:], endMarker)
	if endStart == -1 {
		return fmt.Errorf("could not find end section")
	}
	endStart += indexStart

	// Generate new index
	newIndex := generateIndex(talks, lang)

	// Generate statistics
	stats := generateStats(talks, lang)

	// Remove old stats if exists
	statsMarker := "## 📊"
	oldStatsStart := strings.Index(contentStr[:indexStart], statsMarker)
	if oldStatsStart != -1 {
		oldStatsEnd := strings.Index(contentStr[oldStatsStart+5:], "\n## ")
		if oldStatsEnd != -1 {
			oldStatsEnd += oldStatsStart + 5 + 1
			contentStr = contentStr[:oldStatsStart] + contentStr[oldStatsEnd:]
			// Recalculate positions
			indexStart = strings.Index(contentStr, indexMarker)
			endStart = strings.Index(contentStr[indexStart:], endMarker) + indexStart
		}
	}

	// Build new content
	newContent := contentStr[:indexStart] + stats + newIndex + contentStr[endStart:]

	if err := os.WriteFile(readmePath, []byte(newContent), 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Updated %s\n", readmePath)
	return nil
}

func generateStats(talks []Talk, lang string) string {
	if len(talks) == 0 {
		return ""
	}

	// Calculate statistics
	totalTalks := len(talks)
	pastTalks := 0
	futureTalks := 0
	today := time.Now().Format("2006-01-02")

	topicCount := make(map[string]int)
	yearCount := make(map[string]int)

	for _, talk := range talks {
		if talk.Date < today {
			pastTalks++
		} else {
			futureTalks++
		}

		yearCount[talk.Year]++

		for _, topic := range talk.Topics {
			if topic != "" && topic != "Topic1" && topic != "Topic2" && topic != "Topic3" {
				topicCount[topic]++
			}
		}
	}

	var sb strings.Builder

	if lang == "es" {
		sb.WriteString("## 📊 Estadísticas\n\n")
		sb.WriteString(fmt.Sprintf("- 🎤 **Total de Charlas**: %d\n", totalTalks))
		sb.WriteString(fmt.Sprintf("- ✅ **Pasadas**: %d\n", pastTalks))
		sb.WriteString(fmt.Sprintf("- 🔜 **Próximas**: %d\n", futureTalks))

		if len(yearCount) > 1 {
			sb.WriteString(fmt.Sprintf("- 📅 **Años Activos**: %d\n", len(yearCount)))
		}

		if len(topicCount) > 0 {
			topTopics := getTopN(topicCount, 3)
			sb.WriteString("- 🏷️ **Temas Principales**: ")
			sb.WriteString(formatTopics(topTopics))
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("## 📊 Statistics\n\n")
		sb.WriteString(fmt.Sprintf("- 🎤 **Total Talks**: %d\n", totalTalks))
		sb.WriteString(fmt.Sprintf("- ✅ **Past**: %d\n", pastTalks))
		sb.WriteString(fmt.Sprintf("- 🔜 **Upcoming**: %d\n", futureTalks))

		if len(yearCount) > 1 {
			sb.WriteString(fmt.Sprintf("- 📅 **Active Years**: %d\n", len(yearCount)))
		}

		if len(topicCount) > 0 {
			topTopics := getTopN(topicCount, 3)
			sb.WriteString("- 🏷️ **Top Topics**: ")
			sb.WriteString(formatTopics(topTopics))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

type topicCount struct {
	Topic string
	Count int
}

func getTopN(m map[string]int, n int) []topicCount {
	var sorted []topicCount
	for k, v := range m {
		sorted = append(sorted, topicCount{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		// Primary sort: by count (descending)
		if sorted[i].Count != sorted[j].Count {
			return sorted[i].Count > sorted[j].Count
		}
		// Secondary sort: by name (alphabetically) for deterministic output
		return sorted[i].Topic < sorted[j].Topic
	})

	if len(sorted) > n {
		sorted = sorted[:n]
	}

	return sorted
}

func formatTopics(topics []topicCount) string {
	var parts []string
	for _, t := range topics {
		parts = append(parts, fmt.Sprintf("%s (%d)", t.Topic, t.Count))
	}
	return strings.Join(parts, ", ")
}

func generateIndex(talks []Talk, lang string) string {
	var sb strings.Builder

	if lang == "es" {
		sb.WriteString("## 📑 Índice de Charlas\n\n")
		sb.WriteString("Explora todas las charlas por año, tema y evento. Haz clic en cualquier charla para acceder a la demo completa, código y materiales.\n\n")
	} else {
		sb.WriteString("## 📑 Talks Index\n\n")
		sb.WriteString("Browse all talks by year, topic, and event. Click on any talk to access the full demo, code, and materials.\n\n")
	}

	// Group by year
	talksByYear := make(map[string][]Talk)
	for _, talk := range talks {
		talksByYear[talk.Year] = append(talksByYear[talk.Year], talk)
	}

	// Sort years descending
	var years []string
	for year := range talksByYear {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(years)))

	// Generate tables by year
	for _, year := range years {
		sb.WriteString(fmt.Sprintf("### %s\n\n", year))
		sb.WriteString(generateTable(talksByYear[year], lang))
		sb.WriteString("\n\n")
	}

	// Coming soon section
	if lang == "es" {
		sb.WriteString("### Próximamente 🚀\n\n")
		sb.WriteString("¡Más charlas y demos se agregarán aquí a medida que sucedan!\n\n")
		sb.WriteString("---\n\n")
		sb.WriteString("## 🏷️ Buscar por Tema\n\n")
	} else {
		sb.WriteString("### Coming Soon 🚀\n\n")
		sb.WriteString("More talks and demos will be added here as they happen!\n\n")
		sb.WriteString("---\n\n")
		sb.WriteString("## 🏷️ Browse by Topic\n\n")
	}

	sb.WriteString(generateTopicsIndex(talks, lang))
	sb.WriteString("\n\n")

	return sb.String()
}

func generateTable(talks []Talk, lang string) string {
	var sb strings.Builder

	if lang == "es" {
		sb.WriteString("| Fecha | Título de la Charla | Temas | Evento/Ubicación | Materiales |\n")
		sb.WriteString("|-------|---------------------|-------|------------------|------------|\n")
	} else {
		sb.WriteString("| Date | Talk Title | Topics | Event/Location | Materials |\n")
		sb.WriteString("|------|------------|--------|----------------|-----------|\n")
	}

	for _, talk := range talks {
		date := talk.Date
		title := fmt.Sprintf("[**%s**](./%s)", talk.Title, talk.Path)
		topics := strings.Join(talk.Topics, ", ")
		event := talk.Event
		materials := fmt.Sprintf("[EN](./%s/README.md) / [ES](./%s/README-es.md)", talk.Path, talk.Path)

		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", date, title, topics, event, materials))
	}

	return sb.String()
}

func generateTopicsIndex(talks []Talk, lang string) string {
	topicsMap := make(map[string][]Talk)

	for _, talk := range talks {
		for _, topic := range talk.Topics {
			topicsMap[topic] = append(topicsMap[topic], talk)
		}
	}

	var topics []string
	for topic := range topicsMap {
		topics = append(topics, topic)
	}
	sort.Strings(topics)

	var sb strings.Builder
	for _, topic := range topics {
		talks := topicsMap[topic]
		var links []string
		for _, talk := range talks {
			link := fmt.Sprintf("[%s (%s)](./%s)", talk.Title, talk.Year, talk.Path)
			links = append(links, link)
		}
		sb.WriteString(fmt.Sprintf("- **%s**: %s\n", topic, strings.Join(links, ", ")))
	}

	return sb.String()
}
