package analysis

type CharacterMention struct {
	Name          string   `json:"name"`
	RoleInChapter string   `json:"role_in_chapter"`
	Traits        []string `json:"traits"`
	StateChange   string   `json:"state_change"`
}

type SceneCandidate struct {
	Location   string   `json:"location"`
	Time       string   `json:"time"`
	Purpose    string   `json:"purpose"`
	Characters []string `json:"characters"`
	KeyEvents  []string `json:"key_events"`
}

type ChapterAnalysis struct {
	ChapterNumber   int                `json:"chapter_number"`
	ChapterTitle    string             `json:"chapter_title"`
	Summary         string             `json:"summary"`
	Characters      []CharacterMention `json:"characters"`
	Locations       []string           `json:"locations"`
	KeyEvents       []string           `json:"key_events"`
	Conflicts       []string           `json:"conflicts"`
	SceneCandidates []SceneCandidate   `json:"scene_candidates"`
}
