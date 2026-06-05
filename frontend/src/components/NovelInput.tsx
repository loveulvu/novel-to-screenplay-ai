type NovelInputProps = {
  value: string;
  onChange: (value: string) => void;
  onUseSample: () => void;
};

export function NovelInput({ value, onChange, onUseSample }: NovelInputProps) {
  return (
    <section className="panel input-panel">
      <div className="section-heading">
        <div>
          <h2>小说文本</h2>
          <p>粘贴至少三章小说，后端会先做章节级分析，再合并 Story Bible，最后生成剧本 YAML。</p>
        </div>
        <button className="secondary-button" type="button" onClick={onUseSample}>
          填入示例小说
        </button>
      </div>
      <textarea
        className="novel-textarea"
        rows={18}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder="请粘贴至少三章小说文本，支持第1章、第一章、Chapter 1 等章节标题。"
      />
    </section>
  );
}
