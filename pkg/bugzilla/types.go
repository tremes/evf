package bugzilla

// Bug is a type representing a single Bugzilla
type Bug struct {
	ID      int64
	Summary string
}

// BugsResponse is type representing the Bugzilla JSON response when
// requesting more bugs
type BugsResponse struct {
	TotalMatches int    `json:"total_matches"`
	Limit        string `json:"limit"`
	Bugs         []Bug  `json:"bugs"`
}

// BugResponse is a type representing the Bugzilla JSON response when
// requesting one bug
type BugResponse struct {
	BugComments map[int64]Comments `json:"bugs"`
}

type Comments struct {
	Comments []Comment `json:"comments"`
}

// SearchParams represents parameters for searching
// particular Bugzillas
type SearchParams struct {
	Product   string `yaml:"product"`
	Component string `yaml:"component"`
	Status    string `yaml:"status"`
	Version   string `yaml:"version"`
	Limit     int
	Offset    int
}

// Comment is a type representing a Bugzilla comment
type Comment struct {
	Creator string `json:"creator"`
	Text    string `json:"text"`
}
