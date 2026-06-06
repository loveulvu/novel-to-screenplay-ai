import type { ReactNode } from "react";
import type { ChapterAnalysis, GenerateResponse, StoryBible } from "@/lib/api";

type ResultSectionsProps = {
  result: GenerateResponse | null;
};

export function ResultSections({ result }: ResultSectionsProps) {
  if (!result) {
    return (
      <section className="panel empty-result">
        <h2>结果区</h2>
        <p>生成成功后，这里会展示 AI 模式、章节分析、Story Bible、YAML 和 Schema 校验结果。</p>
      </section>
    );
  }

  return (
    <>
      <BasicInfo result={result} />
      <ChapterAnalyses analyses={result.chapter_analyses} />
      <StoryBibleView storyBible={result.story_bible} />
    </>
  );
}

function BasicInfo({ result }: { result: GenerateResponse }) {
  const provider = result.meta?.ai_provider ?? "unknown";
  const providerText =
    provider === "real" ? "当前使用真实 LLM 模式" : provider === "mock" ? "当前使用 Mock 模式" : "当前 AI 模式未知";

  return (
    <section className="panel">
      <h2>基础信息</h2>
      <p className={provider === "real" ? "mode-badge real-mode" : "mode-badge mock-mode"}>{providerText}</p>
      <div className="metric-grid">
        <div>
          <span>AI provider</span>
          <strong>{provider}</strong>
        </div>
        <div>
          <span>AI model</span>
          <strong>{result.meta?.ai_model || "未配置"}</strong>
        </div>
        <div>
          <span>章节数量</span>
          <strong>{result.chapter_count}</strong>
        </div>
        <div>
          <span>Schema 校验</span>
          <strong>{result.validation.passed ? "通过" : "未通过"}</strong>
        </div>
      </div>
    </section>
  );
}

function ChapterAnalyses({ analyses }: { analyses: ChapterAnalysis[] }) {
  return (
    <section className="panel">
      <h2>章节分析结果</h2>
      <div className="stack">
        {analyses.map((chapter) => (
          <article className="sub-card" key={`${chapter.chapter_number}-${chapter.chapter_title}`}>
            <h3>
              第 {chapter.chapter_number} 章：{chapter.chapter_title}
            </h3>
            <p>{chapter.summary}</p>

            <FieldList title="角色分析">
              {chapter.characters.map((character) => (
                <li key={character.name}>
                  <strong>{character.name}</strong>：{character.role_in_chapter}
                  {character.traits.length ? `；特征：${character.traits.join("、")}` : ""}
                  {character.state_change ? `；变化：${character.state_change}` : ""}
                </li>
              ))}
            </FieldList>

            <FieldList title="地点" items={chapter.locations} />
            <FieldList title="关键事件" items={chapter.key_events} />
            <FieldList title="冲突" items={chapter.conflicts} />
            <FieldList title="事实锚点" items={chapter.factual_anchors} />

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
          </article>
        ))}
      </div>
    </section>
  );
}

function StoryBibleView({ storyBible }: { storyBible: StoryBible }) {
  return (
    <section className="panel">
      <h2>Story Bible</h2>
      <div className="story-header">
        <h3>{storyBible.title}</h3>
        <p>{storyBible.logline}</p>
      </div>

      <FieldList title="全局角色">
        {storyBible.global_characters.map((character) => (
          <li key={character.id}>
            <strong>{character.name}</strong>（{character.role}）：{character.motivation}
            <span className="muted"> / {character.id}</span>
          </li>
        ))}
      </FieldList>

      <FieldList title="时间线">
        {storyBible.timeline.map((item) => (
          <li key={`${item.chapter_number}-${item.event}`}>
            第 {item.chapter_number} 章：{item.event}
          </li>
        ))}
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
              <span>
                第 {scene.source_chapter} 章 / {scene.location} / {scene.time}
              </span>
              <p>{scene.summary}</p>
              <small>角色：{scene.characters.join("、") || "无"}</small>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function FieldList({ title, items, children }: { title: string; items?: string[]; children?: ReactNode }) {
  const values = items ?? [];

  if (!children && values.length === 0) {
    return null;
  }

  return (
    <div className="field-block">
      <h4>{title}</h4>
      <ul>{children ?? values.map((item) => <li key={item}>{item}</li>)}</ul>
    </div>
  );
}
