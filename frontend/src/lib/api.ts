export type Validation = {
  passed: boolean;
  errors: string[];
};

export type GenerateResponse = {
  chapter_count: number;
  chapter_analyses: unknown[];
  story_bible: Record<string, unknown>;
  screenplay_json: Record<string, unknown>;
  screenplay_yaml: string;
  validation: Validation;
};

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export async function generateScreenplay(novelText: string): Promise<GenerateResponse> {
  const response = await fetch(`${API_BASE_URL}/api/generate`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      novel_text: novelText
    })
  });

  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error ?? "生成失败");
  }

  return data as GenerateResponse;
}
