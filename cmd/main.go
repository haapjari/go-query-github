package main

import (
    "strings"
	"log"
	"github.com/haapjari/go-query-github/pkg/models"
	"github.com/haapjari/go-query-github/pkg/psql"
    "github.com/haapjari/go-query-github/pkg/ghb"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load(".env")	
    if err != nil {
        log.Fatal("Error loading .env file")
    }   

    p := psql.NewPostgreSQL("localhost", 5432, "postgres", "postgres", "postgres")
	db, err := p.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}

	var repos []models.Repo
	db.Find(&repos)

    totalRepos := len(repos) 
 
    for i, repo := range repos {
        parts := strings.Split(repo.Url, "/")
        owner := parts[1]
        url := repo.Url
        name := parts[2]

        log.Default().Printf("Processing Repository %d of %d\n", i+1, totalRepos)
        log.Default().Printf("Owner: %s, Name: %s\n", owner, name)

        g := ghb.NewGitHub() 

        releases, err := g.GetTotalReleasesCount(owner, name)
        if err != nil {
            log.Fatal(err)
        }

        err = p.UpdateRows(db, url, "releases", releases)
        if err != nil {
            log.Fatal(err)
        }

        deployments, err := g.GetTotalDeploymentsCount(owner, name)
        if err != nil {
            log.Fatal(err)
        }

        err = p.UpdateRows(db, url, "deployments", deployments)
        if err != nil {
            log.Fatal(err)
        }

        contributors, err := g.GetTotalContributorsCount(owner, name)
        if err != nil {
            log.Fatal(err)
        }

        err = p.UpdateRows(db, url, "contributors", contributors)
        if err != nil {
            log.Fatal(err)
        }

        notifications, err := g.GetTotalNotificationCount(owner, name)
        if err != nil {
            log.Fatal(err)
        }

        err = p.UpdateRows(db, url, "notifications", notifications)
        if err != nil {
            log.Fatal(err)
        }

        openPullRequests, closedPullRequests, err := g.FetchAllPullRequests(owner, name) 
        if err != nil {
            log.Fatal(err)
        }

        err = p.UpdateRows(db, url, "open_pulls", len(openPullRequests))
        if err != nil {
            log.Fatal(err)
        }

        err = p.UpdateRows(db, url, "closed_pulls", len(closedPullRequests))
        if err != nil {
            log.Fatal(err)
        }

        result, _, err := g.APIClient.Repositories.Get(g.APIClientContext, owner, name)    
        if err != nil {
            log.Fatal(err)
        }

        forks := *result.ForksCount

        err = p.UpdateRows(db, url, "forks", forks)
        if err != nil {
            log.Fatal(err)
        }

        subs := *result.SubscribersCount

        err = p.UpdateRows(db, url, "subscribers", subs)
        if err != nil {
            log.Fatal(err)
        }

        events := result.GetNetworkCount()

        err = p.UpdateRows(db, url, "network_events", events)
        if err != nil {
            log.Fatal(err)
        }

        watchers := result.GetWatchersCount()

        err = p.UpdateRows(db, url, "watchers", watchers)
        if err != nil {
            log.Fatal(err)
        }
        
        remainingRepos := totalRepos - (i + 1)

        log.Default().Printf("Completed Repository %d of %d. %d repositories remaining.\n\n", i+1, totalRepos, remainingRepos)
    }
}
