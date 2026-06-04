package analysis

type ChapterAnalysis struct {
	ChapterNumber   int      `json:"chapter_number"`
	ChapterTitle    string   `json:"chapter_title"`
	Summary         string   `json:"summary"`
	Characters      []string `json:"characters"`
	Locations       []string `json:"locations"`
	KeyEvents       []string `json:"key_events"`
	Conflicts       []string `json:"conflicts"`
	SceneCandidates []string `json:"scene_candidates"`
}
