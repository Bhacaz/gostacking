# gostacking

To know more about the [Stacking workflow](https://stacking.dev/).

Allow to simply links branches together to create a stack of branches.

This tool use **merges** instead of rebases like other tools to be less "destructive" and maybe easier to understand
what is going on.

With the command `sync`,
allow to update all branches to merge it one into the other. 

The `sync` command will do:
1. Checkout the branch.
2. Pull the latest changes.
3. Merge the previous branch into the current one.
4. _Optionally_ push the changes (with the `--push` flag).

## Installation

Only **MacOS** is supported for now via Homebrew.

```bash
brew tap Bhacaz/tap
brew install gostacking
brew link gostacking
```

## Commands

Usage:
`gostacking [command]`

```
Available Commands:
  add         Add a branch to the current stack. If no branch is given, add the current branch.
  checkout    Checkout a branch from a stack.
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a gostacking.
  help        Help about any command
  list        List all stacks.
  new         Create a new gostacking.
  remove      Remove a branch from the current stack. (Branch name or number)
  status      Get current stack.
  switch      Change the current stack.
  sync        Merge all branch in a stack into the current branch.
```

## Example

```bash
git checkout -b feature/1
gostacking new my-stack
gostacking add status
# my-stack
#  1. feature/1

gostacking add feature/2
gostacking status
# my-stack
#  1. feature/1
#  2. feature/2

touch file1.txt
git add file1.txt && git commit -m "Add file1.txt"

gostacking status
# my-stack
#  1. feature/1
#  2. feature/2 *

gostacking sync
# Merge feature/1 into feature/2

git log --oneline -n 1 feature/2
# Merge feature/1 into feature/2 (gostacking)
```

## Notes

https://cobra.dev/

```bash
~/go/bin/cobra-cli add new
go mod tidy
go build -ldflags="-s -w" -o dist
goreleaser release
```

## Release

1. Update `.version`
2. `zsh scripts/release.zsh`

## TODOs

- [ ] Change way prints works, stack.go should return string and cmd should print it.
- [ ] Complete tests.
- [ ] status with `--log` option to show last commit of each branch.