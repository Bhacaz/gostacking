package stack

import (
	"fmt"
	"github.com/Bhacaz/gostacking/internal/color"
	"github.com/Bhacaz/gostacking/internal/git"
	"strings"
	"testing"
)

type PrinterStub struct {
	MessageReceived *[]string
}

func (p PrinterStub) Println(a ...interface{}) {
	strs := make([]string, len(a))
	for i, v := range a {
		strs[i] = fmt.Sprint(v)
	}
	*p.MessageReceived = append(*p.MessageReceived, strings.Join(strs, " "))
	*p.MessageReceived = append(*p.MessageReceived, "\n")
}

type gitExecutorStub struct {
	stubExec func(...string) (string, error)
}

func (g gitExecutorStub) Exec(command ...string) (string, error) {
	return g.stubExec(command...)
}

func (sm StacksManager) printerMessage() string {
	return strings.Join(*sm.printer.(PrinterStub).MessageReceived, "")
}

func StacksManagerForTest(gitExecutor git.InterfaceGitExecutor, messageReceived *[]string) StacksManager {
	return StacksManager{
		stacksPersister: &StacksPersistingStub{},
		gitExecutor:     gitExecutor,
		printer: PrinterStub{
			MessageReceived: messageReceived,
		},
	}
}

func TestCreateStack(t *testing.T) {
	gitExecutor := gitExecutorStub{
		stubExec: func(command ...string) (string, error) {
			return "my_feature_part1", nil
		},
	}
	messageReceived := []string{}
	stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

	result := stacksManager.CreateStack("stack3")
	data, _ := stacksManager.stacksPersister.LoadStacks()

	// Add stack3 to the list of stacks
	if data.Stacks[2].Name != "stack3" {
		t.Errorf("got %s, want %s", data.Stacks[2].Name, "stack3")
	}

	if data.Stacks[2].Branches[0] != "my_feature_part1" {
		t.Errorf("got %s, want %s", data.Stacks[2].Branches[0], "my_feature_part1")
	}

	want := "Stack created " + color.Green("stack3")
	if !strings.Contains(stacksManager.printerMessage(), want) {
		t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
	}

	if result != nil {
		t.Errorf("got Error %s, want none", result)
	}
}

//
//func TestCurrentStackStatus(t *testing.T) {
//	stacksManager := StacksManager{
//		stacksPersister: &StacksPersistingStub{},
//		gitCommands:     gitCommandsStub{},
//	}
//
//	result := stacksManager.CurrentStackStatus(false)
//
//	want := "Current stack: " +
//		color.Green("stack1") +
//		"\nBranches:\n1. " +
//		color.Yellow("branch1") + "\n" +
//		"2. " + color.Yellow("branch2") + "\n"
//	if result != want {
//		t.Errorf("got %s, want %s", result, want)
//	}
//}
