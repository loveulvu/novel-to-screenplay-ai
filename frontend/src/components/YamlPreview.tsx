"use client";

import { useState } from "react";

type YamlPreviewProps = {
  yaml: string;
};

export function YamlPreview({ yaml }: YamlPreviewProps) {
  const [copied, setCopied] = useState(false);

  async function handleCopy() {
    if (!yaml) return;
    await navigator.clipboard.writeText(yaml);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1600);
  }

  function handleDownload() {
    if (!yaml) return;
    const blob = new Blob([yaml], { type: "text/yaml;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "screenplay.yaml";
    link.click();
    URL.revokeObjectURL(url);
  }

  return (
    <section className="panel yaml-panel">
      <div className="preview-actions">
        <div>
          <span className="section-kicker">OUTPUT</span>
          <h2>YAML 剧本</h2>
        </div>
        <div className="button-row">
          <button className="secondary-button" onClick={handleCopy} disabled={!yaml}>{copied ? "已复制" : "复制 YAML"}</button>
          <button className="secondary-button" onClick={handleDownload} disabled={!yaml}>下载 screenplay.yaml</button>
        </div>
      </div>
      <pre className="yaml-output"><code>{yaml || "# 生成后将在这里展示结构化剧本 YAML"}</code></pre>
    </section>
  );
}
