package bugzilla

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Handler struct {
	url    string
	params *BugParams
	token  *string
}

func New(url string, params BugParams, token string) *Handler {
	return &Handler{
		url:    url,
		params: &params,
		token:  &token,
	}
}
func (h *Handler) getBugs(ctx context.Context, params *BugParams) (*BugsResponse, error) {
	data, err := h.request(ctx, "/bug", params, nil)
	if err != nil {
		return nil, err
	}
	var res BugsResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (h *Handler) getAllBugs(ctx context.Context) ([]Bug, error) {
	res, err := h.getBugs(ctx, h.params)
	if err != nil {
		return nil, err
	}
	allBugs := res.Bugs
	page, err := strconv.Atoi(res.Limit)
	if err != nil {
		return nil, err
	}
	limit := page
	for page <= res.TotalMatches {
		h.params.Offset = page
		h.params.Limit = page
		res, err := h.getBugs(ctx, h.params)
		if err != nil {
			return nil, err
		}
		page = page + limit
		allBugs = append(allBugs, res.Bugs...)

	}
	return allBugs, nil
}

func (h *Handler) BugzillaToErrata(ctx context.Context) map[string][]Bug {
	bugs, err := h.getAllBugs(ctx)
	fmt.Printf("Found %d related Bugzilla bugs\n", len(bugs))
	if err != nil {
		fmt.Printf("Can't read all the bugs from the Bugzilla API: %v\n", err)
	}
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

func (h *Handler) findErrataID(ctx context.Context, bzBugID int64) string {
	uri := fmt.Sprintf("/bug/%d/comment", bzBugID)
	data, err := h.request(ctx, uri, nil, h.token)
	if err != nil {
		fmt.Printf("Can't get data of the bug %d: %v\n", bzBugID, err)
	}
	var res BugResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		fmt.Printf("Can't unmarshal data of the bug %d: %v\n", bzBugID, err)
	}
	comments := res.BugComments[bzBugID]
	for _, c := range comments.Comments {
		if c.Creator == "errata-xmlrpc@redhat.com" && strings.Contains(c.Text, "Bug report changed to RELEASE_PENDING status") {
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

func (h *Handler) request(ctx context.Context, uri string, params *BugParams, token *string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", h.url, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if params != nil {
		q := req.URL.Query()
		limit := strconv.Itoa(params.Limit)
		offset := strconv.Itoa(params.Offset)
		q.Add("status", params.Status)
		q.Add("product", params.Product)
		q.Add("component", params.Component)
		q.Add("version", params.Version)
		q.Add("offset", offset)
		q.Add("limit", limit)
		req.URL.RawQuery = q.Encode()
	}
	if token != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	//r := io.LimitReader(res.Body, 8192)
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
