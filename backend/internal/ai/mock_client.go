package ai

import (
	"context"
	"fmt"

	"novel-to-screenplay-ai/internal/analysis"
	"novel-to-screenplay-ai/internal/novel"
	"novel-to-screenplay-ai/internal/screenplay"
	"novel-to-screenplay-ai/internal/story"
)

type MockClient struct{}

func NewMockClient() MockClient {
	return MockClient{}
}

func (m MockClient) AnalyzeChapter(ctx context.Context, chapter novel.Chapter) (analysis.ChapterAnalysis, error) {
	_ = ctx

	summaryByChapter := map[int]string{
		1: "林澈在雨夜收到失踪父亲留下的铜钥匙，并发现旧剧院的地下室仍在运转。",
		2: "许岚陪林澈进入旧剧院，二人在尘封舞台上发现父亲研究的时钟装置。",
		3: "反派顾衡现身索要钥匙，林澈必须决定是否启动时钟装置寻找真相。",
	}
	summary := summaryByChapter[chapter.Number]
	if summary == "" {
		summary = fmt.Sprintf("第%d章推进了钥匙、剧院和父亲失踪之谜。", chapter.Number)
	}

	return analysis.ChapterAnalysis{
		ChapterNumber: chapter.Number,
		ChapterTitle:  chapter.Title,
		Summary:       summary,
		Characters:    []string{"林澈", "许岚", "顾衡"},
		Locations:     []string{"旧剧院", "雨巷", "地下室"},
		KeyEvents: []string{
			fmt.Sprintf("第%d章揭示新的线索。", chapter.Number),
			"铜钥匙与旧剧院产生关联。",
		},
		Conflicts: []string{
			"林澈想查清父亲失踪真相，顾衡试图阻止他。",
		},
		SceneCandidates: []string{
			"雨夜收到钥匙",
			"旧剧院探索",
			"时钟装置对峙",
		},
	}, nil
}

func (m MockClient) MergeStoryBible(ctx context.Context, analyses []analysis.ChapterAnalysis) (story.StoryBible, error) {
	_ = ctx

	timeline := make([]story.TimelineEvent, 0, len(analyses))
	for _, item := range analyses {
		event := item.Summary
		if len(item.KeyEvents) > 0 {
			event = item.KeyEvents[0]
		}
		timeline = append(timeline, story.TimelineEvent{
			ChapterNumber: item.ChapterNumber,
			Event:         event,
		})
	}

	return story.StoryBible{
		Title:   "雨夜旧剧院",
		Logline: "一个青年追查父亲失踪之谜，却在旧剧院发现能改写记忆的时钟装置。",
		GlobalCharacters: []story.Character{
			{ID: "char_lin_che", Name: "林澈", Role: "主角", Motivation: "查清父亲失踪真相"},
			{ID: "char_xu_lan", Name: "许岚", Role: "盟友", Motivation: "保护林澈并揭开剧院秘密"},
			{ID: "char_gu_heng", Name: "顾衡", Role: "阻碍者", Motivation: "夺回铜钥匙并隐藏旧实验"},
		},
		Timeline:     timeline,
		MainConflict: "林澈要公开父亲失踪真相，顾衡要封存剧院里的记忆实验。",
		ScenePlan: []story.ScenePlanItem{
			{ID: "scene_001", SourceChapter: 1, Summary: "林澈在雨巷收到父亲留下的铜钥匙。", Location: "雨巷", Time: "夜晚", Characters: []string{"char_lin_che"}},
			{ID: "scene_002", SourceChapter: 2, Summary: "林澈和许岚进入旧剧院地下室，发现时钟装置。", Location: "旧剧院地下室", Time: "深夜", Characters: []string{"char_lin_che", "char_xu_lan"}},
			{ID: "scene_003", SourceChapter: 3, Summary: "顾衡逼迫林澈交出钥匙，三人在舞台上对峙。", Location: "旧剧院舞台", Time: "黎明前", Characters: []string{"char_lin_che", "char_xu_lan", "char_gu_heng"}},
		},
	}, nil
}

func (m MockClient) GenerateScreenplay(ctx context.Context, bible story.StoryBible) (screenplay.Screenplay, error) {
	_ = ctx

	sourceChapters := make([]int, 0, len(bible.Timeline))
	for _, item := range bible.Timeline {
		sourceChapters = append(sourceChapters, item.ChapterNumber)
	}

	characters := make([]screenplay.Character, 0, len(bible.GlobalCharacters))
	for _, character := range bible.GlobalCharacters {
		characters = append(characters, screenplay.Character{
			ID:   character.ID,
			Name: character.Name,
			Role: character.Role,
		})
	}

	scenes := make([]screenplay.Scene, 0, len(bible.ScenePlan))
	for _, plan := range bible.ScenePlan {
		scenes = append(scenes, screenplay.Scene{
			ID:            plan.ID,
			SourceChapter: plan.SourceChapter,
			Location:      plan.Location,
			Time:          plan.Time,
			Summary:       plan.Summary,
			Characters:    plan.Characters,
			Dialogues: []screenplay.Dialogue{
				{Character: "char_lin_che", Emotion: "克制", Line: "这把钥匙不是遗物，它像是在等我回来。"},
				{Character: "char_xu_lan", Emotion: "紧张", Line: "如果这里真的能保存记忆，我们就不能让顾衡先找到它。"},
			},
			Actions: []string{
				"雨水顺着破损的屋檐滴落，铜钥匙在林澈掌心微微发亮。",
				"远处的舞台灯忽明忽暗，像有人刚刚经过。",
			},
		})
	}

	return screenplay.Screenplay{
		Title:          bible.Title,
		SourceChapters: sourceChapters,
		Characters:     characters,
		Scenes:         scenes,
	}, nil
}
