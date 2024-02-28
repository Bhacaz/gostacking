#!/bin/zsh

set -e
rm -rf completions
mkdir completions

go run . completion zsh >"completions/gostacking.zsh"
