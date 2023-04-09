package main

import (
	"fmt"
	"log"
    //"golang.org/x/oauth2"
	//"github.com/google/go-github/github"
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

	// Get a single row from the database
	// Parse owner and repository from "repository_url" column
	// Query GitHup API: "get_avg_weekly_additions"

	// Read a value
	// var repo models.Repo
	var repos []models.Repo
	db.Find(&repos)

    // Access data in the repos struct

	// This is how to access the data in the repos
	// for _, repo := range repos {
	// 	parts := strings.Split(repo.Url, "/")
	// 	owner := parts[1]
	// 	repository := parts[2]
	// 	fmt.Printf("Owner: %s, Repository: %s\n", owner, repository)
	// }

	// TODO: Update Struct
    // TODO: Update Database
    // TODO: Notifications Count
    // TODO: Possible Optimizations
    // TODO: Contributors Count
    // TODO: Deployments Count
    // TODO: Average Weekly Additions
    // TODO: Average Weekly Deletions

    // AvgWeeklyAdditions int `json:"avg_weekly_additions"`
    // AvgWeeklyDeletions int `json:"avg_weekly_deletions"`

    // -------------------------------------------------------------------- //  
    // -------------------------------------------------------------------- //  
    // -------------------------------------------------------------------- //  

    // Fetch Pull Requests 

    owner := "kubernetes-sigs" 
    repo := "kind"

    g := ghb.NewGitHub() 

    // openPullRequests, closedPullRequests, err := g.FetchAllPullRequests(owner, repo) 
    // if err != nil {
    //     log.Fatal(err)
    // }

    // -------------------------------------------------------------------- //  
    // -------------------------------------------------------------------- //  
    // -------------------------------------------------------------------- //  

    result, _, err := g.APIClient.Repositories.Get(g.APIClientContext, owner, repo)    
    if err != nil {
        log.Fatal(err)
    }

    forks := *result.ForksCount
    // subs := *result.SubscribersCount
    // events := result.GetNetworkCount()
    // watchers := result.GetWatchersCount()

    fmt.Printf("Forks: %d\n", forks)

    // -------------------------------------------------------------------- //  
    // -------------------------------------------------------------------- //  
    // -------------------------------------------------------------------- //  

   
	// err = p.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
}
