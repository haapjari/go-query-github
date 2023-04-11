package models 

import (
	"github.com/google/go-github/github"
)

type Repo struct {
   Url string `json:"url"` 
   OpenIssues int `json:"open_issues"`
   ClosedIssues int `json:"closed_issues"`
   Commits int `json:"commits"`
   SelfWrittenLOC int `json:"self_written_loc"`
   LibraryLOC int `json:"library_loc"`
   CreationDate string `json:"creation_date"`
   Stargazers int `json:"stargazers"`
   LatestRelease string `json:"latest_release"`
   Forks int `json:"forks"`
   OpenPulls int `json:"open_pulls"`
   Releases int `json:"releases"`
   ClosedPulls int `json:"closed_pulls"`
   NetworkEvents int `json:"network_events"`
   Subscribers int `json:"subscribers"`
   Contributors int `json:"contributors"`
   Deployments int `json:"deployments"`
   Watchers int `json:"watchers"`
   Notifications int `json:"notifications"`
}

type PullRequest struct {
    URL             string       `json:"url"`
    ID              int          `json:"id"`
    NodeID          string       `json:"node_id"`
    HTMLURL         string       `json:"html_url"`
    DiffURL         string       `json:"diff_url"`
    PatchURL        string       `json:"patch_url"`
    IssueURL        string       `json:"issue_url"`
    Number          int          `json:"number"`
    State           string       `json:"state"`
    Locked          bool         `json:"locked"`
    Title           string       `json:"title"`
    User            User         `json:"user"`
    Body            string       `json:"body"`
    CreatedAt       string       `json:"created_at"`
    UpdatedAt       string       `json:"updated_at"`
    ClosedAt        string       `json:"closed_at"`
    MergedAt        string       `json:"merged_at"`
    MergeCommitSHA  string       `json:"merge_commit_sha"`
    Assignee        User         `json:"assignee"`
    Assignees       []User       `json:"assignees"`
    RequestedReviewers []User    `json:"requested_reviewers"`
    Head            PullRequestBranch `json:"head"`
}

type User struct {
    Login           string       `json:"login"`
    ID              int          `json:"id"`
    NodeID          string       `json:"node_id"`
    AvatarURL       string       `json:"avatar_url"`
    GravatarID      string       `json:"gravatar_id"`
    URL             string       `json:"url"`
    HTMLURL         string       `json:"html_url"`
    FollowersURL    string       `json:"followers_url"`
    FollowingURL    string       `json:"following_url"`
    GistsURL        string       `json:"gists_url"`
    StarredURL      string       `json:"starred_url"`
    SubscriptionsURL string     `json:"subscriptions_url"`
    OrganizationsURL string      `json:"organizations_url"`
    ReposURL        string       `json:"repos_url"`
    EventsURL       string       `json:"events_url"`
    ReceivedEventsURL string    `json:"received_events_url"`
    Type            string       `json:"type"`
    SiteAdmin       bool         `json:"site_admin"`
}

type PullRequestBranch struct {
    Label string `json:"label"`
    Ref   string `json:"ref"`
    SHA   string `json:"sha"`
    User  User   `json:"user"`
}

type ContributorStats struct {
    Total       int      `json:"total"`
    Author      github.User  `json:"author"`
    Weeks       []WeekStats `json:"weeks"`
}

type WeekStats struct {
    W int `json:"w"`
    A int `json:"a"`
    D int `json:"d"`
    C int `json:"c"`
}
