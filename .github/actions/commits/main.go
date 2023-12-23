package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

func main() {
	token := flag.String("INPUTS_ACCESS_TOKEN", "", "Access token")
	owner := flag.String("INPUTS_REPO_OWNER", "", "Owner/Organization of the repostory")
	repo := flag.String("INPUTS_REPO_NAME", "", "Name of the repository")
	/*
		// get environment variables
		token := os.Getenv("INPUT_ACCESS_TOKEN")
		if token == "" {
			log.Fatal("Error loading ACCESS_TOKEN")
		}
		owner := os.Getenv("INPUT_REPO_OWNER")
		if owner == "" {
			log.Fatal("Error loading REPO_OWNER")
		}
		repo := os.Getenv("INPUT_REPO_NAME")
		if repo == "" {
			log.Fatal("Error loading REPO_NAME")
		}
	*/
	ctx := context.Background()
	//	client := github.NewClient(nil)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	branches, _, err := client.Repositories.ListBranches(ctx, *owner, *repo, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	seenCommits := make(map[string]bool)

	for _, branch := range branches {
		commits, _, err := client.Repositories.ListCommits(ctx, *owner, *repo, &github.CommitsListOptions{
			SHA:   branch.GetName(),
			Since: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Until: time.Now(),
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(branch.GetName())
		for _, commit := range commits {
			printCommitInfo(commit)
			commitSHA := commit.GetSHA()
			if _, exists := seenCommits[commitSHA]; exists {
				fmt.Printf("Duplicate commit SHA: %s\n", commitSHA)
			} else {
				seenCommits[commitSHA] = true
				fmt.Printf("New commit SHA: %s\n", commitSHA)
			}
		}
	}
}

func printCommitInfo(commit *github.RepositoryCommit) {
	fmt.Printf("SHA: %s\n", commit.GetSHA())
	fmt.Printf("Author: %s\n", commit.GetAuthor().GetLogin())
	fmt.Printf("Message: %s\n", commit.GetCommit().GetMessage())
	fmt.Printf("Date: %s\n", commit.GetCommit().GetAuthor().GetDate())
	fmt.Printf("URL: %s\n", commit.GetHTMLURL())
	fmt.Println("----")
}
