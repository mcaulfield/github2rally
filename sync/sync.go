// Sync Github to Rally Tool

package sync

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github2rally/rally"
	"log"
	"os"
	"regexp"
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

	// Create matching defects in Rally or update existing
	for _, issue := range issues {
		if issue.PullRequestLinks != nil {
			continue
		}
		newDefect := constructDefect(rc, issue)
		existingDefect := findDefect(rc, newDefect)
		if existingDefect == nil {
			log.Println("Adding new defect to Rally: ", newDefect.Name)
			err = rc.CreateDefect(newDefect)
		} else {
			log.Println("Updating existing defect in Rally: ", newDefect.Name)
			err = rc.UpdateDefect(existingDefect, newDefect)
		}
		if err != nil {
			log.Println(err)
		}
	}
}

func constructDefect(rc *rally.Client, issue *github.Issue) *rally.Defect {
	// Query rally for users that match submitter and assignee
	submitter := matchUser(rc, issue.User)
	assignee := matchUser(rc, issue.Assignee)

	// Construct a new defect struct
	d := rally.Defect{}
	d.Name = fmt.Sprintf("%s #%v", *issue.Title, *issue.Number)
	if submitter != "" {
		d.SubmittedBy = &rally.Ref{URL: submitter, Type: "User"}
	}
	if assignee != "" {
		d.Owner = &rally.Ref{URL: assignee, Type: "User"}
	}
	d.Description = *issue.HTMLURL
	return &d
}

func matchUser(rc *rally.Client, user *github.User) string {
	if user == nil {
		return ""
	}
	// Map github usernames to rally users by relying on a matching Middle Name
	q := fmt.Sprintf("(MiddleName = %v)", *user.Login)
	qr, err := rc.QueryUser(q)
	if err != nil {
		log.Println(err)
		return ""
	}
	if qr.Count == 0 {
		return ""
	}
	return qr.Results[0].URL
}

func findDefect(rc *rally.Client, d *rally.Defect) *rally.Ref {
	re := regexp.MustCompile("#[0-9]+")
	newNum := re.FindString(d.Name)
	q := fmt.Sprintf("(Name contains %v)", newNum)
	qr, err := rc.QueryDefect(q)
	if err != nil {
		log.Println(err)
		return nil
	}
	if qr.Count == 0 {
		return nil
	}
	for _, ref := range qr.Results {
		oldNum := re.FindString(ref.Name)
		if oldNum == newNum {
			return &ref
		}
	}
	return nil
}
