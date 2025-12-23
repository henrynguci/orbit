package tui

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func getLastModifiedTime(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "Unknown"
	}

	modTime := info.ModTime()
	now := time.Now()

	if modTime.Year() == now.Year() && modTime.YearDay() == now.YearDay() {
		return fmt.Sprintf("Today %02d:%02d", modTime.Hour(), modTime.Minute())
	}

	if modTime.Year() == now.Year() && modTime.YearDay() == now.YearDay()-1 {
		return fmt.Sprintf("Yesterday %02d:%02d", modTime.Hour(), modTime.Minute())
	}

	return modTime.Format("02/01/2006 15:04")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncateString(s string, maxLen int) string {
	if strings.Contains(s, "\x1b[") {
		return s
	}
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

type dashboardRow struct {
	Workspace string
	Project   string
	Status    string
	Path      string
}
