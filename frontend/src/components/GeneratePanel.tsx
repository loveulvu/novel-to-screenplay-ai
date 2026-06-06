type GeneratePanelProps = {
  loading: boolean;
  onGenerate: () => void;
};

export function GeneratePanel({ loading, onGenerate }: GeneratePanelProps) {
  return (
    <div className="generate-panel">
      <button className="primary-button" onClick={onGenerate} disabled={loading}>
        {loading ? "正在生成..." : "生成剧本 YAML"}
        <span aria-hidden="true">→</span>
      </button>
      {loading ? (
        <p className="loading-text">正在进行章节分析、故事合并、剧本生成和事实一致性检查，请稍候。</p>
      ) : (
        <p className="helper-text">生成过程将调用后端配置的 AI provider。</p>
      )}
    </div>
  );
}
