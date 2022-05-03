package bugzilla

type Bug struct {
	ID      int64
	Summary string
}

type BugsResponse struct {
	TotalMatches int    `json:"total_matches"`
	Limit        string `json:"limit"`
	Bugs         []Bug  `json:"bugs"`
}

type BugResponse struct {
	BugComments map[int64]Comments `json:"bugs"`
}

type Comments struct {
	Comments []Comment `json:"comments"`
}

type BugParams struct {
	Product   string `yaml:"product"`
	Component string `yaml:"component"`
	Status    string `yaml:"status"`
	Version   string `yaml:"version"`
	Limit     int
	Offset    int
}
type Comment struct {
	Creator string `json:"creator"`
	Text    string `json:"text"`
}
