// Sync Github to Rally Tool

package sync

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github2rally/rally"
	"log"
	"os"
)

// SyncDefects copies github issues to rally defects.
func SyncDefects(rallyApiKey string, owner string, repo string) {
	// Set up github and rally api clients
	gc := github.NewClient(nil)
	rc := rally.NewClient(rallyApiKey)

	// Query github for open issues
	ctx := context.Background()
	issues, _, err := gc.Issues.ListByRepo(ctx, owner, repo, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	for _, issue := range issues {

		// Query rally for a matching defect
		q := fmt.Sprintf("(Name contains #%v)", *issue.Number)
		qr, err := rc.QueryDefect(q)
		if err != nil {
			log.Println(err)
			continue
		}
		if qr.Count != 0 {
			log.Println("Defect already exists in Rally:", *issue.Title, qr.Results[0].Ref)
			continue
		}
		log.Println("Adding new Defect to Rally:", *issue.Title)

		// Query rally for a user that matches the github issue user
		submitter := ""
		if issue.User != nil {
			q = fmt.Sprintf("(MiddleName = %v)", *issue.User.Login)
			qr, err = rc.QueryUser(q)
			if err != nil {
				log.Println(err)
				continue
			}
			if qr.Count != 0 {
				submitter = qr.Results[0].Ref
			}
		}

		// Query rally for a user that matches the github assignee
		assignee := ""
		if issue.Assignee != nil {
			q = fmt.Sprintf("(MiddleName = %v)", *issue.Assignee.Login)
			qr, err = rc.QueryUser(q)
			if err != nil {
				log.Println(err)
			}
			if qr.Count != 0 {
				assignee = qr.Results[0].Ref
			}
		}

		// Construct a new defect struct
		d := rally.Defect{}
		d.Name = fmt.Sprintf("%s #%v", *issue.Title, *issue.Number)
		if submitter != "" {
			d.SubmittedBy = &rally.Ref{Ref: submitter, Type: "User"}
		}
		if assignee != "" {
			d.Owner = &rally.Ref{Ref: assignee, Type: "User"}
		}
		d.ScheduleState = "Defined"
		d.State = "Submitted"

		// Create new rally defect
		err = rc.CreateDefect(&d)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Added new Defect: ", d)
	}
}
