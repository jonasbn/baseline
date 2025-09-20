package types

import (
	"testing"
	"time"
)

func TestRepository(t *testing.T) {
	repo := Repository{
		Name:        "test-repo",
		FullName:    "testorg/test-repo",
		CloneURL:    "https://github.com/testorg/test-repo.git",
		SSHURL:      "git@github.com:testorg/test-repo.git",
		HTTPSURL:    "https://github.com/testorg/test-repo",
		Description: "A test repository",
		Private:     false,
		UpdatedAt:   time.Now(),
		Language:    "Go",
		Owner:       "testorg",
	}

	if repo.Name != "test-repo" {
		t.Errorf("Expected repo name 'test-repo', got '%s'", repo.Name)
	}

	if repo.Owner != "testorg" {
		t.Errorf("Expected repo owner 'testorg', got '%s'", repo.Owner)
	}
}

func TestCloneOptions(t *testing.T) {
	options := CloneOptions{
		Directory:    "./test",
		Organization: "testorg",
		Threads:      4,
		Verbose:      true,
		Source:       "github",
	}

	if options.Directory != "./test" {
		t.Errorf("Expected directory './test', got '%s'", options.Directory)
	}

	if options.Threads != 4 {
		t.Errorf("Expected threads 4, got %d", options.Threads)
	}
}

func TestSourceType(t *testing.T) {
	github := SourceGitHub
	bitbucket := SourceBitbucket

	if string(github) != "github" {
		t.Errorf("Expected SourceGitHub to be 'github', got '%s'", string(github))
	}

	if string(bitbucket) != "bitbucket" {
		t.Errorf("Expected SourceBitbucket to be 'bitbucket', got '%s'", string(bitbucket))
	}
}
