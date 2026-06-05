package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"novel-to-screenplay-ai/internal/analysis"
	"novel-to-screenplay-ai/internal/novel"
	"novel-to-screenplay-ai/internal/screenplay"
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

func BuildRepairJSONPrompt(originalPrompt string, rawOutput string, parseError error, schema string) string {
	return fmt.Sprintf(`你上一次返回的内容无法解析为目标 Go struct JSON。请基于原始任务、解析错误和目标字段要求，返回修复后的合法 JSON。

要求：
- 只返回修复后的 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 字段名必须符合 json tag
- 不要省略必填字段
- 数组字段必须输出 JSON 数组，不能输出字符串

目标字段要求：
%s

json.Unmarshal 错误：
%s

原始任务：
%s

模型第一次返回内容：
%s`, schema, parseError.Error(), truncateText(originalPrompt, repairPromptLimit), truncateText(rawOutput, repairPromptLimit))
}

func BuildRepairScreenplayPrompt(originalPrompt string, current screenplay.Screenplay, validationErrors []string) string {
	payload, _ := json.MarshalIndent(current, "", "  ")
	return fmt.Sprintf(`你上一次返回的 screenplay.Screenplay JSON 已经可以解析，但没有通过后端 Validate。请基于原始任务、当前 JSON 和校验错误，返回修复后的 screenplay.Screenplay JSON。

要求：
- 只返回修复后的 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 字段名必须符合 json tag
- 不要省略必填字段
- 不要删除 validator 需要的字段
- source_chapters 必须非空，每项 number > 0 且 title 非空
- characters 必须非空，每项 id、name、role 非空
- scenes 必须非空，每项 id、location、time、summary 非空
- 每个 scene.characters 必须非空
- 每个 scene.dialogues 必须非空
- 每条 dialogue.character 和 dialogue.line 必须非空

校验错误：
%s

原始任务：
%s

当前 screenplay JSON：
%s`, strings.Join(validationErrors, "\n"), truncateText(originalPrompt, repairPromptLimit), truncateText(string(payload), repairPromptLimit))
}

func chapterAnalysisSchemaDescription() string {
	return `analysis.ChapterAnalysis:
{
  "chapter_number": number,
  "chapter_title": string,
  "summary": string,
  "characters": [
    {
      "name": string,
      "role_in_chapter": string,
      "traits": [string],
      "state_change": string
    }
  ],
  "locations": [string],
  "key_events": [string],
  "conflicts": [string],
  "scene_candidates": [
    {
      "location": string,
      "time": string,
      "purpose": string,
      "characters": [string],
      "key_events": [string]
    }
  ]
}`
}

func storyBibleSchemaDescription() string {
	return `story.StoryBible:
{
  "title": string,
  "logline": string,
  "global_characters": [
    {
      "id": string,
      "name": string,
      "role": string,
      "motivation": string
    }
  ],
  "timeline": [
    {
      "chapter_number": number,
      "event": string
    }
  ],
  "main_conflict": string,
  "scene_plan": [
    {
      "id": string,
      "source_chapter": number,
      "summary": string,
      "location": string,
      "time": string,
      "characters": [string]
    }
  ]
}`
}

func screenplaySchemaDescription() string {
	return `screenplay.Screenplay:
{
  "title": string,
  "source_chapters": [
    {
      "number": number,
      "title": string,
      "summary": string
    }
  ],
  "characters": [
    {
      "id": string,
      "name": string,
      "role": string,
      "description": string
    }
  ],
  "scenes": [
    {
      "id": string,
      "source_chapter": number,
      "location": string,
      "time": string,
      "summary": string,
      "characters": [string],
      "dialogues": [
        {
          "character": string,
          "emotion": string,
          "line": string
        }
      ],
      "actions": [string]
    }
  ]
}`
}
