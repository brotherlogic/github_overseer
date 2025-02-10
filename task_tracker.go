package main

import (
	"context"
	"fmt"
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
	log.Printf("Tracking tthe asks for %v", repo)

	files, err := client.ListFiles(ctx, &ghbpb.ListFilesRequest{
		User: "brotherlogic",
		Repo: repo,
	})

	if err != nil {
		return fmt.Errorf("Unable to list files %w", err)
	}

	for _, file := range files.GetFiles() {
		if strings.HasSuffix(file.GetName(), ".md") {
			createOrUpdateConfig(ctx, repo, file.GetName(), file.GetHash(), config)
		}
	}

	// Look for tasks
	for _, tDoc := range config.TrackedDocuments {
		err = processDocument(ctx, tDoc, client)
		if err != nil {
			return fmt.Errorf("Error processing document %w", err)
		}
	}

	return nil
}

func processDocument(ctx context.Context, tDoc *pb.TrackedDocument, client ghbclient.GithubridgeClient) error {
	// Download the doc
	data, err := client.GetFile(ctx, &ghbpb.GetFileRequest{
		User: "brotherlogic",
		Repo: tDoc.GetRepo(),
		Path: tDoc.GetPath(),
	})

	if err != nil {
		return err
	}

	file := string(data.GetContent())
	inTasks := false
	index := int32(1)
	var tasks []*pb.Task
	for _, line := range strings.Split(file, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "## Tasks") {
			inTasks = true
		}

		if inTasks {
			if strings.HasPrefix(strings.TrimSpace(line), "1.") {
				taskBody := strings.Split(strings.TrimSpace(line), "1.")[1]
				tasks = append(tasks, &pb.Task{
					Task:        taskBody,
					IndexNumber: index,
					IssueId:     -1,
				})
				index++
			}
		}
	}

	nfile, err := buildTasks(ctx, file, tasks, client)
	if err != nil {
		return err
	}

	// Write out nfile into the repo
	log.Printf("WRITING %v", nfile)

	return nil
}

func buildTasks(ctx context.Context, file string, tasks []*pb.Task, client ghbclient.GithubridgeClient) (string, error) {
	for i, task := range tasks {
		// We have a task entered for this one
		if task.IssueId <= 0 {
			// Create a task and store it
			issue, err := client.CreateIssue(ctx, &ghbpb.CreateIssueRequest{
				User:  "brotherlogic",
				Repo:  task.GetBouncedRepo(),
				Title: task.GetTask(),
			})
			if err != nil {
				return "", err

			}
			task.IndexNumber = int32(i)
			task.IssueId = int32(issue.GetIssueId())
			task.State = pb.Task_TASK_STATE_UNKNOWN
		}

		if task.State == pb.Task_TASK_STATE_IN_PROGRESS {
			// Check to see if we've finished this task
			issue, err := client.GetIssue(ctx, &ghbpb.GetIssueRequest{
				User: "brotherlogic",
				Repo: task.GetBouncedRepo(),
				Id:   task.GetIssueId()})
			if err != nil {
				return "", err
			}

			if issue.GetState() == ghbpb.IssueState_ISSUE_STATE_CLOSED {
				task.State = pb.Task_TASK_STATE_COMPLETE
			}
		}

		// Replace the task line with the correct one given the state of the tasklist
		replaceInFile(ctx, file, task.GetTask(), buildTask(task))
	}

	return file, nil
}

func buildTask(task *pb.Task) string {
	strikeout := ""
	if task.GetState() == pb.Task_TASK_STATE_COMPLETE {
		strikeout = "~~"
	}
	return fmt.Sprintf("1. %v%v [%v%v#%v]%v", strikeout, task.GetTask(), "brotherlogic", task.GetBouncedRepo(), task.GetIssueId(), strikeout)
}

func replaceInFile(ctx context.Context, file, base, replace string) (string, error) {
	var nlines []string
	for _, line := range strings.Split(file, "\n") {
		if strings.Contains(line, base) {
			nlines = append(nlines, replace)
		} else {
			nlines = append(nlines, line)
		}
	}

	return strings.Join(nlines, "\n"), nil
}
