package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	withV3()
}

func withV3() {
	ctx := context.Background()
	client := github.NewClient(authenticatedHTTPClient(ctx))
	events, _, err := client.Activity.ListRepositoryEvents(ctx, "rchampourlier", "myenv", &github.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, e := range events {
		fmt.Printf("%s - %s: %s\n", e.CreatedAt, *e.Type, string(*e.RawPayload))
		if *e.Type == "PushEvent" {
			pp, err := e.ParsePayload()
			pe := pp.(*github.PushEvent)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(*pe)
		}
	}
}

func withV4() {
	ctx := context.Background()
	client := githubv4.NewClient(authenticatedHTTPClient(ctx))

	var query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
	}
	err := client.Query(context.Background(), &query, nil)
	if err != nil {
	}
	fmt.Println("    Login:", query.Viewer.Login)
	fmt.Println("CreatedAt:", query.Viewer.CreatedAt)
}

func authenticatedHTTPClient(ctx context.Context) *http.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, src)
	return httpClient
}
