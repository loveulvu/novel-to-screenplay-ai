# Demo 讲解稿

## 1. 输入

用户在前端粘贴至少三章小说文本。系统支持中文章节标题和英文 `Chapter 1` 格式。

## 2. 章节切分

后端先执行 Parse，把长文本切分为章节。如果少于三章，直接返回错误，避免在过短文本上生成没有全局结构的剧本。

## 3. 章节级分析

每章进入 Analyze 阶段，产出 `ChapterAnalysis`。这一阶段像 Map：每章独立提取摘要、人物、地点、事件、冲突、候选场景和 `factual_anchors`。

## 4. 全局合并

Merge 阶段把章节级结果合并为 `StoryBible`，统一角色、时间线、主冲突和分场计划。

## 5. 剧本生成

Generate 阶段根据 `StoryBible`、章节分析和事实锚点产出结构化 `Screenplay` JSON。随后 Fidelity Check 检查事实一致性，必要时进行一次 Fidelity Repair，再经过 Validate 校验必要字段。

## 6. YAML 导出

Export 阶段把 JSON 转成 YAML。前端展示 YAML、Schema 校验结果和 Fidelity Check 质量报告，并提供复制和下载 `screenplay.yaml`。

## 7. 后续扩展

当前支持 mock 和 real 两种模式。继续扩展时仍应保持固定 Pipeline，不转成复杂 agent。
