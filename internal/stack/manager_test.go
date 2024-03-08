package stack

import (
	"fmt"
	"github.com/Bhacaz/gostacking/internal/color"
	"github.com/Bhacaz/gostacking/internal/git"
	"reflect"
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
		stacks:      stacksDataMock(),
		gitExecutor: gitExecutor,
		printer: PrinterStub{
			MessageReceived: messageReceived,
		},
	}
}

func TestCreateStack(t *testing.T) {
	t.Run("create stack", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "my_feature_part1", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.CreateStack("stack3")
		stacksManager.stacks.LoadStacks()
		data := *stacksManager.stacks

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
	})

	t.Run("when currentBranchName return error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "", fmt.Errorf("git command error")
			},
		}
		messageReceived := []string{}
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.CreateStack("stack3")

		if result == nil {
			t.Errorf("got none, want Error")
		}
	})
}

func TestAddBranch(t *testing.T) {
	t.Run("when passing empty string", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch command[0] {
				case "rev-parse":
					return "my_current_branch", nil
				default:
					t.Errorf("unwanted git command: %s", command[0])
					return "", nil
				}
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.AddBranch("")

		want := "Branch " + color.Yellow("my_current_branch") + " added to " + color.Green("stack1")
		got := stacksManager.printerMessage()
		if !strings.Contains(got, want) {
			t.Errorf("got \"%s\", want \"%s\"", got, want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("when passing empty string AND currentBranchName return error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch command[0] {
				case "rev-parse":
					return "", fmt.Errorf("git command error")
				default:
					t.Errorf("unwanted git command: %s", command[0])
					return "", nil
				}
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.AddBranch("")

		if result == nil {
			t.Errorf("got none, want Error")
		}
	})

	t.Run("when passing branch name", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch strings.Join(command, " ") {
				case "rev-parse --verify my_branch":
					return "", nil
				default:
					t.Errorf("unwanted git command: %s", command)
					return "", nil
				}
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.AddBranch("my_branch")

		want := "Branch " + color.Yellow("my_branch") + " added to " + color.Green("stack1")
		got := stacksManager.printerMessage()
		if !strings.Contains(got, want) {
			t.Errorf("got \"%s\", want \"%s\"", got, want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("when passing branch name that does not exist", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch strings.Join(command, " ") {
				case "rev-parse --verify non_existing_branch":
					return "", fmt.Errorf("branch does not exist")
				default:
					t.Errorf("unwanted git command: %s", command)
					return "", nil
				}
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.AddBranch("non_existing_branch")
		want := "Branch " + color.Yellow("non_existing_branch") + " does not exist"
		got := stacksManager.printerMessage()
		if !strings.Contains(got, want) {
			t.Errorf("got \"%s\", want \"%s\"", got, want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})
}

func TestList(t *testing.T) {
	t.Run("list stacks", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.List()
		want := "Current stack: " + color.Green("stack1")
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})
}

func TestListStacksForCompletion(t *testing.T) {
	t.Run("list stacks for completion", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.ListStacksForCompletion("st")
		want := []string{"stack1", "stack2"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("got %s, want %s", result, want)
		}
	})
}

func TestSwitchByName(t *testing.T) {
	t.Run("switch stack by name", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.SwitchByName("stack2")
		want := "Switched to stack " + color.Green("stack2")
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
		data := *stacksManager.stacks
		if stacksManager.stacks.CurrentStack != "stack2" {
			t.Errorf("got %s, want %s", data.CurrentStack, "stack2")
		}
	})

	t.Run("switch stack with empty args using currentBranchName", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "branch3", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.SwitchByName("")
		want := "Switched to stack " + color.Green("stack2")
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
		data := *stacksManager.stacks
		if stacksManager.stacks.CurrentStack != "stack2" {
			t.Errorf("got %s, want %s", data.CurrentStack, "stack2")
		}
	})

	t.Run("switch stack with empty args using currentBranchName and currentBranchName return error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "", fmt.Errorf("git command error")
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.SwitchByName("")

		if result == nil {
			t.Errorf("got none, want Error")
		}
	})

	t.Run("switch stack with empty args and no stack was found for the currentBranchName", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "non_existing_branch", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.SwitchByName("")

		if result == nil {
			t.Errorf("got none, want Error")
		}
	})

	t.Run("switch stack by name when stack does not exist", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.SwitchByName("non_existing_stack")

		if result == nil {
			t.Errorf("show have no error, got %s", result)
		}
	})
}

func TestSwitchByNumber(t *testing.T) {
	t.Run("switch stack by number", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.SwitchByNumber(2)
		want := "Switched to stack " + color.Green("stack2")
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
		data := *stacksManager.stacks
		if stacksManager.stacks.CurrentStack != "stack2" {
			t.Errorf("got %s, want %s", data.CurrentStack, "stack2")
		}
	})

	t.Run("switch stack by number when number is invalid", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.SwitchByNumber(3)
		if result == nil {
			t.Errorf("got none, want Error")
		}
	})
}

func TestRemoveByName(t *testing.T) {
	t.Run("remove branch by name", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.RemoveByName("branch1")
		want := "Branch " + color.Yellow("branch1") + " removed from " + color.Green("stack1")
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
		data := *stacksManager.stacks
		got := len(data.Stacks[0].Branches)
		if got != 1 {
			t.Errorf("got %d, want %d", got, 1)
		}

		if data.Stacks[0].Branches[0] != "branch2" {
			t.Errorf("got %s, want %s", data.Stacks[0].Branches[0], "branch2")
		}
	})

	t.Run("remove branch by name when branch does not exist", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.RemoveByName("non_existing_branch")

		if result == nil {
			t.Errorf("show have no error, got %s", result)
		}
	})
}

func TestStacksManager_RemoveByNumber(t *testing.T) {
	t.Run("remove branch by number", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.RemoveByNumber(1)
		want := "Branch " + color.Yellow("branch1") + " removed from stack " + color.Green("stack1")
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
		data := *stacksManager.stacks
		got := len(data.Stacks[0].Branches)
		if got != 1 {
			t.Errorf("got %d, want %d", got, 1)
		}

		if data.Stacks[0].Branches[0] != "branch2" {
			t.Errorf("got %s, want %s", data.Stacks[0].Branches[0], "branch2")
		}
	})

	t.Run("remove branch by number when number is invalid", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.RemoveByNumber(3)
		if result == nil {
			t.Errorf("got none, want Error")
		}
	})
}

func TestStacksManager_Delete(t *testing.T) {
	t.Run("delete stack", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		err := stacksManager.Delete("stack1")
		want := "Stack " + color.Green("stack1") + " deleted"
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}
		data := *stacksManager.stacks
		got := len(data.Stacks)
		if got != 1 {
			t.Errorf("got %d, want %d", got, 1)
		}
	})

	t.Run("delete stack when stack does not exist", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		result := stacksManager.Delete("non_existing_stack")

		if result == nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("delete the last stack", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		err := stacksManager.Delete("stack1")
		err = stacksManager.Delete("stack2")
		want := "Stack " + color.Green("stack2") + " deleted"
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}
		data := *stacksManager.stacks
		got := len(data.Stacks)
		if got != 0 {
			t.Errorf("got %d, want %d", got, 0)
		}
	})
}

func TestStacksManager_CheckoutByName(t *testing.T) {
	t.Run("checkout branch by name", func(t *testing.T) {
		var gitCommandsReceived []string

		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch strings.Join(command, " ") {
				case "checkout branch1":
					gitCommandsReceived = append(gitCommandsReceived, strings.Join(command, " "))
					return "Switched to branch branch1", nil
				case "rev-parse --verify branch1":
					gitCommandsReceived = append(gitCommandsReceived, strings.Join(command, " "))
					return "", nil
				default:
					t.Errorf("unwanted git command: %s", command)
				}
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.CheckoutByName("branch1")

		if gitCommandsReceived[0] != "rev-parse --verify branch1" {
			t.Errorf("got %s, want %s", gitCommandsReceived[0], "rev-parse --verify branch1")
		}

		if gitCommandsReceived[1] != "checkout branch1" {
			t.Errorf("got %s, want %s", gitCommandsReceived[1], "checkout branch1")
		}

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}
	})

	t.Run("when branch does not exist", func(t *testing.T) {
		var gitCommandsReceived []string

		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch strings.Join(command, " ") {
				case "rev-parse --verify non_existing_branch":
					gitCommandsReceived = append(gitCommandsReceived, strings.Join(command, " "))
					return "", fmt.Errorf("branch does not exist")
				default:
					t.Errorf("unwanted git command: %s", command)
				}
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.CheckoutByName("non_existing_branch")

		if gitCommandsReceived[0] != "rev-parse --verify non_existing_branch" {
			t.Errorf("got %s, want %s", gitCommandsReceived[0], "rev-parse --verify non_existing_branch")
		}

		if err == nil {
			t.Errorf("got none, want Error")
		}
	})

	t.Run("when checkout return an error", func(t *testing.T) {
		var gitCommandsReceived []string

		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch strings.Join(command, " ") {
				case "rev-parse --verify branch1":
					gitCommandsReceived = append(gitCommandsReceived, strings.Join(command, " "))
					return "", nil
				case "checkout branch1":
					gitCommandsReceived = append(gitCommandsReceived, strings.Join(command, " "))
					return "", fmt.Errorf("checkout error")
				default:
					t.Errorf("unwanted git command: %s", command)
				}
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.CheckoutByName("branch1")

		if gitCommandsReceived[0] != "rev-parse --verify branch1" {
			t.Errorf("got %s, want %s", gitCommandsReceived[0], "rev-parse --verify branch1")
		}

		if gitCommandsReceived[1] != "checkout branch1" {
			t.Errorf("got %s, want %s", gitCommandsReceived[1], "checkout branch1")
		}

		if err == nil {
			t.Errorf("got none, want Error")
		}
	})
}
