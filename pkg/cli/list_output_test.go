package cli_test

import "fmt"

func formatListLine(id string, priority int, status, title, parents, depends string) string {
	return fmt.Sprintf("%-9s  %1d  %-4s  %-60s  %-13s  %-13s\n", id, priority, status, title, parents, depends)
}

func truncateListTitleForTest(title string) string {
	const maxTitleLength = 60
	const ellipsis = "..."

	runes := []rune(title)
	if len(runes) <= maxTitleLength {
		return title
	}

	return string(runes[:maxTitleLength-len([]rune(ellipsis))]) + ellipsis
}
