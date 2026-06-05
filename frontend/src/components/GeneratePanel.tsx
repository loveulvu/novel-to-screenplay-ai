type GeneratePanelProps = {
  loading: boolean;
  onGenerate: () => void;
};

export function GeneratePanel({ loading, onGenerate }: GeneratePanelProps) {
  return (
    <div className="generate-panel">
      <button className="primary-button" onClick={onGenerate} disabled={loading}>
        生成剧本 YAML
      </button>
      {loading ? <p className="loading-text">正在进行章节分析、故事合并和剧本生成，请稍候</p> : null}
    </div>
  );
}
