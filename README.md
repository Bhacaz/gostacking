# gostacking

Allow to simply links branches together to create a stack of branches.

This tool use **merges** instead of rebases like other tools to be less "destructive" and maybe more easy to understand
what is going on. Down side are more commits and no way to change order with the tool.

With the command `sync`,
allow to update all branches to merge it one into the other. `sync` will do a `git pull` first, then
a `git merge` branch 1 into branch 2, branch 2 into branch 3, etc.

## Installation

⚠️ Binary are not signed, you will need to allow it to run on your system. The tool is still in development.

Go to [releases](https://github.com/Bhacaz/gostacking/releases/latest) and download the binary for your OS.

Optionally, add the binary to your path.

## Commands

Usage:
`gostacking [command]`

Available Commands:
* `add`         Add a branch to the current stack. If no branch is given, add the current branch.
* `delete`      Delete a gostacking.
* `help`        Help about any command
* `list`        List all stacks.
* `new`         Create a new gostacking.
* `status`      Get current stack.
* `switch`      Change the current stack. Using name or index.
* `sync`        Merge all branch in a stack into the current branch.

## Example

```bash
gostacking new my-stack
gostacking add status
# my-stack
#  1. main
gostacking add feature/1
gostacking add feature/2
gostacking status
# my-stack
#  1. main
#  2. feature/1
#  3. feature/2
git checkout feature/1
touch file1.txt
git add file1.txt && git commit -m "Add file1.txt"
gostacking sync
# Merge main into feature/1
# Merge feature/1 into feature/2
git log --oneline -n 1 feature/2
# gostacking - Merge feature/1 into feature/2
```

## Notes

https://cobra.dev/

```bash
~/go/bin/cobra-cli add new
go mod tidy
go build -ldflags="-s -w" -o bin
```

## TODOs

- [x] Add new stack
- [x] Status (current stack and list branches)
- [x] Add current branch to stack
- [x] Add specific branch to stack
- [x] List, show all stacks
- [x] Switch, switch to stack
- [X] Remove, remove stack
- [x] Sync all branch in a stack (merge each one into the other)
- [ ] Change way prints works, stack.go should return string and cmd should print it
- [ ] Add colors
- Update readme