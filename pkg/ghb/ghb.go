package ghb

import (
	"context"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"log"
	"net/http"
	"os"
)

type GitHub struct {
	HTTPClient *http.Client
    APIClient *github.Client
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

