package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/google/go-github/v57/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Commit struct {
	Sha     string `json:"sha"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Date    string `json:"date"`
	Url     string `json:"url"`
	Branch  string `json:"branch"`
}

func main() {
	godotenv.Load()

	// get environment variables
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Error loading GITHUB_TOKEN")
	}

	owner := os.Getenv("REPO_OWNER")
	if owner == "" {
		log.Fatal("Error loading REPO_OWNER")
	}

	repo := os.Getenv("REPO_NAME")
	if repo == "" {
		log.Fatal("Error loading REPO_NAME")
	}

	tz := os.Getenv("LOCATION")
	if tz == "" {
		log.Fatal("Error loading TIMEZONE")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatal(err)
	}

	sinceUTCStr := os.Getenv("SINCE")
	if sinceUTCStr == "" {
		log.Fatal("Error loading SINCE")
	}
	sinceUTC, err := time.Parse("2006-01-02T15:04:05Z", sinceUTCStr)
	if err != nil {
		log.Fatal(err)
	}
	since := sinceUTC.In(loc)

	untilUTCStr := os.Getenv("UNTIL")
	if untilUTCStr == "" {
		log.Fatal("Error loading UNTIL")
	}
	untilUTC, err := time.Parse("2006-01-02T15:04:05Z", untilUTCStr)
	if err != nil {
		log.Fatal(err)
	}
	until := untilUTC.In(loc)

	commits, err := gh_commits(token, owner, repo, since, until, loc)
	if err != nil {
		log.Fatal(err)
	}

	df := dataframe.LoadStructs(*commits)
	groupd := df.GroupBy("Date", "Author").
		Aggregation([]dataframe.AggregationType{dataframe.Aggregation_COUNT}, []string{"Date"})

	g := groupd.Maps()
	//	fmt.Println(g)

	ndf := dataframe.LoadMaps(g)

	fmt.Println(since, until, ndf)
}

func gh_commits(token, owner, repo string, since, until time.Time, loc *time.Location) (*[]Commit, error) {
	commitArray := make([]Commit, 0)

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	branches, _, err := client.Repositories.ListBranches(ctx, owner, repo, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	seenCommits := make(map[string]bool)

	for _, branch := range branches {
		commits, _, err := client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
			SHA:   branch.GetName(),
			Since: since,
			Until: until,
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, commit := range commits {
			commitSHA := commit.GetSHA()
			if _, exists := seenCommits[commitSHA]; exists {
				continue
			} else {
				seenCommits[commitSHA] = true
				commitData := Commit{
					Sha:     commit.GetSHA(),
					Author:  commit.GetAuthor().GetLogin(),
					Message: commit.GetCommit().GetMessage(),
					Date:    commit.GetCommit().GetAuthor().GetDate().In(loc).Format("2006-01-02"),
					Url:     commit.GetHTMLURL(),
					Branch:  branch.GetName(),
				}
				commitArray = append(commitArray, commitData)
			}
		}
	}
	return &commitArray, nil
}

/*
func printCommitInfo(commit *github.RepositoryCommit) {
	fmt.Printf("SHA: %s\n", commit.GetSHA())
	fmt.Printf("Author: %s\n", commit.GetAuthor().GetLogin())
	fmt.Printf("Message: %s\n", commit.GetCommit().GetMessage())
	fmt.Printf("Date: %s\n", commit.GetCommit().GetAuthor().GetDate())
	fmt.Printf("URL: %s\n", commit.GetHTMLURL())
	fmt.Println("----")
}
*/
