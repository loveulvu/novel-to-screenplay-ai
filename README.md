# Novel to Screenplay AI

七牛云 × XEngineer 暑期实训营第三批次议题三：**AI 小说转剧本工具**。

> 一个基于多阶段 AI Workflow 的小说转结构化 YAML 剧本工具。

用户输入 3～5 章小说文本后，系统先按章节解析并进行逐章结构化分析，再合并全局 Story Bible，最后生成经过事实一致性检查与结构校验的 YAML 剧本。

## 核心流程

```text
Novel Text
→ Parse Chapters
→ Analyze Each Chapter
→ Extract Factual Anchors
→ Merge Story Bible
→ Generate Screenplay JSON
→ Fidelity Check & Repair
→ Schema Validate
→ Export YAML
```

## 为什么不是简单 LLM 套壳

本项目没有直接一次性把全文交给模型生成 YAML，而是先进行章节级结构化分析，再合并 Story Bible，最后通过事实锚点和质量检查约束最终剧本生成。

这种固定 AI Workflow 更适合长文本改编：每个阶段都有明确输入输出，能够减少细节遗漏、人物漂移、事实幻觉和输出截断，也更方便展示、调试和替换模型。

## 核心功能

- 多章节小说输入
- 中文章 / 节 / 英文 `Chapter` 解析
- 分页章节合并
- 逐章 LLM 分析
- Story Bible 合并
- Factual Anchors 事实锚点
- Fidelity Check 与一次 Fidelity Repair
- YAML Schema Validate
- YAML 预览、复制和下载
- Mock / Real AI Provider 切换

## 关键设计

### Chapter Analysis

每章单独分析人物、地点、事件、冲突和候选场景，用于保留长文本细节。该阶段对应 Map-Reduce 中的 Map。

### Factual Anchors

记录原文中必须保留的硬事实，例如关键数字、人物关系、地点、事件结果和关键短句，用于约束后续剧本生成。

### Story Bible

将多章分析结果合并为全局故事资料，统一人物、时间线、主线冲突和分场计划。该阶段对应 Map-Reduce 中的 Reduce。

### 双重质量门控

- `Fidelity Check`：检查无依据事实、人物关系错误、关键数字错误和章节归属混乱；必要时执行一次定向 Repair。
- `Schema Validate`：检查最终剧本字段和结构是否完整，确保结果可被程序读取。

## 技术栈

- Frontend: Next.js + TypeScript
- Backend: Go + Gin
- LLM: OpenAI-compatible Chat Completions API
- Output: YAML
- Validation: Go struct + custom validator
- Intermediate State: JSON / Go struct

## 项目结构

```text
backend/                  Go API 与固定 AI Workflow
frontend/                 Next.js AI workbench
docs/architecture.md      架构与设计取舍
docs/schema.md            YAML Schema 设计说明
docs/demo_script.md       3～5 分钟录制脚本
examples/                 输入、中间态与最终 YAML 示例
```

## 本地运行

### 1. Mock 模式

Mock 模式不需要 API Key，适合本地开发和 Demo。

```bash
cd backend
go run ./cmd/server
```

另开终端启动前端：

```bash
cd frontend
npm install
npm run dev
```

默认地址：

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

### 2. Real 模式

根据 `.env.example` 创建本地 `.env`，不要提交 `.env` 或真实 API Key：

```env
AI_PROVIDER=real
AI_API_KEY=your_api_key_here
AI_BASE_URL=https://your-provider-compatible-api/v1
AI_MODEL=your-model-name
AI_TIMEOUT_SECONDS=180
```

环境变量说明：

- `AI_PROVIDER`：`mock` 或 `real`，未配置时默认使用 mock。
- `AI_API_KEY`：Real 模式需要的 API Key。
- `AI_BASE_URL`：OpenAI-compatible API 地址，可包含 `/v1`。
- `AI_MODEL`：使用的模型名称。
- `AI_TIMEOUT_SECONDS`：单次 LLM HTTP 请求超时时间，默认 `180` 秒。长文本分析时可适当增大，必须为正整数。

配置完成后启动后端和前端：

```bash
cd backend
go run ./cmd/server
```

```bash
cd frontend
npm run dev
```

如后端地址不同，可设置前端环境变量：

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080 npm run dev
```

## 测试

```bash
cd frontend
npm run build
```

```bash
cd backend
go test ./...
```

## 示例文件

- `examples/input_novel.txt`：三章小说输入
- `examples/chapter_analysis.json`：逐章分析结果
- `examples/story_bible.json`：全局 Story Bible
- `examples/output_screenplay.yaml`：最终 YAML 剧本

## 当前限制

- 生成质量仍受 LLM 能力影响
- Fidelity Check 不能保证 100% 无幻觉
- 推荐输入 3～5 章
- 当前不支持 PDF / Word
- 当前不支持登录、历史记录和多人协作
- 当前不包含数据库、文件上传、流式输出或桌面端

## 安全说明

- `.env`、`*.env` 和本地环境文件已加入 `.gitignore`
- 仓库只提交 `.env.example` 占位配置
- 不要在 README、docs、examples 或源码中写入真实 API Key
