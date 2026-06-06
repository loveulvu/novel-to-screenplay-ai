import type { ReactNode } from "react";
import { Card, Empty, Tag } from "antd";
import type { ChapterAnalysis, GenerateResponse, StoryBible } from "@/lib/api";

type ResultSectionsProps = {
  result: GenerateResponse | null;
  overviewOnly?: boolean;
  detailsOnly?: boolean;
};

export function ResultSections({ result, overviewOnly = false, detailsOnly = false }: ResultSectionsProps) {
  if (!result) {
    return (
      <Card className="tool-card empty-result">
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={
            <div className="empty-copy">
              <strong>Waiting for generated result</strong>
              <span>Submit novel text to generate structured screenplay YAML.</span>
            </div>
          }
        />
      </Card>
    );
  }

  if (overviewOnly) {
    return <BasicInfo result={result} />;
  }

  if (detailsOnly) {
    return (
      <section className="detail-area">
        <div className="detail-heading">
          <div>
            <span className="section-kicker">DETAILS</span>
            <h2>详细生成结果</h2>
          </div>
          <p>按章节、事实锚点与全局故事资料分区展示。</p>
        </div>
        <div className="detail-grid">
          <ChapterAnalyses analyses={result.chapter_analyses} />
          <div className="detail-side">
            <FactualAnchors analyses={result.chapter_analyses} />
            <StoryBibleView storyBible={result.story_bible} />
          </div>
        </div>
      </section>
    );
  }

  return null;
}

function BasicInfo({ result }: { result: GenerateResponse }) {
  const provider = result.meta?.ai_provider ?? "unknown";

  return (
    <Card className="tool-card">
      <div className="panel-title-row">
        <div>
          <span className="section-kicker">OVERVIEW</span>
          <h2>基础信息</h2>
        </div>
        <Tag className={provider === "real" ? "mode-badge real-mode" : "mode-badge mock-mode"}>
          {provider === "real" ? "Real LLM" : provider === "mock" ? "Mock Mode" : "Unknown"}
        </Tag>
      </div>
      <div className="metric-grid">
        <Metric label="AI provider" value={provider} />
        <Metric label="AI model" value={result.meta?.ai_model || "未配置"} />
        <Metric label="chapter_count" value={result.chapter_count} />
        <Metric label="Schema 校验" value={result.validation.passed ? "通过" : "失败"} status={result.validation.passed} />
        <Metric
          label="Fidelity Check"
          value={result.fidelity_result.passed ? "通过" : "有风险"}
          status={result.fidelity_result.passed}
        />
      </div>
    </Card>
  );
}

function Metric({ label, value, status }: { label: string; value: string | number; status?: boolean }) {
  return (
    <div className="metric">
      <span>{label}</span>
      <strong className={status === undefined ? "" : status ? "validation-ok" : "validation-bad"}>{value}</strong>
    </div>
  );
}

function ChapterAnalyses({ analyses }: { analyses: ChapterAnalysis[] }) {
  return (
    <Card className="tool-card">
      <div className="card-heading">
        <span className="section-kicker">CHAPTER ANALYSIS</span>
        <h2>章节分析</h2>
      </div>
      <div className="stack">
        {analyses.map((chapter) => (
          <details className="chapter-card" key={`${chapter.chapter_number}-${chapter.chapter_title}`} open={chapter.chapter_number === 1}>
            <summary>
              <span>CH.{String(chapter.chapter_number).padStart(2, "0")}</span>
              <strong>{chapter.chapter_title}</strong>
              <small>{chapter.scene_candidates.length} 个候选场景</small>
            </summary>
            <div className="chapter-content">
              <p className="chapter-summary">{chapter.summary}</p>
              <FieldList title="角色分析">
                {chapter.characters.map((character) => (
                  <li key={character.name}>
                    <strong>{character.name}</strong>：{character.role_in_chapter}
                    {character.traits.length ? `；特征：${character.traits.join("、")}` : ""}
                    {character.state_change ? `；变化：${character.state_change}` : ""}
                  </li>
                ))}
              </FieldList>
              <div className="chapter-columns">
                <FieldList title="地点" items={chapter.locations} />
                <FieldList title="关键事件" items={chapter.key_events} />
                <FieldList title="冲突" items={chapter.conflicts} />
              </div>
              <div className="field-block anchor-inline">
                <h4>事实锚点</h4>
                {chapter.factual_anchors.length ? (
                  <ul>{chapter.factual_anchors.map((anchor) => <li key={anchor}>{anchor}</li>)}</ul>
                ) : (
                  <p className="muted">本章暂无事实锚点。</p>
                )}
              </div>
              <div className="field-block">
                <h4>候选场景</h4>
                <div className="mini-grid">
                  {chapter.scene_candidates.map((scene, index) => (
                    <div className="mini-card" key={`${scene.location}-${index}`}>
                      <strong>{scene.location || "未指定地点"}</strong>
                      <span>{scene.time || "未指定时间"}</span>
                      <p>{scene.purpose}</p>
                      <small>角色：{scene.characters.join("、") || "无"}</small>
                      <small>事件：{scene.key_events.join("、") || "无"}</small>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </details>
        ))}
      </div>
    </Card>
  );
}

function FactualAnchors({ analyses }: { analyses: ChapterAnalysis[] }) {
  return (
    <Card className="tool-card anchors-panel">
      <div className="card-heading">
        <span className="section-kicker">FACTUAL ANCHORS</span>
        <h2>事实锚点</h2>
        <p>用于约束最终剧本并支撑 Fidelity Check。</p>
      </div>
      <div className="stack compact-stack">
        {analyses.map((chapter) => (
          <article className="anchor-group" key={`${chapter.chapter_number}-${chapter.chapter_title}-anchors`}>
            <h3>第 {chapter.chapter_number} 章 · {chapter.chapter_title}</h3>
            {chapter.factual_anchors.length ? (
              <ul>{chapter.factual_anchors.map((anchor) => <li key={anchor}>{anchor}</li>)}</ul>
            ) : (
              <p className="muted">本章暂无事实锚点。</p>
            )}
          </article>
        ))}
      </div>
    </Card>
  );
}

function StoryBibleView({ storyBible }: { storyBible: StoryBible }) {
  return (
    <Card className="tool-card">
      <div className="card-heading">
        <span className="section-kicker">STORY BIBLE</span>
        <h2>Story Bible</h2>
      </div>
      <div className="story-header">
        <h3>{storyBible.title}</h3>
        <p>{storyBible.logline}</p>
      </div>
      <FieldList title="全局角色">
        {storyBible.global_characters.map((character) => (
          <li key={character.id}><strong>{character.name}</strong>（{character.role}）：{character.motivation}</li>
        ))}
      </FieldList>
      <FieldList title="时间线">
        {storyBible.timeline.map((item) => <li key={`${item.chapter_number}-${item.event}`}>第 {item.chapter_number} 章：{item.event}</li>)}
      </FieldList>
      <div className="field-block">
        <h4>主线冲突</h4>
        <p>{storyBible.main_conflict}</p>
      </div>
      <div className="field-block">
        <h4>分场计划</h4>
        <div className="mini-grid">
          {storyBible.scene_plan.map((scene) => (
            <div className="mini-card" key={scene.id}>
              <strong>{scene.id}</strong>
              <span>第 {scene.source_chapter} 章 / {scene.location} / {scene.time}</span>
              <p>{scene.summary}</p>
              <small>角色：{scene.characters.join("、") || "无"}</small>
            </div>
          ))}
        </div>
      </div>
    </Card>
  );
}

function FieldList({ title, items, children }: { title: string; items?: string[]; children?: ReactNode }) {
  const values = items ?? [];
  if (!children && values.length === 0) return null;
  return (
    <div className="field-block">
      <h4>{title}</h4>
      <ul>{children ?? values.map((item) => <li key={item}>{item}</li>)}</ul>
    </div>
  );
}
