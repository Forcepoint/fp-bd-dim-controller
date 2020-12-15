package backup

import (
	"errors"
	"fmt"
	structs2 "fp-dynamic-elements-manager-controller/internal/backup/structs"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"os"
	"time"
)

type HistoryCommitter interface {
	Commit(string, int64) error
	RestoreToPoint(string) error
	ListHistory() ([]structs2.History, error)
}

type GitController struct {
	repo   *git.Repository
	logger *structs.AppLogger
}

func NewGitController(logger *structs.AppLogger) *GitController {
	return &GitController{logger: logger, repo: initRepo(logger)}
}

func initRepo(logger *structs.AppLogger) *git.Repository {
	var r *git.Repository
	fs := osfs.New(os.Getenv("DB_BACKUP_DIR"))
	dot, _ := fs.Chroot(".git")
	storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	r, err := git.Init(storage, fs)

	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			r, err = git.Open(storage, fs)
			if err != nil {
				logger.SystemLogger.Error(err, "error opening repo")
				return nil
			}
		} else {
			logger.SystemLogger.Error(err, "error initialising repo")
			return nil
		}
	}

	return r
}

func (g *GitController) Commit(msg string, elementCount int64) (err error) {
	if g.repo == nil {
		g.logger.SystemLogger.Error(errors.New("git repository is nil"), "error getting instance of git repository")
		return
	}
	w, err := g.repo.Worktree()
	if err != nil {
		g.logger.SystemLogger.Error(err, "error getting worktree")
		return
	}

	now := time.Now()
	_, err = w.Commit(fmt.Sprintf("%s %s : %d Elements", msg, now.Format(time.RFC822), elementCount), &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  os.Getenv("DB_BACKUP_NAME"),
			Email: os.Getenv("DB_BACKUP_EMAIL"),
			When:  now,
		},
	})

	if err != nil {
		g.logger.SystemLogger.Error(err, "error committing worktree")
		return
	}

	return nil
}

func (g *GitController) RestoreToPoint(commitHash string) (err error) {
	wt, err := g.repo.Worktree()

	if err != nil {
		return
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})

	if err != nil {
		return
	}

	return nil
}

func (g *GitController) ListHistory() (commits []structs2.History, err error) {
	// ... retrieves the branch pointed by HEAD
	ref, err := g.repo.Head()

	if err != nil {
		return
	}

	// ... retrieves the commit history
	cIter, err := g.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return
	}

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, structs2.History{
			Hash:    c.Hash.String(),
			Message: c.Message,
		})
		return nil
	})

	return
}
