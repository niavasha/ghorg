// Package cmd holds functions associated with cloning all of a given users repos
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/niavasha/ghuser/colorlog"
	"github.com/niavasha/ghuser/config"
	"golang.org/x/oauth2"
)

func getToken() string {
	if len(config.GitHubToken) != 40 {
		colorlog.PrintInfo("Note: GHUSER_GITHUB_TOKEN not set in $HOME/.ghuser, defaulting to keychain")
		fmt.Println()
		cmd := `security find-internet-password -s github.com | grep "acct" | awk -F\" '{ print $4 }'`
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			colorlog.PrintError(fmt.Sprintf("Failed to execute command: %s", cmd))
		}

		token := strings.TrimSuffix(string(out), "\n")

		if len(token) != 40 {
			log.Fatal("Could not find a GitHub token in keychain, create token, set GITHUB_TOKEN in your $HOME/.ghuser")
		}

		return token
	}

	return config.GitHubToken
}

// TODO: Figure out how to use go channels for this
func getAllUserCloneUrls() ([]string, error) {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getToken()},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListByUserOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 100, Page: 0},
	}

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByUser(context.Background(), os.Args[1], opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	cloneUrls := []string{}

	for _, repo := range allRepos {
		if config.GhuserCloneProtocol == "https" {
			cloneUrls = append(cloneUrls, *repo.CloneURL)
		} else {
			cloneUrls = append(cloneUrls, *repo.SSHURL)
		}
	}

	return cloneUrls, nil
}

func createDirIfNotExist() {
	if _, err := os.Stat(config.AbsolutePathToCloneTo + os.Args[1] + "_ghuser"); os.IsNotExist(err) {
		err = os.MkdirAll(config.AbsolutePathToCloneTo, 0700)
		if err != nil {
			panic(err)
		}
	}
}

func repoExistsLocally(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func getAppNameFromURL(url string) string {
	withGit := strings.Split(url, "/")
	appName := withGit[len(withGit)-1]
	split := strings.Split(appName, ".")
	return strings.Join(split[0:len(split)-1], ".")
}

func printRemainingMessages(infoMessages []error, errors []error) {
	if len(infoMessages) > 0 {
		fmt.Println()
		colorlog.PrintInfo("============ Info ============")
		fmt.Println()
		for _, i := range infoMessages {
			colorlog.PrintInfo(i)
		}
		fmt.Println()
	}

	if len(errors) > 0 {
		fmt.Println()
		colorlog.PrintError("============ Issues ============")
		fmt.Println()
		for _, e := range errors {
			colorlog.PrintError(e)
		}
		fmt.Println()
	}
}

// CloneAllReposByUser clones all repos for a given user
func CloneAllReposByUser() {
	resc, errc, infoc := make(chan string), make(chan error), make(chan error)

	createDirIfNotExist()

	if config.GhuserBranch != "master" {
		colorlog.PrintSubtleInfo("***********************************************************")
		colorlog.PrintSubtleInfo("* Ghuser will be running on branch: " + config.GhuserBranch)
		colorlog.PrintSubtleInfo("* To change back to master run $ export GHUSER_BRANCH=master")
		colorlog.PrintSubtleInfo("***********************************************************")
		fmt.Println()
	}

	cloneTargets, err := getAllUserCloneUrls()

	if err != nil {
		colorlog.PrintError(err)
	} else {
		colorlog.PrintInfo(strconv.Itoa(len(cloneTargets)) + " repos found in " + os.Args[1])
		fmt.Println()
	}

	for _, target := range cloneTargets {
		appName := getAppNameFromURL(target)

		go func(repoUrl string, branch string) {
			repoDir := config.AbsolutePathToCloneTo + os.Args[1] + "_ghuser" + "/" + appName

			if repoExistsLocally(repoDir) == true {

				cmd := exec.Command("git", "checkout", branch)
				cmd.Dir = repoDir
				err := cmd.Run()
				if err != nil {
					infoc <- fmt.Errorf("Could not checkout out %s, no changes made Repo: %s Error: %v", branch, repoUrl, err)
					return
				}

				cmd = exec.Command("git", "clean", "-f", "-d")
				cmd.Dir = repoDir
				err = cmd.Run()
				if err != nil {
					errc <- fmt.Errorf("Problem running git clean: %s Error: %v", repoUrl, err)
					return
				}

				cmd = exec.Command("git", "fetch", "-n", "origin", branch)
				cmd.Dir = repoDir
				err = cmd.Run()
				if err != nil {
					errc <- fmt.Errorf("Problem trying to fetch %v Repo: %s Error: %v", branch, repoUrl, err)
					return
				}

				cmd = exec.Command("git", "reset", "--hard", "origin/"+branch)
				cmd.Dir = repoDir
				err = cmd.Run()
				if err != nil {
					errc <- fmt.Errorf("Problem resetting %s Repo: %s Error: %v", branch, repoUrl, err)
					return
				}
			} else {
				cmd := exec.Command("git", "clone", repoUrl, repoDir)
				err := cmd.Run()
				if err != nil {
					errc <- fmt.Errorf("Problem trying to clone Repo: %s Error: %v", repoUrl, err)
					return
				}
			}

			resc <- repoUrl
		}(target, config.GhuserBranch)
	}

	errors := []error{}
	infoMessages := []error{}

	for i := 0; i < len(cloneTargets); i++ {
		select {
		case res := <-resc:
			colorlog.PrintSuccess("Success " + res)
		case err := <-errc:
			errors = append(errors, err)
		case info := <-infoc:
			infoMessages = append(infoMessages, info)
		}
	}

	printRemainingMessages(infoMessages, errors)

	colorlog.PrintSuccess(fmt.Sprintf("Finished! %s%s_ghuser", config.AbsolutePathToCloneTo, os.Args[1]))
}

// TODO: Clone via http or ssh flag

// Could clone all repos on a user
// users, _, err := client.Useranizations.List(context.Background(), "willnorris", nil)
