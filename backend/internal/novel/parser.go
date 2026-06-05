package novel

import (
	"regexp"
	"strings"
)

var (
	chineseChapterHeadingPattern = regexp.MustCompile(`^第([0-9]+|[一二三四五六七八九十百千万零〇两]+)[章节][:：-]?(.*)$`)
	englishChapterHeadingPattern = regexp.MustCompile(`(?i)^chapter\s+([0-9]+)(?:\s+|[:：]\s*|$)(.*)$`)
	pageMarkPattern              = regexp.MustCompile(`\s*\(第\d+/\d+页\)\s*$`)
)

func ParseChapters(input string) []Chapter {
	lines := strings.Split(input, "\n")
	chapters := make([]Chapter, 0)
	var current *Chapter
	var body []string

	flush := func() {
		if current == nil {
			return
		}
		current.Text = strings.TrimSpace(strings.Join(body, "\n"))
		chapters = append(chapters, *current)
		body = nil
	}

	for _, line := range lines {
		title, ok := parseChapterHeading(line)
		if ok {
			flush()
			current = &Chapter{
				Number: len(chapters) + 1,
				Title:  title,
			}
			continue
		}

		if current == nil {
			continue
		}
		body = append(body, line)
	}

	flush()
	return chapters
}

func parseChapterHeading(line string) (string, bool) {
	if title, ok := parseChineseChapterHeading(line); ok {
		return title, true
	}
	return parseEnglishChapterHeading(line)
}

func parseChineseChapterHeading(line string) (string, bool) {
	line = normalizeHeadingLine(line)

	matches := chineseChapterHeadingPattern.FindStringSubmatch(line)
	if matches == nil {
		return "", false
	}

	title := strings.TrimSpace(matches[2])
	title = pageMarkPattern.ReplaceAllString(title, "")
	title = strings.TrimSpace(title)

	if looksLikeBodyText(title) {
		return "", false
	}
	if title == "" {
		title = strings.TrimSpace(matches[0])
	}
	return title, true
}
func normalizeHeadingLine(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "\uFEFF")
	line = strings.ReplaceAll(line, "\u3000", "")
	line = strings.ReplaceAll(line, " ", "")
	return line
}

func parseEnglishChapterHeading(line string) (string, bool) {
	matches := englishChapterHeadingPattern.FindStringSubmatch(line)
	if matches == nil {
		return "", false
	}

	title := strings.TrimSpace(matches[2])
	if looksLikeBodyText(title) {
		return "", false
	}
	if title == "" {
		title = strings.TrimSpace(matches[0])
	}
	return title, true
}

func looksLikeBodyText(title string) bool {
	return strings.ContainsAny(title, "，。！？；,!?;")
}
