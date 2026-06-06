"use client";

import { useEffect, useState } from "react";
import { Card, Empty, Segmented, Spin } from "antd";
import { OverviewContent } from "@/components/ResultSections";
import { ValidationResult } from "@/components/ValidationResult";
import { YamlPreview } from "@/components/YamlPreview";
import type { GenerateResponse } from "@/lib/api";

type OutputSection = "Overview" | "Quality" | "YAML";

type OutputPanelProps = {
  result: GenerateResponse | null;
  loading: boolean;
};

export function OutputPanel({ result, loading }: OutputPanelProps) {
  const [section, setSection] = useState<OutputSection>("Overview");

  useEffect(() => {
    if (result) setSection("Overview");
  }, [result]);

  return (
    <Card className="tool-card output-card">
      <div className="output-card-header">
        <div>
          <span className="section-kicker">OUTPUT</span>
          <h2>生成结果</h2>
        </div>
        <Segmented<OutputSection>
          options={["Overview", "Quality", "YAML"]}
          value={section}
          onChange={setSection}
          disabled={!result || loading}
        />
      </div>

      <div className="output-card-content">
        {loading ? (
          <div className="output-loading">
            <Spin />
            <strong>正在生成结构化剧本</strong>
            <span>章节分析、故事合并与质量检查正在进行。</span>
          </div>
        ) : !result ? (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <div className="empty-copy">
                <strong>Waiting for generated result</strong>
                <span>Submit novel text to generate structured screenplay YAML.</span>
              </div>
            }
          />
        ) : section === "Overview" ? (
          <OverviewContent result={result} />
        ) : section === "Quality" ? (
          <ValidationResult validation={result.validation} fidelityResult={result.fidelity_result} embedded />
        ) : (
          <YamlPreview yaml={result.screenplay_yaml} embedded />
        )}
      </div>
    </Card>
  );
}
