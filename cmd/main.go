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

        if repo.Releases == 0 {
            releases, err := g.GetTotalReleasesCount(owner, name)
            if err != nil {
                log.Default().Printf("%s", err)
            }

            if releases == 0 {
                err = p.UpdateRows(db, url, "releases", releases)
                if err != nil {
                    log.Default().Printf("%s", err)
                }
            }
        }

        if repo.Deployments == 0 { 
            deployments, err := g.GetTotalDeploymentsCount(owner, name)
            if err != nil {
                log.Default().Printf("%s", err)
            }

            if deployments == 0 {
                err = p.UpdateRows(db, url, "deployments", deployments)
                if err != nil {
                    log.Default().Printf("%s", err)
                }
            }
        }

        if repo.Contributors == 0 {
            contributors, err := g.GetTotalContributorsCount(owner, name)
            if err != nil {
                log.Default().Printf("%s", err)
            }

            if contributors == 0 {
                err = p.UpdateRows(db, url, "contributors", contributors)
                if err != nil {
                    log.Default().Printf("%s", err)
                }
            }
        }

        if repo.Notifications == 0 {
            notifications, err := g.GetTotalNotificationCount(owner, name)
            if err != nil {
                log.Default().Printf("%s", err)
            }

            if notifications == 0 {
                err = p.UpdateRows(db, url, "notifications", notifications)
                if err != nil {
                    log.Default().Printf("%s", err)
                }
            }
        }

        if repo.OpenPulls == 0 && repo.ClosedPulls == 0 {
            openPullRequests, closedPullRequests, err := g.FetchAllPullRequests(owner, name)
            if err != nil {
                log.Default().Printf("%s", err)
            }

            if len(openPullRequests) == 0 && len(closedPullRequests) == 0 {
                err = p.UpdateRows(db, url, "open_pulls", len(openPullRequests))
                if err != nil {
                    log.Default().Printf("%s", err)
                }

                err = p.UpdateRows(db, url, "closed_pulls", len(closedPullRequests))
                if err != nil {
                    log.Default().Printf("%s", err)
                }
            }
        }

		result, _, err := g.APIClient.Repositories.Get(g.APIClientContext, owner, name)
		if err != nil {
			log.Default().Printf("%s", err)
		}

		if result != nil {
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
		}

		remainingRepos := totalRepos - (i + 1)

		log.Default().Printf("Completed Repository %d of %d. %d repositories remaining.\n\n", i+1, totalRepos, remainingRepos)
	}
}
