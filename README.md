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

# Update
brew update && brew upgrade gostacking
```

## Usage

```bash
$ gostacking [command]
```

### Available Commands

```
add         Add a branch to the current stack
checkout    Checkout a branch from a stack
delete      Delete a gostacking
help        Help about any command
list        List all stacks
new         Create a new gostacking
remove      Remove a branch from the current stack
status      Get current stack
switch      Change the current stack
sync        Merge all branches into the others
```

## Example

```bash
git checkout -b feature/1

gostacking new my-stack

gostacking list
# Current stack: my-stack
# 1. my-stack

gostacking add
# Branch feature/1 added to stack my-stack

gostacking status
# my-stack
# 1. feature/1

git checkout -b feature/2
gostacking add feature/2
# Branch feature/2 added to stack my-stack

gostacking status
# my-stack
# 1. feature/1
# 2. feature/2

gostacking checkout 1

touch file1.txt
git add file1.txt && git commit -m "Add file1.txt"

gostacking status
# my-stack
# 1. feature/1
# 2. feature/2 *

gostacking sync
# Syncing my-stack
# Fetching ...
# Branch: feature/1
#     Checkout ...
#     Pull ...
# Branch: feature/2
#     Checkout ...
#     Pull ...
#     Merging feature/1

gostacking status --log
# my-stack
# 1. feature/1
#      Add file1.txt - 21e656719d - 3 minute ago
# 2. feature/2
#      Merge feature/1 into feature/2 (gostacking) - f8178d7384 - 1 minute ago
```

## Release

1. Update the version in file `VERSION`
2. `zsh scripts/release.zsh`

## Notes

https://cobra.dev/

```bash
~/go/bin/cobra-cli add new
go mod tidy
go build -ldflags="-s -w" -o dist
goreleaser release --snapshot --clean
```

## TODOs

- [ ] Add flag `--position` (`-p`) to `add` **command** to add a branch at a specific position.
- [ ] Add completion suggestion list of branches (with a max).
- [ ] Add option to merge HEAD into the first branch.
    - `--merge-head` (`-s`) to `sync` **command** to merge the HEAD branch into the first branch.
    - Detect default HEAD branch `git symbolic-ref refs/remotes/origin/HEAD --short | sed "s@^$(git remote)/@@"`. If it is not found `git remote set-head origin main`
- [ ] Change way prints works, stack.go should return string and cmd should print it.
- [ ] Complete tests.
