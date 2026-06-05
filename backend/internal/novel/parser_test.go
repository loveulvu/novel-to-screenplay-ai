package novel

import "testing"

func TestParseChaptersStandardChineseHeadings(t *testing.T) {
	input := `第1章 开端
这里是开端正文
第二章 发展
这里是发展正文
第十二章 雨夜
这里是雨夜正文`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "开端", "这里是开端正文")
	assertChapter(t, chapters[1], 2, "发展", "这里是发展正文")
	assertChapter(t, chapters[2], 3, "雨夜", "这里是雨夜正文")
}

func TestParseChaptersChineseHeadingsWithoutSpaces(t *testing.T) {
	input := `第1章开端
这里是开端正文
第一章发展
这里是发展正文
第十二章雨夜
这里是雨夜正文`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "开端", "这里是开端正文")
	assertChapter(t, chapters[1], 2, "发展", "这里是发展正文")
	assertChapter(t, chapters[2], 3, "雨夜", "这里是雨夜正文")
}

func TestParseChaptersChineseSectionHeadings(t *testing.T) {
	input := `第一节：相遇
这里是相遇正文
第二节 追踪
这里是追踪正文
第三节对峙
这里是对峙正文`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "相遇", "这里是相遇正文")
	assertChapter(t, chapters[1], 2, "追踪", "这里是追踪正文")
	assertChapter(t, chapters[2], 3, "对峙", "这里是对峙正文")
}

func TestParseChaptersEnglishHeadings(t *testing.T) {
	input := `Chapter 1 Beginning
Chapter one body
chapter 2 Middle
Chapter two body
Chapter 3: Ending
Chapter three body`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "Beginning", "Chapter one body")
	assertChapter(t, chapters[1], 2, "Middle", "Chapter two body")
	assertChapter(t, chapters[2], 3, "Ending", "Chapter three body")
}

func TestParseChaptersIgnoresPreface(t *testing.T) {
	input := `前言
这是一段简介，不应该变成未命名章节。

第1章 开端
这里是开端正文
第二章 发展
这里是发展正文
第三章 结尾
这里是结尾正文`

	chapters := ParseChapters(input)

	if len(chapters) != 3 {
		t.Fatalf("expected 3 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "开端", "这里是开端正文")
	assertChapter(t, chapters[1], 2, "发展", "这里是发展正文")
	assertChapter(t, chapters[2], 3, "结尾", "这里是结尾正文")
}

func TestParseChaptersLessThanThreeChapters(t *testing.T) {
	input := `第1章 开端
这里是开端正文
第二章 发展
这里是发展正文`

	chapters := ParseChapters(input)

	if len(chapters) != 2 {
		t.Fatalf("expected 2 chapters, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "开端", "这里是开端正文")
	assertChapter(t, chapters[1], 2, "发展", "这里是发展正文")
}

func TestParseChaptersNoHeadingsReturnsEmptySlice(t *testing.T) {
	input := `前言
这里只是简介。
没有任何正式章节标题。`

	chapters := ParseChapters(input)

	if len(chapters) != 0 {
		t.Fatalf("expected empty chapter slice, got %d", len(chapters))
	}
}

func TestParseChaptersMergesPaginatedSameSection(t *testing.T) {
	input := `第一节：雨夜来信 (第1/2页)
第一页正文
第一节：雨夜来信 (第2/2页)
第二页正文`

	chapters := ParseChapters(input)

	if len(chapters) != 1 {
		t.Fatalf("expected 1 chapter, got %d", len(chapters))
	}
	assertChapter(t, chapters[0], 1, "雨夜来信", "第一页正文\n第二页正文")
}

func assertChapter(t *testing.T, chapter Chapter, number int, title string, text string) {
	t.Helper()

	if chapter.Number != number {
		t.Fatalf("expected chapter number %d, got %d", number, chapter.Number)
	}
	if chapter.Title != title {
		t.Fatalf("expected chapter title %q, got %q", title, chapter.Title)
	}
	if chapter.Text != text {
		t.Fatalf("expected chapter text %q, got %q", text, chapter.Text)
	}
}
