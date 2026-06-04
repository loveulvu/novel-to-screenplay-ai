type NovelInputProps = {
  value: string;
  onChange: (value: string) => void;
};

export function NovelInput({ value, onChange }: NovelInputProps) {
  return (
    <section className="panel">
      <h2>小说文本</h2>
      <textarea
        className="novel-textarea"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder="请粘贴至少三章小说文本，支持第1章、第一章、Chapter 1 等章节标题。"
      />
    </section>
  );
}
