package stack

import (
	"github.com/Bhacaz/gostacking/internal/color"
	"github.com/Bhacaz/gostacking/internal/git"
	"testing"
)

// TODO replace GitCmd with a mock
func TestNew(t *testing.T) {
	stacksManager := StacksManager{
		stacksPersister: &StacksPersistingStub{},
		gitCmd:          git.Cmd(),
	}

	result := stacksManager.CreateStack("stack3")

	data, _ := stacksManager.stacksPersister.LoadStacks()

	// Add stack3 to the list of stacks
	if data.Stacks[2].Name != "stack3" {
		t.Errorf("got %s, want %s", data.Stacks[2].Name, "stack3")
	}

	// Return the message for CLI
	want := "CreateStack stack created " + color.Green("stack3")
	if result != want {
		t.Errorf("got %s, want %s", result, want)
	}
}
