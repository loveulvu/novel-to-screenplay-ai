package novel

import (
	"regexp"
	"strings"
)

var (
	chineseChapterHeadingPattern = regexp.MustCompile(`^\s*第\s*([0-9０-９]+|[一二三四五六七八九十百千万零〇两]+)\s*([章节])\s*[:：\-—]?\s*(.*)$`)
	englishChapterHeadingPattern = regexp.MustCompile(`(?i)^\s*chapter\s+([0-9]+)(?:\s+|[:：]\s*|$)(.*)$`)
	pageMarkPattern              = regexp.MustCompile(`\s*[\(（]\s*第\s*[0-9０-９]+\s*[/／]\s*[0-9０-９]+\s*页\s*[\)）]\s*$`)
)

func ParseChapters(input string) []Chapter {
	lines := strings.Split(input, "\n")
	chapters := make([]Chapter, 0)
	var current *Chapter
	var currentKey string
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
		chapterKey, title, ok := parseChapterHeading(line)
		if ok {
			if current != nil && chapterKey == currentKey {
				continue
			}

			flush()
			current = &Chapter{
				Number: len(chapters) + 1,
				Title:  title,
			}
			currentKey = chapterKey
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

func parseChapterHeading(line string) (string, string, bool) {
	if chapterKey, title, ok := parseChineseChapterHeading(line); ok {
		return chapterKey, title, true
	}
	return parseEnglishChapterHeading(line)
}

func parseChineseChapterHeading(line string) (string, string, bool) {
	cleanedLine := normalizeHeadingLine(line)
	matches := chineseChapterHeadingPattern.FindStringSubmatch(cleanedLine)
	if matches == nil {
		return "", "", false
	}

	ordinal := normalizeFullWidthDigits(matches[1])
	kind := matches[2]
	title := cleanupHeadingTitle(matches[3])
	if title == "" {
		title = cleanupHeadingTitle(matches[0])
	}
	return "cn:" + ordinal + ":" + kind + ":" + title, title, true
}

func parseEnglishChapterHeading(line string) (string, string, bool) {
	matches := englishChapterHeadingPattern.FindStringSubmatch(normalizeHeadingLine(line))
	if matches == nil {
		return "", "", false
	}

	title := cleanupHeadingTitle(matches[2])
	if title == "" {
		title = cleanupHeadingTitle(matches[0])
	}
	return "en:" + matches[1] + ":chapter:" + title, title, true
}

func normalizeHeadingLine(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "\uFEFF")
	line = strings.ReplaceAll(line, "\u3000", " ")
	return line
}

func cleanupHeadingTitle(title string) string {
	title = normalizeHeadingLine(title)
	title = pageMarkPattern.ReplaceAllString(title, "")
	title = strings.TrimSpace(title)
	title = strings.Trim(title, "-—:：")
	return strings.TrimSpace(title)
}

func normalizeFullWidthDigits(value string) string {
	replacer := strings.NewReplacer(
		"０", "0",
		"１", "1",
		"２", "2",
		"３", "3",
		"４", "4",
		"５", "5",
		"６", "6",
		"７", "7",
		"８", "8",
		"９", "9",
	)
	return replacer.Replace(value)
}
