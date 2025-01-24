package main

import (
	"context"
	"log"
	"strings"

	pb "github.com/brotherlogic/github_overseer/proto"
	ghbpb "github.com/brotherlogic/githubridge/proto"

	ghbclient "github.com/brotherlogic/githubridge/client"
)

func createOrUpdateConfig(ctx context.Context, repo, path, hash string, config *pb.Config) {
	for _, td := range config.GetTrackedDocuments() {
		if td.GetUser() == "brotherlogic" &&
			td.GetRepo() == repo &&
			td.GetPath() == path {
			//Found this in the tracked document
			td.LatestHash = hash
			return
		}
	}

	// Add this new file
	config.TrackedDocuments = append(config.TrackedDocuments, &pb.TrackedDocument{
		User:       "brotherlogic",
		Repo:       repo,
		Path:       path,
		LatestHash: hash,
	})
}

func trackTasks(ctx context.Context, repo string, config *pb.Config, client ghbclient.GithubridgeClient) error {
	log.Printf("Tracking tasks for %v", repo)

	files, err := client.ListFiles(ctx, &ghbpb.ListFilesRequest{
		User: "brotherlogic",
		Repo: repo,
	})

	if err != nil {
		return err
	}

	for _, file := range files.GetFiles() {
		if strings.HasSuffix(file, ".md") {
			createOrUpdateConfig(ctx, repo, file.GetName(), file.GetHash(), config)
		}
	}
}
