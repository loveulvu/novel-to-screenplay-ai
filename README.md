# novel-to-screenplay-ai

面向小说作者的 AI 剧本化改编工具。用户粘贴 3 到 5 章小说文本后，系统会按章节做结构化分析，合并全局 Story Bible，再生成可校验、可复制、可下载的 YAML 剧本初稿。

本项目不是一次性把全文丢给模型生成 YAML 的简单 prompt 套壳，而是固定 AI Workflow / Pipeline：先做章节级 Map 分析，再做 Reduce 合并，最后用事实锚点、事实一致性检查和 Schema 校验约束剧本生成。

输出定位是“忠实于原文事实的剧本化改编初稿”：允许压缩叙述、生成保守的改编对白和可拍摄动作，但关键数字、人物关系、地点、事件结果与章节来源必须可回溯。

## 核心亮点

- 多章节长文本处理
- Map-Reduce 风格 AI Workflow
- Chapter-level structured analysis
- Story Bible 全局合并
- Factual Anchors 事实锚点
- Fidelity Check 事实一致性检查
- Fidelity Repair 一次事实修复
- Schema Validate 结构校验
- YAML Export 导出

## 架构流程

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

这条链路是固定 Pipeline，不是完全自主 agent。`Analyze Each Chapter` 是 Map 阶段，每章单独结构化分析；`Merge Story Bible` 是 Reduce 阶段，把多章人物、时间线、冲突和场景计划合并为全局资料。

## MVP 功能

- 粘贴至少三章小说文本
- 识别 `第1章`、`第一章`、`第一节`、`Chapter 1` 等章节标题
- 输出章节级 `ChapterAnalysis`
- 抽取每章 `factual_anchors`
- 合并 `StoryBible`
- 生成结构化 `Screenplay` JSON
- 执行 Fidelity Check 和必要的一次 Fidelity Repair
- 执行 Schema Validate
- 展示、复制、下载最终 YAML

本版本不包含登录、数据库、历史记录、多人协作、PDF/Word 解析、文件上传和复杂编辑器。

## 题目对应关系

七牛云 × XEngineer 暑期实训营第三批次议题三关注“AI 小说转剧本工具”。本项目对应的 MVP 是：输入多章节小说，经过长文本章节级分析、全局 Story Bible 合并、结构化剧本生成、事实一致性检查和 YAML 导出，形成可演示、可解释的小说转剧本工作流。

## 技术栈

- 前端：Next.js + TypeScript
- 后端：Go + Gin
- 中间态：JSON / Go struct
- 最终输出：YAML
- AI：mock client + OpenAI-compatible real client

## AI Provider 配置

后端默认使用 mock 模式；当 `AI_PROVIDER` 为空或为 `mock` 时，`/api/generate` 会使用内置 `MockClient`。

如需使用真实 OpenAI-compatible Chat Completions API，请在本地 `.env` 中配置：

```bash
AI_PROVIDER=real
AI_API_KEY=your_api_key_here
AI_BASE_URL=https://your-provider-compatible-api/v1
AI_MODEL=your-model-name
AI_TIMEOUT_SECONDS=180
```

不要提交 `.env`，只提交 `.env.example`。`AI_TIMEOUT_SECONDS` 用于控制每次真实 LLM HTTP 请求的超时时间，未配置时默认 180 秒。长文本在合并 Story Bible、生成剧本或事实检查阶段可能耗时更久，可以适当调大；配置必须是正整数秒数。

## 本地运行

Mock 模式无需配置密钥，启动后端：

```bash
cd backend
go run ./cmd/server
```

Real 模式先参考 `.env.example` 创建本地 `.env`，设置 `AI_PROVIDER=real` 及 API 配置，再使用同一命令启动。真实模式会执行逐章分析、Story Bible 合并、剧本生成和 Fidelity Check；任何模式都不会绕过 Schema Validate。

启动前端：

```bash
cd frontend
npm install
npm run dev
```

默认后端地址是 `http://localhost:8080`，默认前端地址是 `http://localhost:3000`。如后端地址不同，可设置：

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080 npm run dev
```

## 接口测试

```bash
curl http://localhost:8080/api/health
```

`POST /api/generate` 成功后返回章节数、章节分析、Story Bible、剧本 JSON、剧本 YAML、Schema 校验结果、Fidelity Check 结果和 meta 信息。

少于 3 章时返回：

```json
{
  "error": "at least 3 chapters are required"
}
```

## Demo 说明

演示时建议展示：

- 输入三章小说
- 章节分析
- Factual Anchors 事实锚点
- Story Bible
- YAML 剧本
- Schema 校验结果
- Fidelity Check 结果
- 复制 / 下载 YAML

示例输入位于 `examples/input_novel.txt`，中间态示例位于 `examples/chapter_analysis.json` 和 `examples/story_bible.json`，最终输出示例位于 `examples/output_screenplay.yaml`。

## 当前限制

- 真实生成质量仍取决于 LLM
- Fidelity Check 可降低事实偏差风险，但不能保证 100% 无幻觉
- 推荐输入 3 到 5 章
- 当前不支持 PDF/Word 解析
- 当前不支持历史记录和多人协作
