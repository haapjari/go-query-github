package main

import (
    "strings"
	"log"
	"github.com/haapjari/go-query-github/pkg/psql"
	"github.com/haapjari/go-query-github/pkg/models"
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

    // file, err := os.Open("data/cleaned_dataset.csv")
    // if err != nil {
    //     panic(err)
    // }
    // defer file.Close() 

    // // 3. Parse the .csv file into a slice of Repo structs
	// reader := csv.NewReader(file)
	// lines, err := reader.ReadAll()
	// if err != nil {
	// 	panic(err)
	// }

	// var repos []models.Repo
	// for _, line := range lines[1:] { // Skip the header line
	// 	openIssues, _ := strconv.Atoi(line[1])
	// 	closedIssues, _ := strconv.Atoi(line[2])
	// 	commits, _ := strconv.Atoi(line[3])
	// 	selfWrittenLOC, _ := strconv.Atoi(line[4])
	// 	libraryLOC, _ := strconv.Atoi(line[5])
	// 	stargazers, _ := strconv.Atoi(line[7])

	// 	repos = append(repos, models.Repo{
	// 		Url:            line[0],
	// 		OpenIssues:     openIssues,
	// 		ClosedIssues:   closedIssues,
	// 		Commits:        commits,
	// 		SelfWrittenLOC: selfWrittenLOC,
	// 		LibraryLOC:     libraryLOC,
	// 		CreationDate:   line[6],
	// 		Stargazers:     stargazers,
	// 		LatestRelease:  line[8],
	// 	})
	// }

    // for _, repo := range repos {
	// 	result := db.Create(&repo)
	// 	if result.Error != nil {
	// 		fmt.Printf("Error inserting repo (%s): %v\n", repo.Url, result.Error)
	// 	}
	// }

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
            log.Default().Printf("%s", err) 
        }

        err = p.UpdateRows(db, url, "releases", releases)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        deployments, err := g.GetTotalDeploymentsCount(owner, name)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        err = p.UpdateRows(db, url, "deployments", deployments)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        contributors, err := g.GetTotalContributorsCount(owner, name)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        err = p.UpdateRows(db, url, "contributors", contributors)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        notifications, err := g.GetTotalNotificationCount(owner, name)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        err = p.UpdateRows(db, url, "notifications", notifications)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        openPullRequests, closedPullRequests, err := g.FetchAllPullRequests(owner, name) 
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        err = p.UpdateRows(db, url, "open_pulls", len(openPullRequests))
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        err = p.UpdateRows(db, url, "closed_pulls", len(closedPullRequests))
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        result, _, err := g.APIClient.Repositories.Get(g.APIClientContext, owner, name)    
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        forks := *result.ForksCount

        err = p.UpdateRows(db, url, "forks", forks)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        subs := *result.SubscribersCount

        err = p.UpdateRows(db, url, "subscribers", subs)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        events := result.GetNetworkCount()

        err = p.UpdateRows(db, url, "network_events", events)
        if err != nil {
            log.Default().Printf("%s", err) 
        }

        watchers := result.GetWatchersCount()

        err = p.UpdateRows(db, url, "watchers", watchers)
        if err != nil {
            log.Default().Printf("%s", err) 
        }
        
        remainingRepos := totalRepos - (i + 1)

        log.Default().Printf("Completed Repository %d of %d. %d repositories remaining.\n\n", i+1, totalRepos, remainingRepos)
    }
}
