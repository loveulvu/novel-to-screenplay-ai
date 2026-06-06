# 架构说明

本项目第一版是无登录、无数据库、可本地运行的 MVP。系统采用固定 AI Workflow / Pipeline，而不是完全自主 agent。每一步都有明确输入输出，方便讲解、调试和后续替换真实 LLM。

## 主流程

```text
Novel Text
-> Parse Chapters
-> Analyze Each Chapter
-> Extract Factual Anchors
-> Merge Story Bible
-> Generate Screenplay JSON
-> Fidelity Check & Repair
-> Schema Validate
-> Export YAML
```

代码主链路保持稳定：

```text
ParseChapters
-> AnalyzeChapters
-> Extract Factual Anchors
-> MergeStoryBible
-> GenerateScreenplay
-> FidelityCheck
-> FidelityRepair
-> ValidateScreenplay
-> ToYAML
```

当前支持 `ai.MockClient` 和 OpenAI-compatible `RealClient`，通过 `AI_PROVIDER` 切换。

## Parse

位置：`backend/internal/novel`

解析用户粘贴的长文本，识别章节标题并切分章节。章节标题前的前言、简介和空行不会被自动生成未命名章节；如果全文没有章节标题，会返回空切片，再由 handler 按少于三章返回 400。

## Analyze

位置：`backend/internal/analysis`

Analyze 类似 Map 阶段：每章单独分析，产出 `ChapterAnalysis`。这一层不会直接生成剧本，而是把每章转成更适合改编的结构化信息：

- 章节摘要
- 单章角色提及 `CharacterMention`
- 地点
- 关键事件
- 冲突
- 候选场景 `SceneCandidate`
- 硬事实锚点 `factual_anchors`

这样做可以避免直接从长篇原文跳到最终剧本，也方便解释长文本处理过程。

## Merge

位置：`backend/internal/story`

Merge 类似 Reduce 阶段：把多章分析结果合并为全局故事资料 `StoryBible`。这一层负责统一跨章节信息，例如角色、主冲突、时间线和分场计划。

## Generate

位置：`backend/internal/screenplay`

Generate 根据 `StoryBible`、`ChapterAnalysis` 和 `factual_anchors` 生成结构化 `Screenplay`。最终剧本包含：

- 标题
- 来源章节 `source_chapters`
- 全局角色表 `characters`
- 分场剧本 `scenes`

## Fidelity Check & Repair

位置：`backend/internal/fidelity` 和 `backend/internal/ai`

Fidelity Check 用于检查 unsupported claims，例如关键数字、资质等级、人物关系、地点、道具、观察方式、身体反应、章节归属、actions 或 dialogues 是否被改写或凭空补写。

当检查结果包含 `medium` 或 `high` 风险时，系统会触发一次 Fidelity Repair。Repair 只修复 issues 指出的事实问题，修复后再检查一次；如果仍有风险，不无限重试，而是把 `fidelity_result` 返回给前端展示。

## Validate

位置：`backend/internal/screenplay/validator.go`

Validate 用于防止 LLM 输出缺字段或结构不完整。即使当前是 mock 数据，也保留校验层，方便后续替换真实模型后约束输出质量。

当前重点校验：

- `title`
- `source_chapters.number/title`
- `characters.id/name/role`
- `scenes.id/location/time/summary`
- `scenes.characters`
- `scenes.dialogues`
- `dialogues.character/line`

`dialogue.emotion` 不强制，因为有些台词可以是中性表达。`actions` 暂时不强制，因为有些场景可以主要由对白推进。

## Export

位置：`backend/internal/screenplay/yaml.go`

Export 把内部 JSON / Go struct 转换为最终 YAML 字符串，供前端展示、复制和下载。中间态使用 JSON / Go struct 是为了便于程序解析、校验和 API 传输；最终输出 YAML 是为了更适合人工阅读、编辑，也符合题目要求。

## 为什么不是 Agent

本项目不是完全自主 agent，不让模型自由决定工具调用、循环步骤或长期记忆。它是固定 AI Workflow / Pipeline：

- Analyze：逐章分析，类似 Map。
- Merge：合并多章结果，类似 Reduce。
- Generate：按 Story Bible、章节分析和事实锚点生成剧本。
- Fidelity Check & Repair：降低事实偏差风险。
- Validate：检查结构完整性。
- Export：输出 YAML。

这种设计更适合 MVP 演示：结构清楚、边界明确、容易解释，也方便后续替换或扩展真实 LLM client。
