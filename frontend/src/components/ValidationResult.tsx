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
          <h4>Schema issues</h4>
          <ul className="validation-errors">
            {validation.errors.map((error) => (
              <li key={error}>{error}</li>
            ))}
          </ul>
        </div>
      ) : null}

      {fidelityResult && fidelityResult.issues.length > 0 ? (
        <div className="field-block">
          <h4>Fidelity issues</h4>
          <ul className="validation-errors">
            {fidelityResult.issues.map((issue) => (
              <li key={`${issue.field}-${issue.problem}`}>
                <strong>{issue.severity}</strong> / {issue.field}：{issue.problem}
                {issue.suggestion ? ` 建议：${issue.suggestion}` : ""}
              </li>
            ))}
          </ul>
        </div>
      ) : null}
    </section>
  );
}
