package worker

import (
	"context"
	"sync"

	"github.com/jonasbn/baseline/internal/git"
	"github.com/jonasbn/baseline/internal/types"
)

// WorkerPool manages concurrent Git operations
type WorkerPool struct {
	numWorkers int
	gitOps     *git.GitOps
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, verbose bool) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		gitOps:     git.NewGitOps(verbose),
	}
}

// CloneRepositories clones repositories concurrently
func (wp *WorkerPool) CloneRepositories(ctx context.Context, repositories []types.Repository, targetDir string) <-chan types.CloneResult {
	resultChan := make(chan types.CloneResult, len(repositories))
	repoChan := make(chan types.Repository, len(repositories))

	// Send all repositories to the channel
	go func() {
		defer close(repoChan)
		for _, repo := range repositories {
			select {
			case repoChan <- repo:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < wp.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range repoChan {
				select {
				case <-ctx.Done():
					return
				default:
					// Skip if repository already exists
					if wp.gitOps.RepositoryExists(repo, targetDir) {
						resultChan <- types.CloneResult{
							Repository: repo,
							Success:    true,
							Error:      nil,
							Duration:   0,
						}
						continue
					}

					result := wp.gitOps.CloneRepository(repo, targetDir)
					resultChan <- result
				}
			}
		}()
	}

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

// UpdateRepositories updates repositories concurrently
func (wp *WorkerPool) UpdateRepositories(ctx context.Context, repositories []types.Repository, targetDir string) <-chan types.UpdateResult {
	resultChan := make(chan types.UpdateResult, len(repositories))
	repoChan := make(chan types.Repository, len(repositories))

	// Send all repositories to the channel
	go func() {
		defer close(repoChan)
		for _, repo := range repositories {
			select {
			case repoChan <- repo:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < wp.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range repoChan {
				select {
				case <-ctx.Done():
					return
				default:
					// Skip if repository doesn't exist
					if !wp.gitOps.RepositoryExists(repo, targetDir) {
						resultChan <- types.UpdateResult{
							Repository: repo,
							Success:    false,
							Updated:    false,
							Error:      nil,
							Duration:   0,
						}
						continue
					}

					result := wp.gitOps.UpdateRepository(repo, targetDir)
					resultChan <- result
				}
			}
		}()
	}

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}
