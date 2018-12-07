package main

import (
	"testing"

	"github.com/niavasha/ghuser/config"
)

func TestDefaultBranch(t *testing.T) {
	branch := config.GhuserBranch
	if branch != "master" {
		t.Errorf("Default branch should be master")
	}
}
