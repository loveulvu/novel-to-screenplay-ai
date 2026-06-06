"use client";

import { useState } from "react";
import { GeneratePanel } from "@/components/GeneratePanel";
import { NovelInput } from "@/components/NovelInput";
import { ResultSections } from "@/components/ResultSections";
import { ValidationResult } from "@/components/ValidationResult";
import { YamlPreview } from "@/components/YamlPreview";
import { generateScreenplay } from "@/lib/api";
import type { GenerateResponse } from "@/lib/api";

const sampleText = `第1章：开端
林舟在雨夜进入咖啡馆，遇见许晚。许晚递给他一封没有署名的信，信中提到一座废弃剧院。林舟认出信纸上的暗纹，和父亲失踪前留下的笔记完全一致。

第二章：追踪
林舟和许晚来到废弃剧院，发现舞台下方藏着一间旧档案室。档案里记录着林舟父亲多年前调查过的失踪案件，也提到一个反复出现的名字：顾衡。

第三章：对峙
顾衡出现在剧院，试图夺走档案。他警告林舟继续追查只会害了许晚。林舟必须决定是保护许晚离开，还是继续追查父亲失踪的真相。`;

const workflowSteps = [
  "Parse：识别章节并切分长文本",
  "Analyze：逐章提取角色、冲突和候选场景",
  "Merge：合并为全局 Story Bible",
  "Generate：按 Schema 生成剧本 JSON",
  "Validate：校验必填字段",
  "Export：导出 YAML"
];

export default function Home() {
  const [novelText, setNovelText] = useState("");
  const [result, setResult] = useState<GenerateResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleGenerate() {
    setLoading(true);
    setError("");
    setResult(null);

    try {
      const data = await generateScreenplay(novelText);
      setResult(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "生成失败，请稍后重试");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="page-shell">
      <section className="workspace">
        <div className="header">
          <div>
            <p className="eyebrow">七牛云 × XEngineer 暑期实训营</p>
            <h1>AI 小说转剧本工具</h1>
            <p className="subtitle">展示从长文本章节分析到 Story Bible，再到结构化 YAML 剧本的完整链路。</p>
          </div>
          <GeneratePanel loading={loading} onGenerate={handleGenerate} />
        </div>

        <section className="panel workflow-panel">
          <h2>长文本处理流水线</h2>
          <div className="workflow-steps">
            {workflowSteps.map((step) => (
              <span key={step}>{step}</span>
            ))}
          </div>
        </section>

        {error ? <div className="error-box">{error}</div> : null}

        <div className="layout-grid">
          <NovelInput value={novelText} onChange={setNovelText} onUseSample={() => setNovelText(sampleText)} />
          <div className="result-column">
            <ResultSections result={result} />
            <YamlPreview yaml={result?.screenplay_yaml ?? ""} />
            <ValidationResult validation={result?.validation ?? null} fidelityResult={result?.fidelity_result ?? null} />
          </div>
        </div>
      </section>
    </main>
  );
}
