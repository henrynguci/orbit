package tui

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func getAuthMethod(url string) transport.AuthMethod {
	if strings.HasPrefix(url, "git@") || strings.Contains(url, "ssh://") {
		auth, err := ssh.NewSSHAgentAuth("git")
		if err == nil {
			return auth
		}
	}
	return nil
}

func cloneRepository(url, path string) error {
	auth := getAuthMethod(url)

	opts := &git.CloneOptions{
		URL:      url,
		Progress: nil,
		Depth:    1,
	}

	if auth != nil {
		opts.Auth = auth
	}

	_, err := git.PlainClone(path, false, opts)
	return err
}
