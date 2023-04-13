package ghb

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type GitHub struct {
	HTTPClient       *http.Client
	APIClient        *github.Client
	APIClientContext context.Context
}

func NewGitHub() *GitHub {
	gh := &GitHub{}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)

	// create http client with OAuth2 token source
	gh.HTTPClient = oauth2.NewClient(context.Background(), tokenSource)

	// authentication for the github client
	gh.APIClientContext = context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)

	tc := oauth2.NewClient(gh.APIClientContext, ts)

	gh.APIClient = github.NewClient(tc)

	return gh
}

func (gh *GitHub) GetTotalReleasesCount(repoOwner, repoName string) (int, error) {
	releasesCount := 0
	opt := &github.ListOptions{PerPage: 100}

	log.Default().Printf("Fetching releases for %s/%s", repoOwner, repoName)

	for {
		releases, resp, err := gh.APIClient.Repositories.ListReleases(gh.APIClientContext, repoOwner, repoName, opt)
		if err != nil {
			return 0, err
		}

		releasesCount += len(releases)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage

		log.Default().Printf("Page: %d\n", opt.Page)
	}

	return releasesCount, nil
}

func (gh *GitHub) GetTotalContributorsCount(repoOwner, repoName string) (int, error) {
	contributors := make(map[int64]struct{})
	page := 1
	perPage := 100

	log.Default().Printf("Fetching contributors for %s/%s", repoOwner, repoName)

	for {
		opts := &github.CommitsListOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}

		commits, resp, err := gh.APIClient.Repositories.ListCommits(gh.APIClientContext, repoOwner, repoName, opts)
		if err != nil {
			return 0, err
		}

		for _, commit := range commits {
			if commit.Author != nil && commit.Author.ID != nil {
				contributors[*commit.Author.ID] = struct{}{}
			}
		}

		if resp.NextPage == 0 {
			break
		}

		page = resp.NextPage

		log.Default().Printf("Page: %d\n", page)

        // In order not to get blocked from GitHub API.
		time.Sleep(2 * time.Second)
	}

	return len(contributors), nil
}

func (gh *GitHub) FetchAllPullRequests(owner, repo string) ([]*github.PullRequest, []*github.PullRequest, error) {
	log.Default().Printf("Fetching pull requests for repository: %s\n", repo)

	opt := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allPullRequests, openPullRequests, closedPullRequests []*github.PullRequest

	page := 1
	for {
		pullRequests, resp, err := gh.APIClient.PullRequests.List(gh.APIClientContext, owner, repo, opt)
		if err != nil {
			return nil, nil, err
		}

		allPullRequests = append(allPullRequests, pullRequests...)

		for _, pr := range pullRequests {
			if pr.GetState() == "open" {
				openPullRequests = append(openPullRequests, pr)
			} else if pr.GetState() == "closed" {
				closedPullRequests = append(closedPullRequests, pr)
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage

		log.Default().Printf("Page: %d\n", page)
		page++
	}

	return openPullRequests, closedPullRequests, nil
}

func (gh *GitHub) GetAverageWeeklyAdditionsAndDeletions(repoOwner, repoName string) (float64, float64, error) {
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	commits, _, err := gh.APIClient.Repositories.ListCommits(gh.APIClientContext, repoOwner, repoName, opt)
	if err != nil {
		return 0, 0, err
	}

	// Calculate total additions and deletions
	totalAdditions := 0
	totalDeletions := 0
	commitsCounted := 0
	for _, commit := range commits {
		if commit.Stats != nil {
			totalAdditions += *commit.Stats.Additions
			totalDeletions += *commit.Stats.Deletions
			commitsCounted++
		}
	}

	// Calculate average weekly additions and deletions (assuming last 100 commits are within 1 week)
	averageWeeklyAdditions := float64(totalAdditions) / float64(commitsCounted)
	averageWeeklyDeletions := float64(totalDeletions) / float64(commitsCounted)

	return averageWeeklyAdditions, averageWeeklyDeletions, nil
}
