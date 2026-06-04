"use client";

import { useState } from "react";
import { GeneratePanel } from "@/components/GeneratePanel";
import { NovelInput } from "@/components/NovelInput";
import { ValidationResult } from "@/components/ValidationResult";
import { YamlPreview } from "@/components/YamlPreview";
import { generateScreenplay } from "@/lib/api";
import type { GenerateResponse } from "@/lib/api";

const sampleText = `第1章 雨夜钥匙
林澈在雨夜回到老街，收到一封没有署名的信。信封里只有一把铜钥匙和父亲失踪前留下的半张剧票。

第二章 旧剧院
许岚陪林澈来到废弃多年的海棠剧院。舞台下方传来钟摆声，他们在地下室发现一台仍在运转的旧时钟。

Chapter 3 舞台对峙
顾衡突然出现，要求林澈交出钥匙。他说那台时钟保存着不该被唤醒的记忆，而林澈终于意识到父亲可能还活着。`;

export default function Home() {
  const [novelText, setNovelText] = useState(sampleText);
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
      setError(err instanceof Error ? err.message : "生成失败");
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
          </div>
          <GeneratePanel loading={loading} onGenerate={handleGenerate} />
        </div>

        <div className="grid">
          <NovelInput value={novelText} onChange={setNovelText} />
          <div className="result-column">
            {error ? <div className="error-box">{error}</div> : null}
            <ValidationResult validation={result?.validation ?? null} />
            <YamlPreview yaml={result?.screenplay_yaml ?? ""} />
          </div>
        </div>
      </section>
    </main>
  );
}
