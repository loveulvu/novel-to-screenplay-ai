import type { FidelityResult, Validation } from "@/lib/api";

type ValidationResultProps = {
  validation: Validation | null;
  fidelityResult: FidelityResult | null;
};

export function ValidationResult({ validation, fidelityResult }: ValidationResultProps) {
  if (!validation) {
    return (
      <section className="panel">
        <h2>质量检查</h2>
        <p className="validation-waiting">等待生成。</p>
      </section>
    );
  }

  return (
    <section className="panel">
      <h2>质量检查</h2>
      <p className="quality-description">Schema Validate 检查结构完整性，Fidelity Check 检查内容是否忠实于原文事实。</p>
      <div className="quality-grid">
        <div>
          <span>Schema 校验</span>
          <strong className={validation.passed ? "validation-ok" : "validation-bad"}>
            {validation.passed ? "通过" : "失败"}
          </strong>
        </div>
        <div>
          <span>事实一致性</span>
          <strong className={fidelityResult?.passed ? "validation-ok" : "validation-bad"}>
            {fidelityResult?.passed ? "通过" : "有风险"}
          </strong>
        </div>
      </div>

      {!validation.passed ? (
        <div className="field-block">
          <h4>Schema 问题</h4>
          <ul className="validation-errors">
            {validation.errors.map((error) => (
              <li key={error}>{error}</li>
            ))}
          </ul>
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
                {issue.suggestion ? <small>建议：{issue.suggestion}</small> : null}
              </article>
            ))}
          </div>
        </div>
      ) : null}
    </section>
  );
}
