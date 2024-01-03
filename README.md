# gostacking

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
* `switch`      Change the current stack.
* `sync`        Merge all branch in a stack into the current branch.

## Notes

https://cobra.dev/

```bash
~/go/bin/cobra-cli add new
go mod tidy
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
