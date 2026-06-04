# 架构说明

本项目第一版是无登录、无数据库、可本地运行的 MVP。系统采用固定 AI Workflow / Pipeline，而不是完全自主 agent。每一步都有明确输入输出，方便讲解、调试和替换真实 LLM。

## 流程

```text
Parse -> Analyze -> Merge -> Generate -> Validate -> Export
```

## Parse

位置：`backend/internal/novel`

解析用户粘贴的长文本，识别章节标题并切分章节。当前支持：

- `第1章`
- `第一章`
- `Chapter 1`
- `chapter 1`

少于 3 章时，接口直接返回 400。

## Analyze

位置：`backend/internal/analysis`

对每个章节生成结构化 `ChapterAnalysis`，包括摘要、角色、地点、关键事件、冲突和候选场景。当前由 `ai.MockClient` 返回示例结构，后续可替换真实 LLM。

## Merge

位置：`backend/internal/story`

把多个章节分析合并为全局故事资料 `StoryBible`。这一层解决跨章节一致性问题，例如统一角色、主冲突、时间线和分场计划。

## Generate

位置：`backend/internal/screenplay`

根据 `StoryBible` 生成结构化 `Screenplay` JSON，包括标题、来源章节、角色表和分场剧本。

## Validate

位置：`backend/internal/screenplay/validator.go`

对剧本结果做最小字段校验。当前检查标题、角色、分场，以及每个分场和台词的必要字段。

## Export

位置：`backend/internal/screenplay/yaml.go`

把内部 JSON 结构转换为最终 YAML 字符串，供前端展示、复制和下载。

## 长文本处理思想

项目参考了 Map-Reduce 的长文本处理思想：先对章节分别分析，再把章节结果合并为全局故事资料。第一版没有引入复杂框架，也没有做 agent 自主规划，而是保留清晰、固定、可解释的流水线。
