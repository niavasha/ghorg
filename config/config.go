// Package config sets environment variables to vars for simplified sharing of environment variables between packages
package config

var (
	// GitHubToken used to auth to github, either comes from keychain locally or from the .env
	GitHubToken string
	// AbsolutePathToCloneTo Path to which ghuser will create a new folder to place all cloned repos
	AbsolutePathToCloneTo string
	// GhuserBranch branch that ghuser will checkout
	GhuserBranch string
	// GhuserCloneProtocol is the protocol to clone with
	GhuserCloneProtocol string
)
