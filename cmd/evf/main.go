package main

import (
	"context"
	"evf/pkg/config"
	"evf/pkg/errata"
	"evf/pkg/jira"
	"fmt"
	"time"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Can't load config file:%v\n", err)
		return
	}
	err = c.InputPassword()
	if err != nil {
		fmt.Printf("Error inputing password:%v\n", err)
		return
	}

	fmt.Printf("Searching Errata versions for the Jira query: %s.\n", c.Jira.SearchParams.Jql)

	c.Jira.SearchParams.MaxResults = 50
	c.Jira.SearchParams.StartAt = 0
	jiraClient := jira.NewClient(nil, c.Jira.URL, c.Jira.Token)
	jiraHandler := jira.NewHandler(jiraClient)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	errataHandler, err := errata.New(c.Errata.URL,
		c.Errata.KerberosConf,
		c.Errata.Username,
		c.Errata.Realm,
		c.Errata.Password)
	if err != nil {
		fmt.Printf("Can't initiate Errata handler: %v\n", err)
		return
	}

	jbugs, err := jiraClient.GetAllBugs(ctx, c.Jira.SearchParams)
	fmt.Printf("Found %d related Jira bugs\n", len(jbugs))
	if err != nil {
		fmt.Printf("Can't read all the bugs from the Jira API: %v\n", err)
	}

	//create mapping errata ID -> slice of Jira bugs
	jiraToErrata := jiraHandler.CreateJiraToErrataMap(ctx, jbugs)
	ch := make(chan errata.Errata)
	// iterate over errata IDs and try to find version in X.Y.Z format
	go func() {
		for k := range jiraToErrata {
			syn, err := errataHandler.Synopsis(k)
			if err != nil {
				fmt.Printf("Can't read synopsis for the errata %s: %v\n", k, err)
				continue
			}
			e := errata.Errata{
				ID:       k,
				Synopsis: syn,
			}
			ch <- e
		}
		close(ch)
	}()

	// print results to stdout
	for e := range ch {
		for _, bug := range jiraToErrata[e.ID] {
			fmt.Printf("%s: %s %s\n", bug.Key, bug.Fields.Summary, e.Synopsis)
		}
	}
}
