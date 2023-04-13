package main

import (
	"log"
	"strings"

	"github.com/haapjari/go-query-github/pkg/ghb"
	"github.com/haapjari/go-query-github/pkg/models"
	"github.com/haapjari/go-query-github/pkg/psql"
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
		log.Default().Printf("%s", err)
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
			log.Default().Printf("%s", err)
		}

		contributors, err := g.GetTotalContributorsCount(owner, name)
		if err != nil {
			log.Default().Printf("%s", err)
		}

		err = p.UpdateRows(db, url, "releases", releases)
		if err != nil {
			log.Default().Printf("%s", err)
		}

		err = p.UpdateRows(db, url, "contributors", contributors)
		if err != nil {
			log.Default().Printf("%s", err)
		}

		remainingRepos := totalRepos - (i + 1)

		log.Default().Printf("Completed Repository %d of %d. %d repositories remaining.\n\n", i+1, totalRepos, remainingRepos)
	}
}
