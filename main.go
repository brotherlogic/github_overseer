package main

import (
	"context"
	"log"
	"time"

	ghbclient "github.com/brotherlogic/githubridge/client"
	ghbpb "github.com/brotherlogic/githubridge/proto"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	ghclient, err := ghbclient.GetClientInternal()
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
	}

	// Get all the repos
	repos, err := ghclient.GetRepos(ctx, &ghbpb.GetReposRequest{
		User: "brotherlogic",
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Found %v repos to process", len(repos.GetRepos()))
}
