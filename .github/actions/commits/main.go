package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Output struct {
	SHA     string `json:"sha"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Date    string `json:"date"`
	URL     string `json:"url"`
}

func main() {
	godotenv.Load(".env")
	// get environment variables
	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Error loading ACCESS_TOKEN")
	}
	fromDate := os.Getenv("FROM_DATE")
	if fromDate == "" {
		log.Fatal("Error loading FROM_DATE")
	}
	toDate := os.Getenv("TO_DATE")
	if toDate == "" {
		log.Fatal("Error loading TO_DATE")
	}
	owner := os.Getenv("REPO_OWNER")
	if owner == "" {
		log.Fatal("Error loading REPO_OWNER")
	}
	repo := os.Getenv("REPO_NAME")
	if repo == "" {
		log.Fatal("Error loading REPO_NAME")
	}

	fmt.Println(fromDate, toDate, owner, repo)

	ctx := context.Background()
	//	client := github.NewClient(nil)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	branches, _, err := client.Repositories.ListBranches(ctx, owner, repo, nil)
	if err != nil {
		log.Fatal(err)
	}

	seenCommits := make(map[string]bool)
	output := []Output{}

	since, err := time.Parse("2006-01-02 03:04:05", fromDate+" 00:00:00")
	if err != nil {
		log.Fatal("fromDate parse:", err)
	}
	until, err := time.Parse("2006-01-02 03:04:05", toDate+" 23:59:59")
	if err != nil {
		log.Fatal("toDate parse:", err)
	}

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

		fmt.Println(branch.GetName())
		for _, commit := range commits {
			printCommitInfo(commit)
			commitSHA := commit.GetSHA()
			if _, exists := seenCommits[commitSHA]; exists {
				fmt.Printf("Duplicate commit SHA: %s\n", commitSHA)
			} else {
				seenCommits[commitSHA] = true
				fmt.Printf("New commit SHA: %s\n", commitSHA)
				output = append(output, setOutput(commit))
			}
		}
	}
	json, err := json.Marshal(output)
	if err != nil {
		log.Fatal("json marshal", err)
	}
	// ファイルとして出力
	file, err := os.Create("commit_metrics.json")
	if err != nil {
		log.Fatal("file create:", err)
	}
	defer file.Close()
	// JSONテキストとして書き込み
	_, err = file.Write(json)
	if err != nil {
		log.Fatal("file write:", err)
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

func setOutput(commit *github.RepositoryCommit) Output {
	return Output{
		SHA:     commit.GetSHA(),
		Author:  commit.GetAuthor().GetLogin(),
		Message: commit.GetCommit().GetMessage(),
		Date:    commit.GetCommit().GetAuthor().GetDate().Format("2006-01-02 15:04:05"),
		URL:     commit.GetHTMLURL(),
	}
}
