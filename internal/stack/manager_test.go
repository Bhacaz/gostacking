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

func TestStacksManager_CreateStack(t *testing.T) {
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

func TestStacksManager_CurrentStackStatus(t *testing.T) {
	t.Run("current stack status", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				// Ensure not being behind remote AND no diff with parent branch
				if strings.HasPrefix(joinedCommand, "diff") {
					return "", nil
				}
				return "something", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.CurrentStackStatus(false)

		want := fmt.Sprintf(
			`Current stack: %s
Branches:
1. %s
2. %s`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Yellow("branch2"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("current stack status with log", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				// Ensure not being behind remote AND no diff with parent branch
				if strings.HasPrefix(joinedCommand, "diff") {
					return "", nil
				} else if strings.HasPrefix(joinedCommand, "log") {
					return "log", nil
				}
				return "something", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.CurrentStackStatus(true)

		want := fmt.Sprintf(
			`Current stack: %s
Branches:
1. %s
	log
2. %s
	log`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Yellow("branch2"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("when fetch return an error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if "fetch" == command[0] {
					return "", fmt.Errorf("git command error")
				}
				t.Errorf("unwanted git command should have return: %s", command[0])
				return "something", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.CurrentStackStatus(false)

		if err == nil {
			t.Errorf("got none, want Error")
		}
	})

	t.Run("when branch1 is behind remote", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				// branch1 is behind remove
				if "diff --name-only branch1...origin/branch1" == joinedCommand {
					return "file1.txt", nil
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.CurrentStackStatus(false)

		want := fmt.Sprintf(
			`Current stack: %s
Branches:
1. %s %s
2. %s`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Teal("↓"),
			color.Yellow("branch2"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("when branch1 is behind remote and with log", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				// Ensure not being behind remote AND no diff with parent branch
				if "diff --name-only branch1...origin/branch1" == joinedCommand {
					return "file1.txt", nil
				} else if strings.HasPrefix(joinedCommand, "log") {
					return "log", nil
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.CurrentStackStatus(true)

		want := fmt.Sprintf(
			`Current stack: %s
Branches:
1. %s %s
	log
2. %s
	log`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Teal("↓"),
			color.Yellow("branch2"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("when branch1 remote return an error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				// branch1 is behind remove
				if "diff --name-only branch1...origin/branch1" == joinedCommand {
					return "", fmt.Errorf("git command error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.CurrentStackStatus(false)

		want := fmt.Sprintf(
			`Current stack: %s
Branches:
1. %s
2. %s`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Yellow("branch2"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("when branch2 has diff with parent branch", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				fmt.Println(joinedCommand)
				// branch2 has diff with parent branch
				if "diff --name-only branch2...branch1" == joinedCommand {
					return "file2.txt", nil
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.CurrentStackStatus(false)

		want := fmt.Sprintf(
			`Current stack: %s
Branches:
1. %s
2. %s %s`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Yellow("branch2"),
			color.Red("*"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}

	})

	t.Run("when branch2 has diff return an error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				fmt.Println(joinedCommand)
				if "diff --name-only branch2...branch1" == joinedCommand {
					return "", fmt.Errorf("git diff command error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.CurrentStackStatus(false)

		want := fmt.Sprintf(
			`Current stack: %s
Branches:
1. %s
2. %s`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Yellow("branch2"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}
	})
}

func TestStacksManager_AddBranch(t *testing.T) {
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

		result := stacksManager.AddBranch("", 0)

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

		result := stacksManager.AddBranch("", 0)

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

		result := stacksManager.AddBranch("my_branch", 0)

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

		result := stacksManager.AddBranch("non_existing_branch", 0)
		want := "Branch " + color.Yellow("non_existing_branch") + " does not exist"
		got := stacksManager.printerMessage()
		if !strings.Contains(got, want) {
			t.Errorf("got \"%s\", want \"%s\"", got, want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}
	})

	t.Run("when passing position 1 to be the first branch", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.AddBranch("my_branch", 1)

		want := "Branch " + color.Yellow("my_branch") + " added to " + color.Green("stack1")
		got := stacksManager.printerMessage()
		if !strings.Contains(got, want) {
			t.Errorf("got \"%s\", want \"%s\"", got, want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}

		data := *stacksManager.stacks
		if data.Stacks[0].Branches[0] != "my_branch" {
			t.Errorf("got %s, want %s", data.Stacks[0].Branches[0], "my_branch")
		}
	})

	t.Run("when passing position is greater then len of branches", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		result := stacksManager.AddBranch("my_branch", 100)

		want := "Branch " + color.Yellow("my_branch") + " added to " + color.Green("stack1")
		got := stacksManager.printerMessage()
		if !strings.Contains(got, want) {
			t.Errorf("got \"%s\", want \"%s\"", got, want)
		}

		if result != nil {
			t.Errorf("show have no error, got %s", result)
		}

		data := *stacksManager.stacks
		if data.Stacks[0].Branches[len(data.Stacks[0].Branches)-1] != "my_branch" {
			t.Errorf("got %s, want %s", data.Stacks[0].Branches[0], "my_branch")
		}
	})
}

func TestStacksManager_List(t *testing.T) {
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

func TestStacksManager_ListStacksForCompletion(t *testing.T) {
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

func TestStacksManager_SwitchByName(t *testing.T) {
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

func TestStacksManager_SwitchByNumber(t *testing.T) {
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

func TestStacksManager_RemoveByName(t *testing.T) {
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

		if result != nil {
			t.Errorf("should have no error, got %s", result)
		}

		want := "Branch " + color.Yellow("non_existing_branch") + " does not exist"
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})
}

func TestStacksManager_StacksManager_RemoveByNumber(t *testing.T) {
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
		if result != nil {
			t.Errorf("got none, want Error")
		}
		want := "Invalid branch number"
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})
}

func TestStacksManager_StacksManager_Delete(t *testing.T) {
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

		if gitCommandsReceived[0] != "checkout branch1" {
			t.Errorf("got %s, want %s", gitCommandsReceived[1], "checkout branch1")
		}

		if err != nil {
			t.Errorf("show have no error, got %s", err)
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

		if gitCommandsReceived[0] != "checkout branch1" {
			t.Errorf("got %s, want %s", gitCommandsReceived[1], "checkout branch1")
		}

		if err == nil {
			t.Errorf("got none, want Error")
		}
	})
}

func TestStacksManager_CheckoutByNumber(t *testing.T) {
	t.Run("checkout branch by number", func(t *testing.T) {
		var gitCommandsReceived []string
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch strings.Join(command, " ") {
				case "checkout branch1":
					gitCommandsReceived = append(gitCommandsReceived, strings.Join(command, " "))
					return "Switched to branch branch1", nil
				default:
					t.Errorf("unwanted git command: %s", command)
				}
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.CheckoutByNumber(1)

		if gitCommandsReceived[0] != "checkout branch1" {
			t.Errorf("got %s, want %s", gitCommandsReceived[0], "checkout branch1")
		}

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}
	})

	t.Run("when number is invalid", func(t *testing.T) {
		var messageReceived []string
		stacksManager := StacksManagerForTest(nil, &messageReceived)

		err := stacksManager.CheckoutByNumber(3)

		if err == nil {
			t.Errorf("got none, want Error")
		}
	})

	t.Run("when checkout return an error", func(t *testing.T) {
		var gitCommandsReceived []string
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				switch strings.Join(command, " ") {
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

		err := stacksManager.CheckoutByNumber(1)

		if gitCommandsReceived[0] != "checkout branch1" {
			t.Errorf("got %s, want %s", gitCommandsReceived[0], "checkout branch1")
		}

		if err == nil {
			t.Errorf("got none, want Error")
		}
	})
}

func TestStacksManager_Sync(t *testing.T) {
	t.Run("when unstaged changes", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "file1.txt", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, false)

		want := "Unstaged changes. Please commit or stash them"
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}
	})

	t.Run("when get currentBranchName return error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				if "rev-parse --abbrev-ref HEAD" == joinedCommand {
					return "", fmt.Errorf("git command error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, false)

		if err == nil {
			t.Errorf("got none, want Error")
		}

		want := "failed to get current branch"
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
		}
	})

	t.Run("when fetch return an error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if "fetch" == command[0] {
					return "", fmt.Errorf("git command error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, false)

		if err == nil {
			t.Errorf("got none, want Error")
		}
		want := "failed to fetch"
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
		}
	})

	t.Run("sync", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, false)

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}

		want := fmt.Sprintf(
			`Syncing %s
Fetching...
Branch: %s
	Checkout...
	Pull...
Branch: %s
	Checkout...
	Pull...
	Merging %s`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Yellow("branch2"),
			color.Yellow("branch1"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})

	t.Run("when checkout return an error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				if "checkout branch1" == joinedCommand {
					return "", fmt.Errorf("checkout error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, false)

		if err == nil {
			t.Errorf("got none, want Error")
		}
		want := "failed to checkout " + color.Yellow("branch1")
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
		}
	})

	t.Run("when pull return an error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				if "pull" == joinedCommand {
					return "", fmt.Errorf("pull error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, false)

		if err == nil {
			t.Errorf("got none, want Error")
		}
		want := "failed to pull"
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
		}
	})

	t.Run("when pull return `There is no tracking information`", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				joinedCommand := strings.Join(command, " ")
				if "pull" == joinedCommand {
					return "There is no tracking information", fmt.Errorf("pull error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, false)

		if err != nil {
			t.Errorf("should have no error, got %s", err)
		}

		want := fmt.Sprintf("Syncing %s", color.Green("stack1"))

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})

	t.Run("when merge return an error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				//joinedCommand := strings.Join(command, " ")
				if command[0] == "merge" {
					return "", fmt.Errorf("merge error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, false)

		if err == nil {
			t.Errorf("got none, want Error")
		}

		want := "failed to merge"
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
		}
	})

	t.Run("when sync with push flag", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(true, false)

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}

		want := fmt.Sprintf(
			`Syncing %s
Fetching...
Branch: %s
	Checkout...
	Pull...
	Pushing...
Branch: %s
	Checkout...
	Pull...
	Merging %s
	Pushing...
`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Yellow("branch2"),
			color.Yellow("branch1"),
		)

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}

	})

	t.Run("when push return `has no upstream branch`", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if command[0] == "push" {
					return "has no upstream branch", fmt.Errorf("push error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(true, false)

		if err != nil {
			t.Errorf("should have no error, got %s", err)
		}

		want := fmt.Sprintf("Syncing %s", color.Green("stack1"))

		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})

	t.Run("when push return error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if command[0] == "push" {
					return "", fmt.Errorf("push error")
				}
				return "", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(true, false)

		if err == nil {
			t.Errorf("should have error, got %s", err)
		}
		want := "failed to push"
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
		}
	})

	t.Run("when sync with merge default branch", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if command[0] == "symbolic-ref" {
					return "origin/main", nil
				}
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, true)

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}

		want := fmt.Sprintf(
			`Syncing %s
Fetching...
Branch: %s
	Checkout...
	Pull...
	Merging %s
Branch: %s
	Checkout...
	Pull...
	Merging %s
`,
			color.Green("stack1"),
			color.Yellow("branch1"),
			color.Yellow("origin/main"),
			color.Yellow("branch2"),
			color.Yellow("branch1"),
		)
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})

	t.Run("when defaultBranchWithRemote return an error", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if command[0] == "symbolic-ref" {
					return "", fmt.Errorf("symbolic-ref error")
				}
				return "", nil
			},
		}

		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Sync(false, true)

		if err == nil {
			t.Errorf("should have error, got %s", err)
		}
		want := "To set it try:"
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
		}
	})
}

func TestStacksManager_Tree(t *testing.T) {
	t.Run("tree", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "abcdef Some commit message - 3 minutes ago", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Tree()

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}

		want := fmt.Sprintf(
			`Current stack: %s 

%s
%s Some commit message - 3 minutes ago
%s%s
%s Some commit message - 3 minutes ago

`,
			color.Green("stack1"),
			color.Red("* branch1"),
			color.Red("| ")+color.DarkYellow("abcdef"),
			color.Red("|\\\n"),
			color.Red("| ")+color.Purple("* branch2"),
			color.Red("| ")+color.Purple("| ")+color.DarkYellow("abcdef"),
		)
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got\n\"%s\"\nwant\n\"%s\"", stacksManager.printerMessage(), want)
		}
	})
}

func TestStacksManager_Publish(t *testing.T) {
	t.Run("publish current branch and first branch of the stack", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if command[0] == "remote" {
					return "git@github.com:User/AwesomeRepo.git", nil
				}
				return "branch1", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Publish()

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}

		want := "Publishing " +
			color.Yellow("branch1") +
			"..." +
			"\nhttps://github.com/User/AwesomeRepo/compare/branch1?expand=1"
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})

	t.Run("when the current branch not part of the current stack", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				return "invalid_branch", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Publish()

		if err == nil {
			t.Errorf("got none, want Error")
		}

		want := "current branch " + color.Yellow("invalid_branch") + " is not part of the current stack " + color.Green("stack1")
		if !strings.Contains(err.Error(), want) {
			t.Errorf("got \"%s\", want \"%s\"", err.Error(), want)
		}
	})

	t.Run("publish current branch and second branch of the stack", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if command[0] == "remote" {
					return "git@github.com:User/AwesomeRepo.git", nil
				}
				return "branch2", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Publish()

		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}

		want := "Publishing " +
			color.Yellow("branch2") +
			"..." +
			"\nhttps://github.com/User/AwesomeRepo/compare/branch1...branch2?expand=1"
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})

	t.Run("when the remote is not on GitHub", func(t *testing.T) {
		gitExecutor := gitExecutorStub{
			stubExec: func(command ...string) (string, error) {
				if command[0] == "remote" {
					return "git@gitlab.com:User/AwesomeRepo.git", nil
				}
				return "branch1", nil
			},
		}
		var messageReceived []string
		stacksManager := StacksManagerForTest(gitExecutor, &messageReceived)

		err := stacksManager.Publish()
		if err != nil {
			t.Errorf("show have no error, got %s", err)
		}

		want := "Publishing " +
			color.Yellow("branch1") + "...\nRemote is not on GitHub. Sorry."
		if !strings.Contains(stacksManager.printerMessage(), want) {
			t.Errorf("got \"%s\", want \"%s\"", stacksManager.printerMessage(), want)
		}
	})
}
