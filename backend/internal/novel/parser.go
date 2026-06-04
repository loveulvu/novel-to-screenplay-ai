package novel

import (
	"regexp"
	"strings"
)

var chapterHeadingPattern = regexp.MustCompile(`(?i)^\s*((第\s*([0-9]+|[一二三四五六七八九十百零〇两]+)\s*章)|(chapter\s+([0-9]+)))(?:\s+|[:：-]\s*|$)(.*)$`)

func ParseChapters(input string) []Chapter {
	lines := strings.Split(input, "\n")
	var chapters []Chapter
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
		matches := chapterHeadingPattern.FindStringSubmatch(line)
		if matches != nil {
			flush()
			number := len(chapters) + 1
			title := strings.TrimSpace(matches[6])
			if title == "" {
				title = strings.TrimSpace(matches[1])
			}
			current = &Chapter{
				Number: number,
				Title:  title,
			}
			continue
		}

		if current == nil {
			if strings.TrimSpace(line) == "" {
				continue
			}
			current = &Chapter{
				Number: 1,
				Title:  "未命名章节",
			}
		}
		body = append(body, line)
	}

	flush()
	return chapters
}
