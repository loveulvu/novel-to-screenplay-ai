package ai

import (
	"encoding/json"
	"fmt"

	"novel-to-screenplay-ai/internal/analysis"
	"novel-to-screenplay-ai/internal/novel"
	"novel-to-screenplay-ai/internal/story"
)

const systemPrompt = "你是小说改编剧本的结构化分析助手。你必须只输出合法 JSON，不输出 markdown、解释、注释或 YAML。"

func BuildChapterAnalysisPrompt(chapter novel.Chapter) string {
	return fmt.Sprintf(`请根据下面的单章小说文本，输出严格符合 analysis.ChapterAnalysis 的 JSON。

必须包含字段：
- chapter_number
- chapter_title
- summary
- characters
- locations
- key_events
- conflicts
- scene_candidates

characters 每一项必须包含：
- name
- role_in_chapter
- traits
- state_change

scene_candidates 每一项必须包含：
- location
- time
- purpose
- characters
- key_events

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 不要省略字段
- 字段名必须使用上述 snake_case JSON 字段名
- traits、locations、key_events、conflicts、characters 等数组字段必须始终输出 JSON 数组，即使只有一个元素也不能输出字符串
- 如果信息不足，也要给出合理的空数组或简短说明，但不要缺字段

characters 示例：
[
  {
    "name": "角色名",
    "role_in_chapter": "本章作用",
    "traits": ["性格特征1", "性格特征2"],
    "state_change": "状态变化"
  }
]

章节编号：%d
章节标题：%s
章节正文：
%s`, chapter.Number, chapter.Title, chapter.Text)
}

func BuildStoryBiblePrompt(analyses []analysis.ChapterAnalysis) string {
	payload, _ := json.MarshalIndent(analyses, "", "  ")
	return fmt.Sprintf(`请根据所有章节分析结果，合并成 story.StoryBible JSON。

必须输出字段：
- title
- logline
- global_characters
- timeline
- main_conflict
- scene_plan

global_characters 每一项必须包含：
- id
- name
- role
- motivation

timeline 每一项必须包含：
- chapter_number
- event

scene_plan 每一项必须包含：
- id
- source_chapter
- summary
- location
- time
- characters

重点：
- 统一人物名称，避免同一人物在不同章节漂移
- 为每个主要角色生成稳定 id，建议使用 char_ 开头的小写英文/拼音/下划线
- 合并时间线
- 提炼主线冲突
- 生成适合剧本改编的 title 和 logline
- 生成至少 3 个适合剧本改编的 scene_plan

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 字段名必须符合上述 JSON 字段名
- 不要缺字段

章节分析 JSON：
%s`, string(payload))
}

func BuildScreenplayPrompt(bible story.StoryBible) string {
	payload, _ := json.MarshalIndent(bible, "", "  ")
	return fmt.Sprintf(`请根据 StoryBible 生成 screenplay.Screenplay JSON。

必须输出字段：
- title
- source_chapters
- characters
- scenes

source_chapters 每一项必须包含：
- number
- title
- summary

characters 每一项必须包含：
- id
- name
- role
- description

scenes 每一项必须包含：
- id
- source_chapter
- location
- time
- summary
- characters
- dialogues
- actions

dialogues 每一项必须包含：
- character
- emotion
- line

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- characters 字段里引用角色 id
- dialogue.character 也使用角色 id，并尽量和 characters 中的 id 一致
- 至少生成 3 个 scenes
- 每个 scene 至少有 1 条 dialogue
- 不要缺少 validator 要求的字段
- source_chapters.title 可以根据章节事件概括，但不能为空

StoryBible JSON：
%s`, string(payload))
}
