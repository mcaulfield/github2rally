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

## Syncing Submitter and Assignee Fields

Correlating Github users and Rally users requires that each Rally user set their Middle Name field in their profile to match their github username.

## To-do
- [X] Set assignee field when Defect is created
- [ ] Update Defect when title or assignee changes
- [X] Link to github issue in Defect description
- [ ] Allow user to set Rally workspace
- [ ] Fix check for github issue number to disambiguate multiple results
