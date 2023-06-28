package run

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

var (
	OutPutPath      string = "."
	GithubServerURL string = "https://github.com/"
)

func getFileFromRef(repository *git.Repository, ref *plumbing.Hash, filename string) ([]byte, error) {

	// Get the commit object corresponding to ref
	commit, err := repository.CommitObject(*ref)
	if err != nil {
		return nil, err
	}

	// Get the tree object from the commit
	tree, err := repository.TreeObject(commit.TreeHash)
	if err != nil {
		return nil, err
	}

	// Get the file from the tree
	obj, err := tree.File(filename)
	if err != nil {
		return nil, err
	}

	reader, err := obj.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data := make([]byte, obj.Size)
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GitGet downloads a file from a git repository
func GitGet(repo, ref, filename string) error {

	url := fmt.Sprintf(GithubServerURL + repo + ".git")

	memFs := memfs.New()
	r, err := git.Clone(memory.NewStorage(), memFs, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		fmt.Println(err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		// Hash: plumbing.NewHash(ref),
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", ref)),
		Force:  true,
	})
	if err != nil {
		return err
	}

	pRef, err := r.Head()
	if err != nil {
		return err
	}

	hash := pRef.Hash()
	bytes, err := getFileFromRef(r, &hash, filename)
	if err != nil {
		return err
	}

	f, err := os.Create(OutPutPath + "/" + filepath.Base(filename))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(string(bytes))
	if err != nil {
		return err
	}

	return nil
}
