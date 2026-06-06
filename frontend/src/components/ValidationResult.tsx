import { Alert, Card } from "antd";
import type { FidelityResult, Validation } from "@/lib/api";

type ValidationResultProps = {
  validation: Validation | null;
  fidelityResult: FidelityResult | null;
  embedded?: boolean;
};

export function ValidationResult({ validation, fidelityResult, embedded = false }: ValidationResultProps) {
  if (!validation) {
    return (
      <Card className="tool-card">
        <div className="card-heading">
          <span className="section-kicker">QUALITY</span>
          <h2>质量检查</h2>
        </div>
        <p className="validation-waiting">生成后将检查 Schema 结构和事实一致性。</p>
      </Card>
    );
  }

  const content = (
    <>
      <div className="card-heading">
        <span className="section-kicker">QUALITY</span>
        <h2>质量检查</h2>
        <p>Schema Validate 检查结构完整性，Fidelity Check 检查内容是否忠实于原文事实。</p>
      </div>
      <div className="quality-grid">
        <Status label="Schema 校验" passed={validation.passed} failedText="失败" />
        <Status label="Fidelity Check" passed={Boolean(fidelityResult?.passed)} failedText="有风险" />
      </div>

      {!validation.passed ? (
        <div className="field-block">
          <h4>Schema 问题</h4>
          <Alert
            type="error"
            showIcon
            message="Schema 校验失败"
            description={<ul className="validation-errors">{validation.errors.map((error) => <li key={error}>{error}</li>)}</ul>}
          />
        </div>
      ) : null}

      {fidelityResult && fidelityResult.issues.length > 0 ? (
        <div className="field-block">
          <h4>事实一致性风险</h4>
          <div className="issue-list">
            {fidelityResult.issues.map((issue) => (
              <article className="issue-item" key={`${issue.field}-${issue.problem}`}>
                <div className="issue-heading">
                  <strong className={`severity severity-${issue.severity}`}>{issue.severity}</strong>
                  <code>{issue.field}</code>
                </div>
                <p>{issue.problem}</p>
                <small>建议：{issue.suggestion || "请根据原文核对。"}</small>
              </article>
            ))}
          </div>
        </div>
      ) : (
        <p className="quality-clear">未发现明显事实一致性问题。</p>
      )}
    </>
  );

  return embedded ? content : <Card className="tool-card">{content}</Card>;
}

function Status({ label, passed, failedText }: { label: string; passed: boolean; failedText: string }) {
  return (
    <div className="quality-status">
      <span>{label}</span>
      <strong className={passed ? "validation-ok" : "validation-bad"}>
        <i />{passed ? "通过" : failedText}
      </strong>
    </div>
  );
}
