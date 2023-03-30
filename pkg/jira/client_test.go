package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getIssue(t *testing.T) {
	issuesResponse := issuesResponse{
		TotalMatches: 1,
		Limit:        50,
		Issues: []Issue{
			{
				Key: "1111",
			},
		},
	}
	response, err := json.Marshal(issuesResponse)
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
	assert.Equal(t, issuesResponse.Issues, issues)
}
