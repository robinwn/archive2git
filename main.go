package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: archive2git <directory_path> <personal_access_token> <author_email>")
		os.Exit(1)
	}

	dir := os.Args[1]
	pat := os.Args[2]
	authorEmail := os.Args[3]

	repo, err := git.PlainOpen(dir)
	if err != nil {
		fmt.Printf("Error opening repository: %v\n", err)
		os.Exit(1)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Printf("Error getting worktree: %v\n", err)
		os.Exit(1)
	}

	status, err := worktree.Status()
	if err != nil {
		fmt.Printf("Error getting status: %v\n", err)
		os.Exit(1)
	}

	if status.IsClean() {
		fmt.Println("No changes to commit")
		os.Exit(0)
	}

	_, err = worktree.Add(".")
	if err != nil {
		fmt.Printf("Error adding changes: %v\n", err)
		os.Exit(1)
	}

	// Get current date and time
	now := time.Now()
	commitMessage := fmt.Sprintf("Auto-commit: Update files - %s", now.Format("2006-01-02 15:04:05"))

	commit, err := worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "archive2git",
			Email: authorEmail,
			When:  now,
		},
	})
	if err != nil {
		fmt.Printf("Error committing changes: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Changes committed: %s\n", commit.String())

	// Get the remote URL
	remote, err := repo.Remote("origin")
	if err != nil {
		fmt.Printf("Error getting remote: %v\n", err)
		os.Exit(1)
	}

	remoteURL := remote.Config().URLs[0]

	// Modify the remote URL to include the PAT
	parts := strings.SplitN(remoteURL, "://", 2)
	if len(parts) != 2 {
		fmt.Println("Invalid remote URL format")
		os.Exit(1)
	}

	authRemoteURL := fmt.Sprintf("https://%s:%s@%s", "PAT", pat, parts[1])

	// Push changes
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "PAT", // This can be any non-empty string
			Password: pat,
		},
		RemoteURL: authRemoteURL,
	})
	if err != nil {
		fmt.Printf("Error pushing changes: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Changes pushed successfully")
}