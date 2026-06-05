import type { Validation } from "@/lib/api";

type ValidationResultProps = {
  validation: Validation | null;
};

export function ValidationResult({ validation }: ValidationResultProps) {
  if (!validation) {
    return (
      <section className="panel">
        <h2>Schema 校验结果</h2>
        <p className="validation-waiting">等待生成。</p>
      </section>
    );
  }

  return (
    <section className="panel">
      <h2>Schema 校验结果</h2>
      {validation.passed ? (
        <p className="validation-ok">Schema 校验通过</p>
      ) : (
        <ul className="validation-errors">
          {validation.errors.map((error) => (
            <li key={error}>{error}</li>
          ))}
        </ul>
      )}
    </section>
  );
}
