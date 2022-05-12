package bugzilla

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Client interface {
	// GetAllBugs returns all the bugs satisfying the provided search params
	GetAllBugs(ctx context.Context, params SearchParams) ([]Bug, error)
	// GetComments returns all the comments for the related Bugzilla ID
	GetComments(ctx context.Context, bzID int64) (*Comments, error)
}

type ClientImpl struct {
	httpCli *http.Client
	url     string
	token   string
}

// bugsResponse is type representing the Bugzilla JSON response when
// requesting more bugs
type bugsResponse struct {
	TotalMatches int    `json:"total_matches"`
	Limit        string `json:"limit"`
	Bugs         []Bug  `json:"bugs"`
}

func NewClient(client *http.Client, url string, token string) Client {
	if client == nil {
		client = http.DefaultClient
	}
	return &ClientImpl{
		httpCli: client,
		url:     url,
		token:   token,
	}
}

func (c *ClientImpl) GetAllBugs(ctx context.Context, params SearchParams) ([]Bug, error) {
	res, err := c.getBug(ctx, params)
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
		params.Offset = page
		params.Limit = page
		res, err := c.getBug(ctx, params)
		if err != nil {
			return nil, err
		}
		page = page + limit
		allBugs = append(allBugs, res.Bugs...)

	}
	return allBugs, nil
}

// getBug makes the HTTP  GET request to the Bugzilla API to get the bugs
// satisfying the provided search parameters.
func (c *ClientImpl) getBug(ctx context.Context, params SearchParams) (*bugsResponse, error) {
	data, err := c.request(ctx, "/bug", &params, &c.token)
	if err != nil {
		return nil, err
	}
	var bugsRes bugsResponse
	err = json.Unmarshal(data, &bugsRes)
	if err != nil {
		return nil, err
	}
	return &bugsRes, nil
}

func (c *ClientImpl) GetComments(ctx context.Context, bzID int64) (*Comments, error) {
	uri := fmt.Sprintf("/bug/%d/comment", bzID)
	data, err := c.request(ctx, uri, nil, &c.token)
	if err != nil {
		return nil, err
	}
	bugsRes := struct {
		BugComments map[int64]Comments `json:"bugs"`
	}{}
	err = json.Unmarshal(data, &bugsRes)
	if err != nil {
		return nil, err
	}
	comments := bugsRes.BugComments[bzID]
	return &comments, nil
}

// request makes a HTTP GET request to the Bugzilla API appending the provided URI path.
// `searchParams` and `token` parameters are optional and can be nil.
func (c *ClientImpl) request(ctx context.Context, uri string, params *SearchParams, token *string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.url, uri)

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
	res, err := c.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	//r := io.LimitReader(res.Body, 8192)
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
