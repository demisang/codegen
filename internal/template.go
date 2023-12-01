package internal

type Placeholders []Placeholder

type Placeholder struct {
	Value       string `json:"value"`
	Replace     string `json:"replace"`
	Description string `json:"description"`
}

type Template struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	TargetDir    string        `json:"target_dir"`
	Description  string        `json:"description"`
	Placeholders []Placeholder `json:"placeholders"`
}

type PreviewListItem struct {
	IsDir   bool   `json:"is_dir"`
	IsNew   bool   `json:"is_new"`
	Path    string `json:"path"`
	Content string `json:"content"`
}
