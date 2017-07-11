// Github to Rally Sync Tool

package main

import (
	"flag"
	"fmt"
	"github2rally/sync"
	"os"
)

func main() {
	key := flag.String("rally-api-key", "", "CA Agile Central (Rally) API Key")
	owner := flag.String("github-owner", "", "Github user or organization")
	repo := flag.String("github-repo", "", "Github repository name")
	flag.Parse()
	if *key == "" || *owner == "" || *repo == "" {
		fmt.Fprintf(os.Stderr, "Missing arg, run with --help for usage\n")
		os.Exit(1)
	}
	sync.SyncDefects(*key, *owner, *repo)
}
