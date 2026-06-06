"use client";

import { useState } from "react";
import { Alert, Card, ConfigProvider, Tag } from "antd";
import { GeneratePanel } from "@/components/GeneratePanel";
import { NovelInput } from "@/components/NovelInput";
import { OutputPanel } from "@/components/OutputPanel";
import { ResultSections } from "@/components/ResultSections";
import { generateScreenplay } from "@/lib/api";
import type { GenerateResponse } from "@/lib/api";

const sampleText = `第1章 雨夜钥匙
林澈在雨夜回到老街。巷口的邮筒早已停用，却在这天晚上吐出一封没有署名的信。
信封里只有一把铜钥匙和半张旧剧票。剧票背面写着父亲熟悉的字迹：海棠剧院，午夜之后。

第二章 旧剧院
许岚不放心林澈独自前往，带着手电陪他走进废弃的海棠剧院。
他们在舞台下方找到一台仍在运转的旧时钟，齿轮之间夹着另一半剧票。

Chapter 3 舞台对峙
顾衡在灯光亮起时出现，要求林澈交出铜钥匙。
林澈握紧钥匙，决定启动时钟，查清父亲消失的真相。`;

const workflowSteps = [
  ["01", "章节分析", "逐章提取人物、事件与场景"],
  ["02", "故事合并", "构建统一的 Story Bible"],
  ["03", "剧本生成", "生成结构化 YAML 剧本"],
  ["04", "质量检查", "校验 Schema 与事实一致性"]
];

export default function Home() {
  const [novelText, setNovelText] = useState("");
  const [result, setResult] = useState<GenerateResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleGenerate() {
    if (!novelText.trim()) {
      setError("请至少输入 2 个章节或分节。");
      return;
    }

    setLoading(true);
    setError("");

    try {
      setResult(await generateScreenplay(novelText));
    } catch (err) {
      const message = err instanceof Error ? err.message : "";
      setError(
        message.includes("at least 3 chapters are required")
          ? "请至少输入 3 个章节或分节。"
          : "生成失败，请检查后端服务、AI 配置或输入章节格式。"
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#1f1f1f",
          colorInfo: "#1f1f1f",
          colorBorder: "#e5e5e5",
          colorText: "#171717",
          colorTextSecondary: "#737373",
          borderRadius: 10,
          fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", "Microsoft YaHei", sans-serif'
        }
      }}
    >
      <main className="page-shell">
        <header className="site-header">
          <div>
            <h1>Novel to Screenplay AI</h1>
            <p>多章节小说 → Story Bible → YAML 剧本</p>
          </div>
          <div className="status-tags" aria-label="系统能力">
            <Tag><i className="status-dot" />Real LLM</Tag>
            <Tag>YAML Schema</Tag>
            <Tag>Fidelity Check</Tag>
          </div>
        </header>

        {error ? (
          <Alert
            className="error-card"
            type="error"
            showIcon
            message={error}
            description="确认服务与输入后，可在左侧重新发起生成。"
            closable
            onClose={() => setError("")}
          />
        ) : null}

        <div className="workspace-grid">
          <aside className="control-column">
            <NovelInput
              value={novelText}
              onChange={(value) => {
                setNovelText(value);
                if (error) setError("");
              }}
              onUseSample={() => {
                setNovelText(sampleText);
                setError("");
              }}
            />
            <GeneratePanel loading={loading} onGenerate={handleGenerate} />
            <Card className="tool-card workflow-card">
              <div className="card-heading">
                <span className="section-kicker">WORKFLOW</span>
                <h2>从小说到剧本</h2>
              </div>
              <ol className="workflow-list">
                {workflowSteps.map(([number, title, description]) => (
                  <li key={number}>
                    <span>{number}</span>
                    <div>
                      <strong>{title}</strong>
                      <p>{description}</p>
                    </div>
                  </li>
                ))}
              </ol>
            </Card>
          </aside>

          <section className="result-column">
            <OutputPanel result={result} loading={loading} />
          </section>
        </div>

        {result && !loading ? <ResultSections result={result} detailsOnly /> : null}
      </main>
    </ConfigProvider>
  );
}
