package stack

import (
    "fmt"
    "github.com/Bhacaz/gostacking/internal/git"
)

func (sm StacksManager) New(stackName string) {
    newStack := Stack{
                    Name: stackName,
                    Branches: []string{git.CurrentBranchName()},
                }

    data, _ := sm.stacksPersister.LoadStacks()
    data.CurrentStack = stackName
    data.Stacks = append(data.Stacks, newStack)

    sm.stacksPersister.SaveStacks(data)
	fmt.Println("New stack created", stackName)
}

//
// func CurrentStackName() string {
//     data, err := LoadStacks()
//     if err != nil {
//         fmt.Println("Error loading JSON:", err)
//         return ""
//     }
//
//     return data.CurrentStack
// }
//
// func CurrentStackStatus() string {
//     data, err := LoadStacks()
//     if err != nil {
//         fmt.Println("Error loading JSON:", err)
//         return ""
//     }
//
//     var displayBranches string
//     branches, _ := data.GetBranchesByName(data.CurrentStack)
//     for i, branch := range branches {
//         // Maybe someday it will be nice to add
//         // git log --pretty=format:'%s - %Cred%h%Creset %C(bold blue)%an%Creset %Cgreen%cr%Creset' -n 1 master
//         displayBranches += fmt.Sprintf("%d. " + color.Yellow(branch) + "\n", i + 1)
//     }
//     return "Current stack: " + color.Green(data.CurrentStack) + "\nBranches:\n" + displayBranches
// }
//
// func Add(branchName string) {
//      if branchName == "" {
//         branchName = git.CurrentBranchName()
//     } else {
//         if !git.BranchExists(branchName) {
//             fmt.Println("Branch", branchName, "does not exist")
//         }
//     }
//
//     data, _ := LoadStacks()
//     stack, _ := data.GetStackByName(data.CurrentStack)
//     stack.Branches = append(stack.Branches, branchName)
//     stack.Branches = slices.Compact(stack.Branches)
//     data.SaveStacks()
//     fmt.Println("Branch", branchName, "added to stack", data.CurrentStack)
// }
//
// func List() {
//     data, _ := LoadStacks()
//     fmt.Println("Current stack:", color.Green(data.CurrentStack))
//     for i, stack := range data.Stacks {
//         fmt.Printf("%d. %s\n", i + 1, color.Yellow(stack.Name))
//     }
// }
//
// func SwitchByName(stackName string) {
//     data, _ := LoadStacks()
//     data.CurrentStack = stackName
//     data.SaveStacks()
//     fmt.Println("Switched to stack", stackName)
// }
//
// func SwitchByNumber(number int) {
//     data, _ := LoadStacks()
//     stack := data.Stacks[number - 1]
//     data.CurrentStack = stack.Name
//     data.SaveStacks()
//     fmt.Println("Switched to stack", stack.Name)
// }
//
// func RemoveByName(branchName string) {
//     data, _ := LoadStacks()
//     stack, _ := data.GetStackByName(data.CurrentStack)
//     var filteredBranches []string
//     for _, branch := range stack.Branches {
//         if branch != branchName {
//             filteredBranches = append(filteredBranches, branch)
//         }
//     }
//
//     if len(filteredBranches) == len(stack.Branches) {
//         fmt.Println("Branch", branchName, "does not exist")
//         return
//     }
//
//     stack.Branches = filteredBranches
//     data.SaveStacks()
//     fmt.Println("Branch", color.Yellow(branchName), "removed from stack", color.Green(data.CurrentStack))
// }
//
// func RemoveByNumber(number int) {
//     data, _ := LoadStacks()
//     stack, _ := data.GetStackByName(data.CurrentStack)
//     if number < 1 || number > len(stack.Branches) {
//         fmt.Println("Invalid branch number")
//         return
//     }
//
//     branchName := stack.Branches[number - 1]
//     stack.Branches = append(stack.Branches[:number - 1], stack.Branches[number:]...)
//     data.SaveStacks()
//     fmt.Println("Branch", color.Yellow(branchName), "removed from stack", color.Green(data.CurrentStack))
// }
//
// func Delete(stackName string) {
//     data, _ := LoadStacks()
//     var filteredStacks []Stack
//     for _, stack := range data.Stacks {
//         if stack.Name != stackName {
//             filteredStacks = append(filteredStacks, stack)
//         }
//     }
//
//     if len(filteredStacks) == len(data.Stacks) {
//         fmt.Println("Stack", stackName, "does not exist")
//         return
//     }
//
//     data.Stacks = filteredStacks
//
//     if data.CurrentStack == stackName {
//         data.CurrentStack = data.Stacks[0].Name
//     }
//
//     data.SaveStacks()
//     fmt.Println("Stack", stackName, "deleted from stack")
//     fmt.Println("Current stack status:")
//     fmt.Println(CurrentStackStatus())
// }
//
// func Sync(push bool) {
//     data, _ := LoadStacks()
//     currentBranch := git.CurrentBranchName()
//     branches, _ := data.GetBranchesByName(data.CurrentStack)
//     git.SyncBranches(branches, currentBranch, push)
// }
//
// func CheckoutByName(branchName string) {
//     if !git.BranchExists(branchName) {
//         fmt.Println("Branch", branchName, "does not exist")
//         return
//     }
//
//     git.Checkout(branchName)
// }
//
// func CheckoutByNumber(number int) {
//     data, _ := LoadStacks()
//     branches, _ := data.GetBranchesByName(data.CurrentStack)
//     if number < 1 || number > len(branches) {
//         fmt.Println("Invalid branch number")
//         return
//     }
//
//     git.Checkout(branches[number - 1])
// }
