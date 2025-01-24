package main

import (
	"context"
	"log"
	"time"

	pb "github.com/brotherlogic/github_overseer/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	ghbclient "github.com/brotherlogic/githubridge/client"
	ghbpb "github.com/brotherlogic/githubridge/proto"
	pstoreclient "github.com/brotherlogic/pstore/client"
	pspb "github.com/brotherlogic/pstore/proto"
)

const (
	CONFIG_KEY = "github.com/brotherlogic/github_overseer/config"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	ghclient, err := ghbclient.GetClientInternal()
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
	}

	// Get a pstore client
	pclient, err := pstoreclient.GetClient()
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
	}

	// Load the config
	val, err := pclient.Read(ctx, &pspb.ReadRequest{
		Key: CONFIG_KEY,
	})
	config := &pb.Config{RepoMap: map[string]string{}}
	if err != nil && status.Code(err) != codes.NotFound {
		log.Fatalf("Failure to read config: %v", err)
	}

	if err == nil {
		proto.Unmarshal(val.GetValue().GetValue(), config)
	}

	// Get all the repos
	repos, err := ghclient.GetRepos(ctx, &ghbpb.GetReposRequest{
		User: "brotherlogic",
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Found %v repos to process", len(repos.GetRepos()))

	// Validate latest hash on each repo
	for _, repo := range repos.GetRepos() {
		recordedHash := config.GetRepoMap()[repo]
		grepo, err := ghclient.GetRepo(ctx, &ghbpb.GetRepoRequest{
			User: "brotherlogic",
			Repo: repo,
		})
		if err != nil {
			log.Fatalf("Error getting repo: %v", err)
		}

		if recordedHash != grepo.GetSha1() {
			trackTasks(ctx, repo, config, ghclient)
			config.RepoMap[repo] = grepo.GetSha1()
		}
	}

	data, err := proto.Marshal(config)
	if err != nil {
		log.Fatalf("Bad marshal: %v", err)
	}
	_, err = pclient.Write(ctx, &pspb.WriteRequest{
		Key:   CONFIG_KEY,
		Value: &anypb.Any{Value: data},
	})
	log.Printf("Finished run: %v", err)
}
