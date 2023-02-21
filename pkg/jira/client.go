package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const defaultStart = 0

type Client interface {
	// GetAllIssues returns all the issues satisfying the provided search params
	GetAllIssues(ctx context.Context, params SearchParams) ([]Issue, error)
}

type ClientImpl struct {
	httpCli *http.Client
	url     string
	token   string
}

type issuesResponse struct {
	TotalMatches int     `json:"total"`
	Limit        int     `json:"maxResults"`
	Issues       []Issue `json:"issues"`
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

func (c *ClientImpl) GetAllIssues(ctx context.Context, params SearchParams) ([]Issue, error) {
	params.StartAt = defaultStart
	res, err := c.getIssue(ctx, params)
	if err != nil {
		return nil, err
	}
	allIssues := res.Issues
	page := res.Limit
	limit := page
	for page <= res.TotalMatches {
		params.StartAt = page
		res, err := c.getIssue(ctx, params)
		if err != nil {
			return nil, err
		}
		page = page + limit
		allIssues = append(allIssues, res.Issues...)
	}
	return allIssues, nil
}

func (c *ClientImpl) getIssue(ctx context.Context, params SearchParams) (*issuesResponse, error) {
	data, err := c.request(ctx, "/rest/api/2/search", &params)
	if err != nil {
		return nil, err
	}
	var issuesRes issuesResponse
	err = json.Unmarshal(data, &issuesRes)
	if err != nil {
		return nil, err
	}
	return &issuesRes, nil
}

func (c *ClientImpl) request(ctx context.Context, uri string, params *SearchParams) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.url, uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if params != nil {
		q := req.URL.Query()
		q.Set("jql", params.Jql)
		q.Set("startAt", strconv.Itoa(params.StartAt))
		q.Set("fields", "summary, comment")
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
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
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
