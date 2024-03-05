package stack

import (
	"fmt"
	"github.com/Bhacaz/gostacking/internal/color"
	"github.com/Bhacaz/gostacking/internal/git"
	"log"
	"slices"
	"strings"
)

type StacksManager struct {
	stacksPersister StacksPersisting
	gitCommands     git.InterfaceCommands
}

func (sm StacksManager) load() StacksData {
	data, err := sm.stacksPersister.LoadStacks()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return data
}

func Manager() StacksManager {
	return StacksManager{
		stacksPersister: StacksPersistingFile{},
		gitCommands:     git.Cmd(),
	}
}

func (sm StacksManager) CreateStack(stackName string) string {
	currentBranch, err := sm.gitCommands.CurrentBranchName()
	if err != nil {
		log.Fatalf(err.Error())
	}

	newStack := Stack{
		Name:     stackName,
		Branches: []string{currentBranch},
	}

	data := sm.load()
	data.CurrentStack = stackName
	data.Stacks = append(data.Stacks, newStack)

	sm.stacksPersister.SaveStacks(data)
	return "Stack created " + color.Green(stackName)
}

// CurrentStackStatus Will show start for:
// 1. Behind remote
// 2. Has diff with previous branch
func (sm StacksManager) CurrentStackStatus(showLog bool) string {
	data := sm.load()
	sm.gitCommands.Fetch()

	var displayBranches string
	branches, _ := data.GetBranchesByName(data.CurrentStack)
	for i, branch := range branches {
		var showStar = false
		if sm.gitCommands.IsBehindRemote(branch) {
			showStar = true
		}

		displayBranches += fmt.Sprintf("%d. "+color.Yellow(branch), i+1)
		if i > 0 {
			hasDiff := sm.gitCommands.BranchDiff(branches[i-1], branch)
			if hasDiff && !showStar {
				showStar = true
			}
		}

		if showStar {
			displayBranches += " " + color.Red("*")
		}
		if showLog {
			displayBranches += "\n\t" + sm.gitCommands.LastLog(branch)
		}

		displayBranches += "\n"
	}
	return "Current stack: " + color.Green(data.CurrentStack) + "\nBranches:\n" + displayBranches
}

func (sm StacksManager) AddBranch(branchName string) {
	if branchName == "" {
		branchName, _ = sm.gitCommands.CurrentBranchName()
	} else {
		if !sm.gitCommands.BranchExists(branchName) {
			log.Fatalf("Branch %s does not exist", branchName)
		}
	}

	data := sm.load()
	stack, _ := data.GetStackByName(data.CurrentStack)
	stack.Branches = append(stack.Branches, branchName)
	stack.Branches = slices.Compact(stack.Branches)
	sm.stacksPersister.SaveStacks(data)
	fmt.Println("Branch", color.Yellow(branchName), "added to stack", color.Green(data.CurrentStack))
}

func (sm StacksManager) List() {
	data := sm.load()
	fmt.Println("Current stack:", color.Green(data.CurrentStack))
	for i, stack := range data.Stacks {
		fmt.Printf("%d. %s\n", i+1, color.Yellow(stack.Name))
	}
}

func (sm StacksManager) ListStacksForCompletion(toComplete string) []string {
	data := sm.load()
	var stacks []string
	for _, stack := range data.Stacks {
		if toComplete == "" || strings.HasPrefix(stack.Name, toComplete) {
			stacks = append(stacks, stack.Name)
		}
	}
	return stacks
}

func (sm StacksManager) SwitchByName(stackName string) {
	data := sm.load()
	var stack *Stack
	var err error
	if stackName == "" {
		currentBranchName, _ := sm.gitCommands.CurrentBranchName()
		stack, err = data.GetStackByBranch(currentBranchName)
		if err != nil {
			fmt.Println("No stack found for branch", currentBranchName)
			return
		}
	} else {
		stack, err = data.GetStackByName(stackName)
		if err != nil {
			fmt.Println("Stack", stackName, "does not exist")
			return
		}
	}
	data.SetCurrentStack(stack.Name)
	fmt.Println("Switched to stack", color.Green(stack.Name))
}

func (sm StacksManager) SwitchByNumber(number int) {
	data := sm.load()
	stack := data.Stacks[number-1]
	data.CurrentStack = stack.Name
	sm.stacksPersister.SaveStacks(data)
	fmt.Println("Switched to stack", color.Green(stack.Name))
}

func (sm StacksManager) RemoveByName(branchName string) {
	data := sm.load()
	stack, _ := data.GetStackByName(data.CurrentStack)
	var filteredBranches []string
	for _, branch := range stack.Branches {
		if branch != branchName {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	if len(filteredBranches) == len(stack.Branches) {
		fmt.Println("Branch", branchName, "does not exist")
		return
	}

	stack.Branches = filteredBranches
	sm.stacksPersister.SaveStacks(data)
	fmt.Println("Branch", color.Yellow(branchName), "removed from stack", color.Green(data.CurrentStack))
}

func (sm StacksManager) RemoveByNumber(number int) {
	data := sm.load()
	stack, _ := data.GetStackByName(data.CurrentStack)
	if number < 1 || number > len(stack.Branches) {
		fmt.Println("Invalid branch number")
		return
	}

	branchName := stack.Branches[number-1]
	stack.Branches = append(stack.Branches[:number-1], stack.Branches[number:]...)
	sm.stacksPersister.SaveStacks(data)
	fmt.Println("Branch", color.Yellow(branchName), "removed from stack", color.Green(data.CurrentStack))
}

func (sm StacksManager) Delete(stackName string) {
	data := sm.load()
	var filteredStacks []Stack
	for _, stack := range data.Stacks {
		if stack.Name != stackName {
			filteredStacks = append(filteredStacks, stack)
		}
	}

	if len(filteredStacks) == len(data.Stacks) {
		fmt.Println("Stack", stackName, "does not exist")
		return
	}

	data.Stacks = filteredStacks

	var newCurrentStack = data.CurrentStack
	if len(filteredStacks) > 0 && data.CurrentStack == stackName {
		newCurrentStack = data.Stacks[0].Name
	} else if len(filteredStacks) == 0 {
		newCurrentStack = ""
	}
	data.CurrentStack = newCurrentStack

	sm.stacksPersister.SaveStacks(data)
	fmt.Println("Stack", color.Green(stackName), "deleted")
}

func (sm StacksManager) Sync(push bool, mergeHead bool) {
	data := sm.load()
	currentBranch, err := sm.gitCommands.CurrentBranchName()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("Syncing", color.Green(data.CurrentStack))
	branches, _ := data.GetBranchesByName(data.CurrentStack)
	sm.gitCommands.SyncBranches(branches, currentBranch, push, mergeHead)
}

func (sm StacksManager) CheckoutByName(branchName string) {
	if !sm.gitCommands.BranchExists(branchName) {
		fmt.Println("Branch", branchName, "does not exist")
		return
	}

	sm.gitCommands.Checkout(branchName)
}

func (sm StacksManager) CheckoutByNumber(number int) {
	data := sm.load()
	branches, _ := data.GetBranchesByName(data.CurrentStack)
	if number < 1 || number > len(branches) {
		fmt.Println("Invalid branch number")
		return
	}

	sm.gitCommands.Checkout(branches[number-1])
}
