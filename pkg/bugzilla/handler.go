package bugzilla

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

const (
	CommentCreator = "errata-xmlrpc@redhat.com"
)

// Handler represents Bugzilla handler type creating
// the mapping between errata and bugzilla
type Handler struct {
	client Client
}

// New creates a new instance of the `Handler` type
func NewHandler(bugzillaClient Client) *Handler {
	return &Handler{
		client: bugzillaClient,
	}
}

// CreateBZToErrataMap creates mapping when key is the errata ID and each errata
// can have number of related Bugzilla bugs
func (h *Handler) CreateBZToErrataMap(ctx context.Context, bugs []Bug) map[string][]Bug {
	errataToBZ := make(map[string][]Bug)
	for _, b := range bugs {
		errataId := h.findErrataID(ctx, b.ID)
		if errataId == "" {
			fmt.Printf("Didn't find the errata for the Bug %d ID\n", b.ID)
			continue
		}
		if bugs, ok := errataToBZ[errataId]; ok {
			bugs = append(bugs, b)
			errataToBZ[errataId] = bugs
		} else {
			errataToBZ[errataId] = []Bug{b}
		}
	}
	return errataToBZ
}

// findErrataID iterates over all the comments of a respective bug and tries to
// find a comment containing errata information
func (h *Handler) findErrataID(ctx context.Context, bzBugID int64) string {
	comments, err := h.client.GetComments(ctx, bzBugID)
	if err != nil {
		fmt.Printf("Can't get data of the bug %d: %v\n", bzBugID, err)
	}
	for _, c := range comments.Comments {
		if c.Creator == CommentCreator && strings.Contains(c.Text, "Bug report changed to RELEASE_PENDING status") {
			r, err := regexp.Compile(`advisory/\d+`)
			if err != nil {
				fmt.Printf("Can't compile the regex pattern: %v\n", err)
			}
			subStr := r.FindString(c.Text)
			errataID := strings.Split(subStr, "/")[1]
			return errataID
		}
	}
	return ""
}
