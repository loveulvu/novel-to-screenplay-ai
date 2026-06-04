# Demo 讲解稿

## 1. 输入

用户在前端粘贴至少三章小说文本。系统支持中文章节标题和英文 `Chapter 1` 格式。

## 2. 章节切分

后端先执行 Parse，把长文本切分为章节。如果少于三章，直接返回错误，避免在过短文本上生成没有全局结构的剧本。

## 3. 章节级分析

每章进入 Analyze 阶段，产出 `ChapterAnalysis`。这一阶段像 Map：每章独立提取摘要、人物、地点、事件、冲突和候选场景。

## 4. 全局合并

Merge 阶段把章节级结果合并为 `StoryBible`，统一角色、时间线、主冲突和分场计划。

## 5. 剧本生成

Generate 阶段根据 `StoryBible` 产出结构化 `Screenplay` JSON，再经过 Validate 校验必要字段。

## 6. YAML 导出

Export 阶段把 JSON 转成 YAML。前端展示 YAML，并提供复制和下载 `screenplay.yaml`。

## 7. 后续扩展

当前使用 `ai.MockClient`。接入真实模型时，只需要替换 `backend/internal/ai` 下的实现，保持 Pipeline 不变。
