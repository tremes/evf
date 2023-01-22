package jira

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
}

type ClientImpl struct {
	httpCli *http.Client
	url     string
	token   string
}

type bugsResponse struct {
	TotalMatches int   `json:"total"`
	Limit        int   `json:"maxResults"`
	Bugs         []Bug `json:"issues"`
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
	page := res.Limit
	params.MaxResults = page
	limit := page
	for page <= res.TotalMatches {
		params.StartAt = page
		res, err := c.getBug(ctx, params)
		if err != nil {
			return nil, err
		}
		page = page + limit
		allBugs = append(allBugs, res.Bugs...)
	}
	return allBugs, nil
}

func (c *ClientImpl) getBug(ctx context.Context, params SearchParams) (*bugsResponse, error) {
	data, err := c.request(ctx, "/rest/api/2/search", &params, &c.token)
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

func (c *ClientImpl) request(ctx context.Context, uri string, params *SearchParams, token *string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.url, uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	q := req.URL.Query()
	q.Set("jql", params.Jql)
	q.Set("startAt", strconv.Itoa(params.StartAt))
	q.Set("maxResults", "50")
	q.Set("validationQuery", "true")
	q.Set("fields", "summary, comment")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	req.Header.Add("Accept", "application/json")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return io.ReadAll(res.Body)
}
