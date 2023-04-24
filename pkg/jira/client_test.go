package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllIssues(t *testing.T) {
	tests := []struct {
		name     string
		response issuesResponse
		exp      []Issue
	}{
		{
			name: "one issue",
			response: issuesResponse{
				TotalMatches: 1,
				Limit:        50,
				Issues: []Issue{
					{
						Key: "1111",
					},
				},
			},
			exp: []Issue{
				{
					Key: "1111",
				},
			},
		},
		{
			name: "multiple issues",
			response: issuesResponse{
				TotalMatches: 1,
				Limit:        50,
				Issues: []Issue{
					{
						Key: "1111",
					},
					{
						Key: "2222",
					},
					{
						Key: "3333",
					},
				},
			},
			exp: []Issue{
				{
					Key: "1111",
				},
				{
					Key: "2222",
				},
				{
					Key: "3333",
				},
			},
		},
		{
			name: "no issue",
			response: issuesResponse{
				TotalMatches: 1,
				Limit:        50,
				Issues:       []Issue{},
			},
			exp: []Issue{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := json.Marshal(test.response)
			assert.NoError(t, err)
			httpServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				if request.URL.Path == "/rest/api/2/search" {
					_, err := writer.Write(response)
					if err != nil {
						assert.NoError(t, err)
					}
				}
			}))
			defer httpServer.Close()
			client := NewClient(http.DefaultClient, httpServer.URL, "")
			ctx := context.Background()

			issues, err := client.GetAllIssues(ctx, SearchParams{})
			assert.NoError(t, err)
			assert.Equal(t, test.exp, issues)
		})
	}
}
