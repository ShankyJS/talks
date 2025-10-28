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
}

type Talk struct {
	Metadata
	Path string
	Year string
}

func main() {
	talks, err := findAllTalks()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	output := generateStats(talks)
	fmt.Print(output)

	// Save to file
	if err := os.WriteFile("stats.txt", []byte(output), 0644); err != nil {
		fmt.Printf("❌ Error writing stats.txt: %v\n", err)
		os.Exit(1)
	}
}

func findAllTalks() ([]Talk, error) {
	var talks []Talk

	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

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

func generateStats(talks []Talk) string {
	var sb strings.Builder

	if len(talks) == 0 {
		sb.WriteString("## 📊 Talk Statistics\n\n")
		sb.WriteString("No talks found yet. Create your first talk with:\n")
		sb.WriteString("```bash\n")
		sb.WriteString("make create-talk DATE=2025-12-01 SLUG=my-first-talk\n")
		sb.WriteString("```\n")
		return sb.String()
	}

	// Calculate statistics
	totalTalks := len(talks)
	pastTalks := 0
	futureTalks := 0
	today := time.Now().Format("2006-01-02")

	talksByYear := make(map[string]int)
	topicCount := make(map[string]int)
	eventCount := make(map[string]int)

	var upcoming []Talk

	for _, talk := range talks {
		if talk.Date < today {
			pastTalks++
		} else {
			futureTalks++
			upcoming = append(upcoming, talk)
		}

		talksByYear[talk.Year]++

		for _, topic := range talk.Topics {
			if topic != "" {
				topicCount[topic]++
			}
		}

		if talk.Event != "" && talk.Event != "Conference/Meetup Name" && talk.Event != "Unknown" {
			eventCount[talk.Event]++
		}
	}

	// Sort upcoming talks by date
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].Date < upcoming[j].Date
	})

	// Generate output
	sb.WriteString("## 📊 Talk Statistics\n\n")
	sb.WriteString(fmt.Sprintf("### 🎤 Total Talks: %d\n\n", totalTalks))
	sb.WriteString(fmt.Sprintf("- **Past Talks**: %d\n", pastTalks))
	sb.WriteString(fmt.Sprintf("- **Upcoming Talks**: %d\n\n", futureTalks))

	// Talks by year
	sb.WriteString("### 📅 Talks by Year\n\n")
	years := make([]string, 0, len(talksByYear))
	for year := range talksByYear {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(years)))

	for _, year := range years {
		count := talksByYear[year]
		bar := strings.Repeat("█", count)
		sb.WriteString(fmt.Sprintf("- **%s**: %d %s\n", year, count, bar))
	}
	sb.WriteString("\n")

	// Top topics
	sb.WriteString("### 🏷️ Most Popular Topics\n\n")
	topics := topNTopics(topicCount, 10)
	for _, t := range topics {
		if t.Topic != "" {
			bar := strings.Repeat("█", t.Count)
			sb.WriteString(fmt.Sprintf("- **%s**: %d %s\n", t.Topic, t.Count, bar))
		}
	}
	sb.WriteString("\n")

	// Events
	sb.WriteString("### 🎪 Events\n\n")
	if len(eventCount) > 0 {
		events := topNTopics(eventCount, 5)
		for _, e := range events {
			sb.WriteString(fmt.Sprintf("- **%s**: %d talks\n", e.Topic, e.Count))
		}
	} else {
		sb.WriteString("No events with talks yet.\n")
	}
	sb.WriteString("\n")

	// Upcoming talks
	if len(upcoming) > 0 {
		sb.WriteString("### 🔜 Upcoming Talks\n\n")
		for i, talk := range upcoming {
			if i >= 5 {
				break
			}
			sb.WriteString(fmt.Sprintf("- **%s**: %s @ %s\n", talk.Date, talk.Title, talk.Event))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

type topicCount struct {
	Topic string
	Count int
}

func topNTopics(m map[string]int, n int) []topicCount {
	var sorted []topicCount
	for k, v := range m {
		sorted = append(sorted, topicCount{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	if len(sorted) > n {
		sorted = sorted[:n]
	}

	return sorted
}
