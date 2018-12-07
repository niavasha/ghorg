package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/niavasha/ghuser/cmd"
	"github.com/niavasha/ghuser/colorlog"
	"github.com/niavasha/ghuser/config"
)

func init() {
	if len(os.Args) <= 1 {
		log.Fatal("You must provide an user to clone from")
	}

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal("Error trying to find users home directory")
	}

	err = godotenv.Load(home + "/.ghuser")
	if err != nil {
		fmt.Println()
		colorlog.PrintSubtleInfo("Could not find a $HOME/.ghuser proceeding with defaults")
	}

	config.GitHubToken = os.Getenv("ghuser_GITHUB_TOKEN")
	config.AbsolutePathToCloneTo = os.Getenv("ghuser_ABSOLUTE_PATH_TO_CLONE_TO")
	config.GhuserBranch = os.Getenv("ghuser_BRANCH")
	config.GhuserCloneProtocol = os.Getenv("ghuser_CLONE_PROTOCOL")

	if config.GhuserBranch == "" {
		config.GhuserBranch = "master"
	}

	if config.AbsolutePathToCloneTo == "" {
		config.AbsolutePathToCloneTo = home + "/Desktop/"
	}

	if config.GhuserCloneProtocol == "" {
		config.GhuserCloneProtocol = "https"
	}

	withTrailingSlash := ensureTrailingSlash(config.AbsolutePathToCloneTo)
	config.AbsolutePathToCloneTo = withTrailingSlash
}

func ensureTrailingSlash(path string) string {
	if string(path[len(path)-1]) == "/" {
		return path
	}

	return path + "/"
}

func asciiTime() {
	colorlog.PrintInfo(
		`
 +-+-+-+-+ +-+-+ +-+-+-+-+-+-+
 |T|I|M|E| |T|O| |G|H|U|S|E|R|
 +-+-+-+-+ +-+-+ +-+-+-+-+-+-+
`)
}

func main() {
	asciiTime()
	cmd.CloneAllReposByUser()
}
