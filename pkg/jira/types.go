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
	Key    string `json:"key"`
	Fields Field  `json:"fields"`
}

type SearchParams struct {
	StartAt    int
	MaxResults int    `yaml:"maxResults"`
	Jql        string `yaml:"jql"`
}
