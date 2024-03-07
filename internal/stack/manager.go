package stack

import (
	"errors"
	"fmt"
	"github.com/Bhacaz/gostacking/internal/color"
	"github.com/Bhacaz/gostacking/internal/git"
	"github.com/Bhacaz/gostacking/internal/printer"
	"log"
	"slices"
	"strings"
)

type StacksManager struct {
	stacksPersister StacksPersisting
	gitExecutor     git.InterfaceGitExecutor
	printer         printer.Printer
}

func (sm StacksManager) load() StacksData {
	data, err := sm.stacksPersister.LoadStacks()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return data
}

func NewManager() StacksManager {
	return StacksManager{
		stacksPersister: StacksPersistingFile{},
		printer:         printer.NewPrinter(),
		gitExecutor:     git.NewExecutor(true),
	}
}

func (sm StacksManager) CreateStack(stackName string) error {
	currentBranch, err := sm.currentBranchName()
	if err != nil {
		return err
	}

	newStack := Stack{
		Name:     stackName,
		Branches: []string{currentBranch},
	}

	data := sm.load()
	data.CurrentStack = stackName
	data.Stacks = append(data.Stacks, newStack)

	sm.stacksPersister.SaveStacks(data)
	sm.printer.Println("Stack created", color.Green(stackName))
	return nil
}

// CurrentStackStatus Will show start for:
// 1. Behind remote
// 2. Has diff with previous branch
func (sm StacksManager) CurrentStackStatus(showLog bool) error {
	data := sm.load()
	err := sm.fetch()
	if err != nil {
		return err
	}

	var displayBranches string
	branches, _ := data.GetBranchesByName(data.CurrentStack)
	for i, branch := range branches {
		var showStar = false
		if sm.isBehindRemote(branch) {
			showStar = true
		}

		displayBranches += fmt.Sprintf("%d. "+color.Yellow(branch), i+1)
		if i > 0 {
			// Don't check diff if star is already shown
			if !showStar {
				hasDiff, err := sm.branchHasDiff(branches[i-1], branch)
				if err != nil {
					displayBranches += fmt.Sprintf("Could not get diff status for %s...%s\n%s", branches[i-1], branch, err.Error())
				}
				if hasDiff {
					showStar = true
				}
			}
		}

		if showStar {
			displayBranches += " " + color.Red("*")
		}
		if showLog {
			displayBranches += "\n\t" + sm.lastLog(branch)
		}

		displayBranches += "\n"
	}
	sm.printer.Println("Current stack: " + color.Green(data.CurrentStack) + "\nBranches:\n" + displayBranches)
	return nil
}

func (sm StacksManager) AddBranch(branchName string) error {
	if branchName == "" {
		currentBranchName, err := sm.currentBranchName()
		if err != nil {
			return err
		}
		branchName = currentBranchName
	} else {
		if !sm.branchExists(branchName) {
			return errors.New("Branch " + color.Yellow(branchName) + " does not exist")
		}
	}

	data := sm.load()
	stack, _ := data.GetStackByName(data.CurrentStack)
	stack.Branches = append(stack.Branches, branchName)
	stack.Branches = slices.Compact(stack.Branches)
	sm.stacksPersister.SaveStacks(data)
	sm.printer.Println("Branch", color.Yellow(branchName), "added to stack", color.Green(data.CurrentStack))
	return nil
}

func (sm StacksManager) List() error {
	data := sm.load()
	sm.printer.Println("Current stack:", color.Green(data.CurrentStack))
	for i, stack := range data.Stacks {
		sm.printer.Println(
			fmt.Sprintf("%d. %s", i+1, color.Yellow(stack.Name)),
		)
	}
	return nil
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

func (sm StacksManager) SwitchByName(stackName string) error {
	data := sm.load()
	var stack *Stack
	var err error
	if stackName == "" {
		currentBranchName, err := sm.currentBranchName()
		if err != nil {
			return err
		}
		stack, err = data.GetStackByBranch(currentBranchName)
		if err != nil {
			return err
		}
	} else {
		stack, err = data.GetStackByName(stackName)
		if err != nil {
			return err
		}
	}
	data.SetCurrentStack(stack.Name)
	sm.printer.Println("Switched to stack", color.Green(stack.Name))
	return nil
}

func (sm StacksManager) SwitchByNumber(number int) error {
	data := sm.load()
	stack := data.Stacks[number-1]
	data.CurrentStack = stack.Name
	sm.stacksPersister.SaveStacks(data)
	sm.printer.Println("Switched to stack", color.Green(stack.Name))
	return nil
}

func (sm StacksManager) RemoveByName(branchName string) error {
	data := sm.load()
	stack, _ := data.GetStackByName(data.CurrentStack)
	var filteredBranches []string
	for _, branch := range stack.Branches {
		if branch != branchName {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	if len(filteredBranches) == len(stack.Branches) {
		return errors.New("Branch " + branchName + " does not exist")
	}

	stack.Branches = filteredBranches
	sm.stacksPersister.SaveStacks(data)
	sm.printer.Println("Branch", color.Yellow(branchName), "removed from stack", color.Green(data.CurrentStack))
	return nil
}

func (sm StacksManager) RemoveByNumber(number int) error {
	data := sm.load()
	stack, _ := data.GetStackByName(data.CurrentStack)
	if number < 1 || number > len(stack.Branches) {
		return errors.New("invalid branch number")
	}

	branchName := stack.Branches[number-1]
	stack.Branches = append(stack.Branches[:number-1], stack.Branches[number:]...)
	sm.stacksPersister.SaveStacks(data)
	sm.printer.Println("Branch", color.Yellow(branchName), "removed from stack", color.Green(data.CurrentStack))
	return nil
}

func (sm StacksManager) Delete(stackName string) error {
	data := sm.load()
	var filteredStacks []Stack
	for _, stack := range data.Stacks {
		if stack.Name != stackName {
			filteredStacks = append(filteredStacks, stack)
		}
	}

	if len(filteredStacks) == len(data.Stacks) {
		return errors.New("stack " + color.Green(stackName) + " does not exist")
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
	sm.printer.Println("Stack", color.Green(stackName), "deleted")
	return nil
}

func (sm StacksManager) Sync(push bool, mergeHead bool) error {
	data := sm.load()
	_, err := sm.currentBranchName()
	if err != nil {
		return err
	}

	sm.printer.Println("Syncing", color.Green(data.CurrentStack))
	//branches, _ := data.GetBranchesByName(data.CurrentStack)
	//sm.gitCommands.SyncBranches(branches, currentBranch, push, mergeHead)
	return errors.New("not implemented")
}

func (sm StacksManager) CheckoutByName(branchName string) error {
	if !sm.branchExists(branchName) {
		return errors.New("branch does not exist")
	}

	return sm.checkout(branchName)
}

func (sm StacksManager) CheckoutByNumber(number int) error {
	data := sm.load()
	branches, _ := data.GetBranchesByName(data.CurrentStack)
	if number < 1 || number > len(branches) {
		return errors.New("invalid branch number")
	}

	return sm.checkout(branches[number-1])
}
