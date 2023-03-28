package jira

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateJiraToErrataMap(t *testing.T) {
	issues := []Issue{
		{
			Key: "1",
			Fields: Field{
				Comments: Comments{
					Comment: []Comment{
						{
							Author: Author{
								EmailAddress: "errata-owner+e-tool@redhat.com",
							},
							Body: "This issue has been added to advisory [errata/advisory/1111]",
						},
					},
				},
			},
		},
		{
			Key: "2",
			Fields: Field{
				Comments: Comments{
					Comment: []Comment{
						{
							Author: Author{
								EmailAddress: "errata-owner+e-tool@redhat.com",
							},
							Body: "This issue has been added to advisory [errata/advisory/2222]",
						},
					},
				},
			},
		},
		{
			Key: "3",
			Fields: Field{
				Comments: Comments{
					Comment: []Comment{
						{
							Author: Author{
								EmailAddress: "errata-owner+e-tool@redhat.com",
							},
							Body: "This issue has been added to advisory [errata/advisory/1111]",
						},
					},
				},
			},
		},
	}
	tests := []struct {
		name     string
		issues   []Issue
		expected map[string][]Issue
	}{
		{
			name:   "one issue",
			issues: []Issue{issues[0]},
			expected: map[string][]Issue{
				"1111": {issues[0]},
			},
		},
		{
			name:     "no issue",
			issues:   []Issue{},
			expected: map[string][]Issue{},
		},
		{
			name:   "multiple issues with differnet ID",
			issues: []Issue{issues[0], issues[1]},
			expected: map[string][]Issue{
				"1111": {issues[0]},
				"2222": {issues[1]},
			},
		},
		{
			name:   "multiple issues with same ID",
			issues: []Issue{issues[0], issues[2]},
			expected: map[string][]Issue{
				"1111": {issues[0], issues[2]},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient(http.DefaultClient, "", "")
			handler := NewHandler(client)
			jiraToErrata := handler.CreateJiraToErrataMap(test.issues)
			assert.Equal(t, test.expected, jiraToErrata)
		})
	}
}
