package ghb

import (
	"context"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"log"
	"net/http"
	"os"
    "fmt"
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

func (gh *GitHub) GetTotalContributorsCount(repoOwner, repoName string) (int, error) {
	// Get all repository contributors
	contributors, _, err := gh.APIClient.Repositories.ListContributors(gh.APIClientContext, repoOwner, repoName, &github.ListContributorsOptions{})
	if err != nil {
		return 0, err
	}

	// Count contributors
	contributorsCount := len(contributors)

	return contributorsCount, nil
}

func (gh *GitHub) GetTotalNotificationCount(owner, repo string) (int, error) {
    	// Get all user's notifications
	notifications, _, err := gh.APIClient.Activity.ListNotifications(gh.APIClientContext, &github.NotificationListOptions{All: true})
	if err != nil {
		return 0, err
	}

	// Filter notifications by the target repository
	r := fmt.Sprintf("%s/%s", owner, repo)
	count := 0
	for _, notification := range notifications {
		if notification.Repository != nil && notification.Repository.GetFullName() == r {
			count++
		}
	}

	return count, nil 
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

