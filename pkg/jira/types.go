package jira

type Author struct {
	EmailAddress string `json:"emailAddress"`
}

type Comment struct {
	Author Author `json:"Author"`
	Body   string `json:"body"`
}

type Comments struct {
	Comment []Comment `json:"comments"`
}

type Field struct {
	Summary  string   `json:"summary"`
	Comments Comments `json:"comment"`
}

type Bug struct {
	ID     string `json:"id"`
	Fields Field  `json:"fields"`
}

type SearchParams struct {
	StartAt       int
	MaxResults    int    `yaml:"maxResults"`
	ValidateQuery bool   `yaml:"validationQuery"`
	Fields        string `yaml:"fields"`
	Expand        string `yaml:"expand"`
	Jql           string `yaml:"jql"`
}
