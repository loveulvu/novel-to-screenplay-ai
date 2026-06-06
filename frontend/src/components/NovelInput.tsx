import { Button, Card, Input } from "antd";

const { TextArea } = Input;

type NovelInputProps = {
  value: string;
  onChange: (value: string) => void;
  onUseSample: () => void;
};

export function NovelInput({ value, onChange, onUseSample }: NovelInputProps) {
  return (
    <Card className="tool-card input-panel">
      <div className="card-heading">
        <span className="section-kicker">SOURCE</span>
        <h2>Source Novel</h2>
        <p>支持第 X 章 / 第 X 节 / Chapter X，推荐 3～5 章。</p>
      </div>
      <TextArea
        className="novel-textarea"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        placeholder={"第1章 标题\n粘贴章节正文...\n\n第二章 标题\n粘贴章节正文...\n\nChapter 3 Title\n粘贴章节正文..."}
        autoSize={{ minRows: 18, maxRows: 28 }}
      />
      <div className="input-footer">
        <span>{value.length.toLocaleString()} 字符</span>
        <Button className="secondary-button" type="default" onClick={onUseSample}>
          填入示例小说
        </Button>
      </div>
    </Card>
  );
}
