package jira

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

const (
	CommentAuthor = "errata-owner+e-tool@redhat.com"
)

type Handler struct {
	client Client
}

func NewHandler(jiraClient Client) *Handler {
	return &Handler{
		client: jiraClient,
	}
}

func (h *Handler) CreateJiraToErrataMap(ctx context.Context, bugs []Bug) map[string][]Bug {
	jiraToErrata := make(map[string][]Bug)
	for _, b := range bugs {
		jiraID := h.FindErrataID(ctx, &b)
		if jiraID == "" {
			fmt.Printf("Didn't find the errata for the %s\n", b.Key)
			continue
		}
		if bugs, ok := jiraToErrata[jiraID]; ok {
			jiraToErrata[jiraID] = append(bugs, b)
		} else {
			jiraToErrata[jiraID] = []Bug{b}
		}
	}
	return jiraToErrata
}

func (h *Handler) FindErrataID(ctx context.Context, jiraBug *Bug) string {
	for _, c := range jiraBug.Fields.Comments.Comment {
		if c.Author.EmailAddress == CommentAuthor && strings.Contains(c.Body, "This issue has been added to advisory") {
			r, err := regexp.Compile(`advisory/\d+`)
			if err != nil {
				fmt.Printf("Can't compile the regex pattern: %v\n", err)
				return ""
			}
			subStr := r.FindString(c.Body)
			errataID := strings.Split(subStr, "/")[1]
			return errataID
		}
	}
	return ""
}
