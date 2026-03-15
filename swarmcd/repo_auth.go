package swarmcd

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/m-adawi/swarm-cd/util"
)

func createRepoAuth(repoName string) (transport.AuthMethod, error) {
	repoConfig := config.RepoConfigs[repoName]

	if strings.HasPrefix(repoConfig.Url, "ssh://") {
		return createSSHAuth(*repoConfig)
	}

	return createHTTPBasicAuth(repoName, *repoConfig)
}

func createSSHAuth(repoConfig util.RepoConfig) (transport.AuthMethod, error) {
	username := repoConfig.Username
	if username == "" {
		username = "git"
	}

	return ssh.NewPublicKeysFromFile(repoConfig.Username, repoConfig.PasswordFile, "")
}

func createHTTPBasicAuth(repoName string, repoConfig util.RepoConfig) (*http.BasicAuth, error) {
	// assume repo is public and no auth is required
	if repoConfig.Username == "" && repoConfig.Password == "" && repoConfig.PasswordFile == "" {
		return nil, nil
	}

	if repoConfig.Username == "" {
		return nil, fmt.Errorf("you must set username for the repo %s", repoName)
	}

	if repoConfig.Password == "" && repoConfig.PasswordFile == "" {
		return nil, fmt.Errorf("you must set one of password or password_file properties for the repo %s", repoName)
	}

	var password string
	if repoConfig.Password != "" {
		password = repoConfig.Password
	} else {
		passwordBytes, err := os.ReadFile(repoConfig.PasswordFile)
		if err != nil {
			return nil, fmt.Errorf("could not read password file %s for repo %s", repoConfig.PasswordFile, repoName)
		}
		// trim newline and whitespaces
		password = strings.TrimSpace(string(passwordBytes))
	}

	return &http.BasicAuth{
		Username: repoConfig.Username,
		Password: password,
	}, nil
}
