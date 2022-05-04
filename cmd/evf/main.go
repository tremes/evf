package main

import (
	"context"
	"evf/pkg/bugzilla"
	"evf/pkg/config"
	"evf/pkg/errata"
	"fmt"
	"time"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Can't load config file:%v\n", err)
		return
	}
	fmt.Printf("Searching Errata versions for the \"%s\" product in version \"%s\" and \"%s\" component\n", c.Bugzilla.Product, c.Bugzilla.Version, c.Bugzilla.Component)
	bzHandler := bugzilla.New(c.Bugzilla.URL, c.Bugzilla.SearchParams, c.Bugzilla.Key)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
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
	//create mapping errata ID -> slice of BZ bugs
	errataToBZ := bzHandler.BugzillaToErrata(ctx)
	m := make(chan errata.Errata)

	// iterate over errata IDs and try to find version in X.Y.Z format
	go func() {
		for k, _ := range errataToBZ {
			syn, err := errataHandler.Synopsis(k)
			if err != nil {
				fmt.Printf("Can't read synopsis for the errata %s: %v\n", k, err)
				continue
			}
			e := errata.Errata{
				ID:       k,
				Synopsis: syn,
			}
			m <- e
		}
		close(m)
	}()

	// print results to stdout
	for e := range m {
		for _, bug := range errataToBZ[e.ID] {
			fmt.Printf("Bug %d: %s %s\n", bug.ID, bug.Summary, e.Synopsis)
		}
	}
}
