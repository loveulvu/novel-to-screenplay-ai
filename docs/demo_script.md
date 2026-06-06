# 3～5 分钟 Demo 演示脚本

## 0:00～0:30 项目定位

画面：打开前端首页，停留在小说输入区和流程说明。

讲解：

> 这是一个面向小说作者的 AI 剧本化改编工具。它不是把全文一次性交给模型，而是通过固定的 Map-Reduce 风格工作流，先逐章分析，再合并 Story Bible，最后生成经过事实一致性检查和 Schema 校验的 YAML 剧本初稿。

## 0:30～1:00 输入三章小说

操作：

1. 点击“填入示例小说”，或粘贴 `examples/input_novel.txt`。
2. 简要展示三个章节标题和长文本输入框。
3. 点击“生成剧本 YAML”。

讲解：

> 输入至少三章小说后，后端首先解析章节。每个章节会独立进入分析阶段，因此长文本不会直接挤进一次生成请求。

## 1:00～2:00 展示章节分析与事实锚点

画面：生成完成后，展开或滚动查看章节分析结果。

重点展示：

- `chapter_number`、`chapter_title`、`summary`
- 章节人物、地点、关键事件和冲突
- `scene_candidates`
- `factual_anchors`

讲解：

> Analyze 是 Map 阶段。系统逐章提取人物、地点、冲突和候选场景，同时记录 factual anchors。事实锚点保存关键数字、人物关系、地点和事件结果，用来约束最终剧本，减少模型改错事实。

## 2:00～2:40 展示 Story Bible

画面：滚动到 Story Bible 区域，展示标题、logline、全局人物、时间线、主冲突和场景计划。

讲解：

> Merge 是 Reduce 阶段。它把多章分析合并为 Story Bible，统一跨章节人物、时间线、主冲突和场景计划，为最终剧本提供全局视角。

## 2:40～3:40 展示 YAML 剧本与质量报告

画面：展示 YAML 区域和 Quality Report。

重点展示：

- `source_chapters` 保留原始章节标题
- `characters` 使用稳定 id
- `scenes` 中的地点、对白和动作
- Schema 校验结果
- Fidelity Check 结果及 issues

讲解：

> Generate 会使用 Story Bible、章节分析和事实锚点生成 Screenplay JSON。随后 Fidelity Check 检查数字、人物关系、地点、道具和章节归属是否存在无依据补写；必要时最多修复一次。Schema Validate 则确保最终结构字段完整。两种检查关注点不同，不能互相替代。

## 3:40～4:10 复制与下载

操作：

1. 点击“复制 YAML”。
2. 点击“下载 screenplay.yaml”。

讲解：

> 中间态使用 JSON 和 Go struct，方便程序解析、API 传输和校验；最终导出 YAML，方便作者阅读、复制、修改和继续创作。

## 4:10～4:30 收尾

讲解：

> 当前 MVP 支持 mock 和 real 两种模式，不包含登录、数据库、历史记录或文件上传。真实生成质量仍取决于模型，但章节级分析、事实锚点、Fidelity Check 和 Schema Validate 让生成过程更可解释、更容易控制。
