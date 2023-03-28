package jira

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	commentAuthor = "errata-owner+e-tool@redhat.com"
)

type Handler struct {
	client Client
}

func NewHandler(jiraClient Client) *Handler {
	return &Handler{
		client: jiraClient,
	}
}

func (h *Handler) CreateJiraToErrataMap(issues []Issue) map[string][]Issue {
	errataToJira := make(map[string][]Issue)
	for _, i := range issues {
		errataID := h.FindErrataID(&i)
		if errataID == "" {
			fmt.Printf("Didn't find the errata for the %s\n", i.Key)
			continue
		}
		errataToJira[errataID] = append(errataToJira[errataID], i)
	}
	return errataToJira
}

func (h *Handler) FindErrataID(jiraIssue *Issue) string {
	errataID := ""
	for _, c := range jiraIssue.Fields.Comments.Comment {
		if c.Author.EmailAddress == commentAuthor && strings.Contains(c.Body, "This issue has been added to advisory") {
			r, err := regexp.Compile(`advisory/\d+`)
			if err != nil {
				fmt.Printf("Can't compile the regex pattern: %v\n", err)
				continue
			}
			subStr := r.FindString(c.Body)
			errataID = strings.Split(subStr, "/")[1]
		}
	}
	return errataID
}
