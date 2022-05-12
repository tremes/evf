package bugzilla

// Bug is a type representing a single Bugzilla
type Bug struct {
	ID      int64
	Summary string
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
