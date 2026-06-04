# novel-to-screenplay-ai

AI 小说转剧本工具 MVP，用于参加七牛云 × XEngineer 暑期实训营第三批次议题三。项目目标是把用户粘贴的至少 3 个章节小说文本，转换为结构化剧本 YAML。

## 题目对应关系

议题要求关注长文本小说到剧本的结构化转换。本项目把流程拆成可解释的固定 Pipeline：章节解析、章节级分析、全局故事资料合并、剧本生成、校验和 YAML 导出。

## MVP 功能

- 粘贴至少三章小说文本。
- 识别 `第1章`、`第一章`、`Chapter 1`、`chapter 1`。
- 返回章节级结构化分析 `ChapterAnalysis`。
- 合并全局故事资料 `StoryBible`。
- 生成结构化 `Screenplay` JSON。
- 校验必要字段。
- 展示、复制、下载最终 YAML。

本版本不包含登录、数据库、历史记录、多人协作、PDF/Word 解析和复杂编辑器。

## 技术栈

- 前端：Next.js + TypeScript
- 后端：Go + net/http
- 接口中间态：JSON
- 最终输出：YAML
- AI 客户端：`ai.MockClient` mock 实现，暂不接真实 LLM API

## 系统流程

```text
Parse -> Analyze -> Merge -> Generate -> Validate -> Export
```

- Parse：切分小说章节。
- Analyze：逐章提取结构化分析。
- Merge：合并成全局 Story Bible。
- Generate：生成剧本 JSON。
- Validate：检查剧本必要字段。
- Export：导出 YAML。

## 本地运行

### 启动后端

```bash
cd backend
go run ./cmd/server
```

默认监听 `http://localhost:8080`。

### 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认访问 `http://localhost:3000`。如后端地址不同，可设置：

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080 npm run dev
```

## 接口测试

### GET /api/health

```bash
curl http://localhost:8080/api/health
```

期望返回：

```json
{
  "status": "ok",
  "service": "novel-to-screenplay-ai"
}
```

### POST /api/generate

```bash
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d "{\"novel_text\":\"第1章 开端\n第一章内容\n第二章 发展\n第二章内容\nChapter 3 结局\n第三章内容\"}"
```

成功后返回章节数、章节分析、Story Bible、剧本 JSON、剧本 YAML 和校验结果。少于 3 章时返回：

```json
{
  "error": "at least 3 chapters are required"
}
```

## Demo 说明

示例输入位于 `examples/input_novel.txt`，对应中间态示例位于 `examples/chapter_analysis.json` 和 `examples/story_bible.json`，最终输出示例位于 `examples/output_screenplay.yaml`。

课堂或答辩演示时，可以先展示 `docs/architecture.md` 解释 Pipeline，再运行前后端，把示例小说粘贴到页面中生成 YAML。

## 接入真实 LLM 的位置

后续接入真实 OpenAI、七牛云或其他 LLM API 时，优先修改：

- `backend/internal/ai/client.go`
- `backend/internal/ai/mock_client.go`

建议新增真实 client，例如 `real_client.go`，实现与 `MockClient` 相同的方法：`AnalyzeChapter`、`MergeStoryBible`、`GenerateScreenplay`。这样可以保持 Parse -> Analyze -> Merge -> Generate -> Validate -> Export 的主流程不变。
