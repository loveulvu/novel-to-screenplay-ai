type GeneratePanelProps = {
  loading: boolean;
  onGenerate: () => void;
};

export function GeneratePanel({ loading, onGenerate }: GeneratePanelProps) {
  return (
    <button className="primary-button" onClick={onGenerate} disabled={loading}>
      {loading ? "生成中..." : "生成剧本 YAML"}
    </button>
  );
}
