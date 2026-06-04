package screenplay

type Character struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Dialogue struct {
	Character string `json:"character"`
	Emotion   string `json:"emotion"`
	Line      string `json:"line"`
}

type Scene struct {
	ID            string     `json:"id"`
	SourceChapter int        `json:"source_chapter"`
	Location      string     `json:"location"`
	Time          string     `json:"time"`
	Summary       string     `json:"summary"`
	Characters    []string   `json:"characters"`
	Dialogues     []Dialogue `json:"dialogues"`
	Actions       []string   `json:"actions"`
}

type Screenplay struct {
	Title          string      `json:"title"`
	SourceChapters []int       `json:"source_chapters"`
	Characters     []Character `json:"characters"`
	Scenes         []Scene     `json:"scenes"`
}

type ValidationResult struct {
	Passed bool     `json:"passed"`
	Errors []string `json:"errors"`
}
