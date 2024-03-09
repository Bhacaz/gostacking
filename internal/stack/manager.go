package stack

import (
	"errors"
	"fmt"
	"github.com/Bhacaz/gostacking/internal/color"
	"github.com/Bhacaz/gostacking/internal/git"
	"github.com/Bhacaz/gostacking/internal/printer"
	"strings"
)

type StacksManager struct {
	stacks      *StacksData
	gitExecutor git.InterfaceGitExecutor
	printer     printer.Printer
}

func NewManager(gitVerbose bool) StacksManager {
	return StacksManager{
		stacks: &StacksData{
			StacksPersister: StacksPersistingFile{},
		},
		printer:     printer.NewPrinter(),
		gitExecutor: git.NewExecutor(gitVerbose),
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

	sm.stacks.LoadStacks()
	sm.stacks.CurrentStack = stackName
	sm.stacks.Stacks = append(sm.stacks.Stacks, newStack)
	sm.stacks.SaveStacks()
	sm.printer.Println("Stack created", color.Green(stackName))
	return nil
}

// CurrentStackStatus Will show start for:
// 1. Behind remote
// 2. Has diff with previous branch
func (sm StacksManager) CurrentStackStatus(showLog bool) error {
	sm.stacks.LoadStacks()
	data := *sm.stacks
	err := sm.fetch()
	if err != nil {
		return err
	}

	var displayBranches string
	branches, _ := data.GetBranchesByName(data.CurrentStack)
	for i, branch := range branches {
		branchStatus := defaultBranchStatus()

		if sm.isBehindRemote(branch) {
			branchStatus.BehindRemote = true
		}
		if sm.aheadRemote(branch) {
			branchStatus.AheadRemote = true
		}

		displayBranches += fmt.Sprintf("%d. "+color.Yellow(branch), i+1)
		if i == 0 {
			branchStatus.BehindDefaultBranch = sm.behindDefaultBranch(branch)
		} else {
			hasDiff, _ := sm.branchHasDiff(branches[i-1], branch)
			if hasDiff {
				branchStatus.HasDiff = true
			}
		}

		displayBranches += branchStatus.Symbols()

		if showLog {
			displayBranches += "\n\t" + sm.lastLog(branch)
		}

		displayBranches += "\n"
	}
	sm.printer.Println("Current stack: " + color.Green(data.CurrentStack) + "\nBranches:\n" + displayBranches)
	return nil
}

func (sm StacksManager) AddBranch(branchName string, position int) error {
	if branchName == "" {
		currentBranchName, err := sm.currentBranchName()
		if err != nil {
			return err
		}
		branchName = currentBranchName
	} else {
		if !sm.branchExists(branchName) {
			sm.printer.Println("Branch " + color.Yellow(branchName) + " does not exist")
			return nil
		}
	}

	sm.stacks.LoadStacks()
	data := *sm.stacks
	stack, _ := data.GetStackByName(data.CurrentStack)

	if position == 0 || position > len(stack.Branches) {
		stack.Branches = append(stack.Branches, branchName)
	} else {
		position--
		newBranches := make([]string, 0, len(stack.Branches)+1)
		for i, branch := range stack.Branches {
			if i == position {
				newBranches = append(newBranches, branchName)
			}
			newBranches = append(newBranches, branch)
		}
		stack.Branches = newBranches
	}

	data.SaveStacks()
	sm.printer.Println("Branch", color.Yellow(branchName), "added to", color.Green(data.CurrentStack))
	return nil
}

func (sm StacksManager) List() error {
	sm.stacks.LoadStacks()
	data := *sm.stacks
	sm.printer.Println("Current stack:", color.Green(data.CurrentStack))
	for i, stack := range data.Stacks {
		sm.printer.Println(
			fmt.Sprintf("%d. %s", i+1, color.Yellow(stack.Name)),
		)
	}
	return nil
}

func (sm StacksManager) ListStacksForCompletion(toComplete string) []string {
	sm.stacks.LoadStacks()
	data := *sm.stacks
	var stacks []string
	for _, stack := range data.Stacks {
		if toComplete == "" || strings.HasPrefix(stack.Name, toComplete) {
			stacks = append(stacks, stack.Name)
		}
	}
	return stacks
}

func (sm StacksManager) SwitchByName(stackName string) error {
	sm.stacks.LoadStacks()
	var stack *Stack
	var err error
	if stackName == "" {
		currentBranchName, err := sm.currentBranchName()
		if err != nil {
			return err
		}
		stack, err = sm.stacks.GetStackByBranch(currentBranchName)
		if err != nil {
			return err
		}
	} else {
		stack, err = sm.stacks.GetStackByName(stackName)
		fmt.Println(stack, err)
		if err != nil {
			return err
		}
	}
	sm.stacks.SetCurrentStack(stack.Name)
	sm.printer.Println("Switched to stack", color.Green(stack.Name))
	return nil
}

func (sm StacksManager) SwitchByNumber(number int) error {
	sm.stacks.LoadStacks()

	if number < 1 || number > len(sm.stacks.Stacks) {
		return errors.New("invalid stack number")
	}

	stack := sm.stacks.Stacks[number-1]
	sm.stacks.SetCurrentStack(stack.Name)
	sm.stacks.SaveStacks()
	sm.printer.Println("Switched to stack", color.Green(stack.Name))
	return nil
}

func (sm StacksManager) RemoveByName(branchName string) error {
	sm.stacks.LoadStacks()
	data := *sm.stacks
	stack, _ := data.GetStackByName(data.CurrentStack)
	var filteredBranches []string
	for _, branch := range stack.Branches {
		if branch != branchName {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	if len(filteredBranches) == len(stack.Branches) {
		sm.printer.Println("Branch", color.Yellow(branchName), "does not exist")
		return nil
	}

	stack.Branches = filteredBranches
	data.SaveStacks()
	sm.printer.Println("Branch", color.Yellow(branchName), "removed from", color.Green(data.CurrentStack))
	return nil
}

func (sm StacksManager) RemoveByNumber(number int) error {
	sm.stacks.LoadStacks()
	data := *sm.stacks
	stack, _ := data.GetStackByName(data.CurrentStack)
	if number < 1 || number > len(stack.Branches) {
		sm.printer.Println("Invalid branch number")
		return nil
	}

	branchName := stack.Branches[number-1]
	stack.Branches = append(stack.Branches[:number-1], stack.Branches[number:]...)
	data.SaveStacks()
	sm.printer.Println("Branch", color.Yellow(branchName), "removed from stack", color.Green(data.CurrentStack))
	return nil
}

func (sm StacksManager) Delete(stackName string) error {
	sm.stacks.LoadStacks()
	var filteredStacks []Stack
	for _, stack := range sm.stacks.Stacks {
		if stack.Name != stackName {
			filteredStacks = append(filteredStacks, stack)
		}
	}

	if len(filteredStacks) == len(sm.stacks.Stacks) {
		return errors.New("stack " + color.Green(stackName) + " does not exist")
	}

	sm.stacks.Stacks = filteredStacks

	var newCurrentStack = sm.stacks.CurrentStack
	if len(filteredStacks) > 0 && sm.stacks.CurrentStack == stackName {
		newCurrentStack = sm.stacks.Stacks[0].Name
	} else if len(filteredStacks) == 0 {
		newCurrentStack = ""
	}

	sm.stacks.SetCurrentStack(newCurrentStack)
	sm.printer.Println("Stack", color.Green(stackName), "deleted")
	return nil
}

func (sm StacksManager) CheckoutByName(branchName string) error {
	return sm.checkout(branchName)
}

func (sm StacksManager) CheckoutByNumber(number int) error {
	sm.stacks.LoadStacks()
	data := *sm.stacks
	branches, _ := data.GetBranchesByName(data.CurrentStack)
	if number < 1 || number > len(branches) {
		return errors.New("invalid branch number")
	}

	return sm.checkout(branches[number-1])
}

func (sm StacksManager) Sync(push bool, mergeDefaultBranch bool) error {
	sm.stacks.LoadStacks()
	data := *sm.stacks
	if sm.unstagedChanges() {
		sm.printer.Println("Unstaged changes. Please commit or stash them")
		return nil
	}

	checkoutBranchEnd, err := sm.currentBranchName()
	if err != nil {
		return err
	}

	sm.printer.Println("Syncing", color.Green(data.CurrentStack))

	sm.printer.Println("Fetching...")
	err = sm.fetch()
	if err != nil {
		return err
	}

	branches, _ := data.GetBranchesByName(data.CurrentStack)

	for i, branch := range branches {
		sm.printer.Println("Branch:", color.Yellow(branch))
		sm.printer.Println("\tCheckout...")
		err = sm.checkout(branch)
		if err != nil {
			return err
		}

		sm.printer.Println("\tPull...")
		err = sm.pullBranch()
		if err != nil {
			return err
		}

		if i == 0 {
			err = sm.syncFirstBranch(branch, push, mergeDefaultBranch)
			if err != nil {
				return err
			}
			continue
		}

		parentBranch := branches[i-1]
		sm.printer.Println("\tMerging", color.Yellow(parentBranch))
		err = sm.merge(branch, parentBranch)
		if err != nil {
			return err
		}
		if push {
			sm.printer.Println("\tPushing...")
			err = sm.pushBranch()
			if err != nil {
				return err
			}
		}
	}

	return sm.checkout(checkoutBranchEnd)
}

func (sm StacksManager) Tree() error {
	//branchColorsSequence := []color.Color{
	//	color.Yellow,
	//	color.Green,
	//	color.Magenta,
	//	color.Purple,
	//}

	sm.stacks.LoadStacks()
	branches, _ := sm.stacks.GetCurrentBranches()
	sm.printer.Println("Current stack:", color.Green(sm.stacks.CurrentStack))
	lastIndex := len(branches) - 1
	treeOutput := ""
	// prepend default branch to the list of branches
	defaultBranch, err := sm.defaultBranchWithRemote()
	if err != nil {
		return err
	}
	branches = append([]string{defaultBranch}, branches...)

	for i, branch := range branches {
		if i == lastIndex {
			treeOutput += branch + "\n"
			continue
		}

		treeOutput += branches[i+1] + "\n"
		commits, err := sm.commitsBetweenBranches(branch, branches[i+1])
		if err != nil {
			return err
		}

		for _, commit := range commits {
			treeOutput += "|\t" + commit + "\n"
		}
	}
	return nil
}

func (sm StacksManager) syncFirstBranch(firstBranch string, push bool, mergeDefaultBranch bool) error {
	if mergeDefaultBranch {
		defaultBranch, err := sm.defaultBranchWithRemote()
		if err != nil {
			return err
		}
		sm.printer.Println("\tMerging", color.Yellow(defaultBranch))
		err = sm.merge(firstBranch, defaultBranch)
		if err != nil {
			return err
		}
	}

	if push {
		sm.printer.Println("\tPushing...")
		return sm.pushBranch()
	}
	return nil
}
