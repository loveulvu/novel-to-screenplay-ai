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

const faithfulAdaptationPrinciples = `1. 不得新增原文没有明确出现的具体事实。禁止新增原文没有的道具、动作、人物关系、地点、数量、因果解释、作弊手段、法术、物品或动机。
2. 可以做剧本化改编，但必须忠实原文事实。允许压缩情节、合并相近场景、把心理活动改写成可表演的低风险动作或短对白；禁止改变事件结果、人物立场、原文因果，禁止新增关键剧情，禁止把模糊暗示改成确定事实。
3. 保留原小说气质。尽量保留原文专有名词、人物姓名、地点名称、修行体系名词、关键设定名词、核心意象和整体气氛。
4. 经典短句可以少量保留，但不能整段照抄。对白优先参考原文已有说法；没有原文台词时可以生成改编型对白，但不能加入新事实。
5. 对不确定信息使用“似乎”“暗示”“存在嫌疑”“可能”“背后有隐情”等模糊表达，不要编造确定细节。

示例：原文只说“古月赤城作弊”或“方源知道古月赤城作弊信息”时，允许写“古月赤城的测试结果存在作弊嫌疑”或“方源知道古月赤城背后有隐情”；禁止写“古月赤城偷偷服下作弊丹”“古月赤城袖中藏着作弊蛊”“古月赤练当场施法帮助作弊”。`

func BuildChapterAnalysisPrompt(chapter novel.Chapter) string {
	return fmt.Sprintf(`请根据下面的单章小说文本，输出严格符合 analysis.ChapterAnalysis 的 JSON。

忠实改编原则：
%s

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
- 只提取原文章节中明确出现或强烈暗示的信息
- characters、locations、key_events、conflicts、scene_candidates 都不得新增原文没有的事实
- scene_candidates 的 purpose 可以概括叙事功能，但不能新增剧情
- 对心理活动可以概括，但不能改写成原文没有发生的外部行为

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
%s`, faithfulAdaptationPrinciples, chapter.Number, chapter.Title, chapter.Text)
}

func BuildStoryBiblePrompt(analyses []analysis.ChapterAnalysis) string {
	payload, _ := json.MarshalIndent(analyses, "", "  ")
	return fmt.Sprintf(`请根据所有章节分析结果，合并成 story.StoryBible JSON。

忠实改编原则：
%s

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
- 合并多章分析时，不要发明新主线
- title / logline 可以概括，但不能改变故事类型和核心冲突
- global_characters 只能来自章节分析中出现过的人物
- timeline 只能基于已分析章节的事件
- main_conflict 可以抽象概括，但不能新增反派计划或隐藏设定
- scene_plan 只能来自 scene_candidates 和 key_events

要求：
- 只返回 JSON
- 不要 markdown
- 不要解释
- 不要输出 YAML
- 字段名必须符合上述 JSON 字段名
- 不要缺字段

章节分析 JSON：
%s`, faithfulAdaptationPrinciples, string(payload))
}

func BuildScreenplayPrompt(bible story.StoryBible) string {
	payload, _ := json.MarshalIndent(bible, "", "  ")
	return fmt.Sprintf(`请根据 StoryBible 生成 screenplay.Screenplay JSON。

这是忠实于原文事实的剧本化改编初稿，不是完全二创。

忠实改编原则：
%s

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
- source_chapters.title 不要重新命名；如果输入中已有章节标题，应原样保留
- 可以生成改编对白，但不得加入新事实
- dialogues.line 如果来自原文经典短句，可以短句保留；如果是改编对白，要符合人物处境和原文气质
- actions 不得新增原文没有的具体道具和行为
- 如果原文只是心理活动，actions 可以写为“沉默、观察、停顿、转身”等低风险表演动作，不要编造新动作
- 不要新增原文没有出现的人物

StoryBible JSON：
%s`, faithfulAdaptationPrinciples, string(payload))
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
