package gitops

import (
	"fmt"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
)

const (
	usernameKey = "username"
	passwordKey = "password"
)

//go:generate mockery -name=GitOperator -output=automock -outpkg=automock -case=underscore
type GitOperator interface {
	Clone(s storage.Storer, worktree billy.Filesystem, o *git.CloneOptions) (*git.Repository, error)
}

type Config struct {
	RepoUrl      string
	Branch       string
	ActualCommit string
	BaseDir      string
	Secret       map[string]interface{}
}

type Manager struct {
	gitOperator GitOperator
}

func NewManager(operator GitOperator) *Manager {
	return &Manager{gitOperator: operator}
}

func (g *Manager) CheckBranchChanges(config Config) (commitHash string, changesOccurred bool, err error) {
	auth, err := convertToBasicAuth(config.Secret)
	if err != nil {
		return commitHash, changesOccurred, errors.Wrap(err, "while parsing auth fields")
	}

	repo, err := g.gitOperator.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           config.RepoUrl,
		ReferenceName: plumbing.NewBranchReferenceName(config.Branch),
		Auth:          auth,
		SingleBranch:  true,
		NoCheckout:    true,
		Tags:          git.NoTags,
	})
	if err != nil {
		return commitHash, changesOccurred, errors.Wrapf(err, "while cloning repository: %s, branch: %s", config.RepoUrl, config.Branch)
	}

	head, err := repo.Head()
	if err != nil {
		return commitHash, changesOccurred, errors.Wrapf(err, "while getting HEAD reference for repository: %s, branch: %s", config.RepoUrl, config.Branch)
	}

	commitHash = head.Hash().String()

	if commitHash != config.ActualCommit {
		changesOccurred = true
	}

	return commitHash, changesOccurred, nil
}

func convertToBasicAuth(secret map[string]interface{}) (*http.BasicAuth, error) {
	if secret == nil {
		return &http.BasicAuth{}, nil
	}

	username, ok := secret[usernameKey].(string)
	if !ok {
		return nil, fmt.Errorf("missing field %s", usernameKey)
	}

	password, ok := secret[passwordKey].(string)
	if !ok {
		return nil, fmt.Errorf("missing field %s", passwordKey)
	}

	return &http.BasicAuth{
		Username: username,
		Password: password,
	}, nil
}
