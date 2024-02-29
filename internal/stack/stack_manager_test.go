package stack

import (
	"github.com/Bhacaz/gostacking/internal/color"
	"github.com/Bhacaz/gostacking/internal/git"
	"testing"
)

type gitCommandsStub struct {
	git.InterfaceCommands
}

func (g gitCommandsStub) CurrentBranchName() (string, error) {
	return "my_feature_part1", nil
}

func (g gitCommandsStub) BranchExists(branchName string) bool {
	return true
}

func (g gitCommandsStub) Checkout(branchName string) {
	// Do nothing
}

func (g gitCommandsStub) SyncBranches(branches []string, checkoutBranchEnd string, push bool) {
	// Do nothing
}

func (g gitCommandsStub) BranchDiff(baseBranch string, branch string) bool {
	return false
}

func (g gitCommandsStub) Fetch() {}

func (g gitCommandsStub) IsBehindRemote(branch string) bool { return false }

func TestCreateStack(t *testing.T) {
	stacksManager := StacksManager{
		stacksPersister: &StacksPersistingStub{},
		gitCommands:     gitCommandsStub{},
	}

	result := stacksManager.CreateStack("stack3")

	data, _ := stacksManager.stacksPersister.LoadStacks()

	// Add stack3 to the list of stacks
	if data.Stacks[2].Name != "stack3" {
		t.Errorf("got %s, want %s", data.Stacks[2].Name, "stack3")
	}

	if data.Stacks[2].Branches[0] != "my_feature_part1" {
		t.Errorf("got %s, want %s", data.Stacks[2].Branches[0], "my_feature_part1")
	}

	// Return the message for CLI
	want := "Stack created " + color.Green("stack3")
	if result != want {
		t.Errorf("got %s, want %s", result, want)
	}
}

func TestCurrentStackStatus(t *testing.T) {
	stacksManager := StacksManager{
		stacksPersister: &StacksPersistingStub{},
		gitCommands:     gitCommandsStub{},
	}

	result := stacksManager.CurrentStackStatus(false)

	want := "Current stack: " +
		color.Green("stack1") +
		"\nBranches:\n1. " +
		color.Yellow("branch1") + "\n" +
		"2. " + color.Yellow("branch2") + "\n"
	if result != want {
		t.Errorf("got %s, want %s", result, want)
	}
}
