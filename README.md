# github2rally - Github to Rally Sync Tool

One-way sync from github issues to defects in CA Agile Central (formerly Rally).

```
Usage of github2rally:
  -github-owner string
        Github user or organization
  -github-repo string
        Github repository name
  -rally-api-key string
        CA Agile Central (Rally) API Key
```

## Getting Started

1. [Install go 1.8 or later](https://golang.org/doc/install) (and make sure $GOPATH/bin is on your $PATH)
2. `go get -u github.com/mcaulfield/github2rally`
3. [Request a CA Agile Central (Rally) API Key](https://rally1.rallydev.com/login/accounts/index.html#/keys)
4. Set your default workspace and project in Rally profile settings
5. Run `github2rally --github-owner <username or org> --github-repo <repo name> --rally-api-key <key from step 3>`

## Syncing Submitter and Assignee Fields

Correlating Github users and Rally users requires that each Rally user set their Middle Name field in their profile to match their github username.

## Todo

- [ ] Configurable Rally workspace/project
- [ ] Configurable Github credentials
- [ ] Support config file and env vars for config

