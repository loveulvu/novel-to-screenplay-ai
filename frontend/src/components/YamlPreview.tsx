"use client";

import { Button, Card, message } from "antd";

type YamlPreviewProps = {
  yaml: string;
  embedded?: boolean;
};

export function YamlPreview({ yaml, embedded = false }: YamlPreviewProps) {
  const [messageApi, contextHolder] = message.useMessage();

  async function handleCopy() {
    if (!yaml) return;
    await navigator.clipboard.writeText(yaml);
    messageApi.success("YAML 已复制");
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

  const content = (
    <>
      {contextHolder}
      <div className="preview-actions">
        <div>
          <span className="section-kicker">OUTPUT</span>
          <h2>YAML 剧本</h2>
        </div>
        <div className="button-row">
          <Button className="secondary-button" onClick={handleCopy} disabled={!yaml}>复制 YAML</Button>
          <Button className="secondary-button" onClick={handleDownload} disabled={!yaml}>下载 screenplay.yaml</Button>
        </div>
      </div>
      <pre className="yaml-output"><code>{yaml || "# 生成后将在这里展示结构化剧本 YAML"}</code></pre>
    </>
  );

  return embedded ? <div className="yaml-panel">{content}</div> : <Card className="tool-card yaml-panel">{content}</Card>;
}
