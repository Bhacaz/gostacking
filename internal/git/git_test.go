package git

import (
    "testing"
    "errors"
)

type GitExecutorStub struct {
    ExecuteStub func(string) (string, error)
}

func (gee GitExecutorStub) ExecuteGitCommand(command string) (string, error) {
    return gee.ExecuteStub(command)
}

func TestCurrentBranchName(t *testing.T) {
    t.Run("Get current branch name", func(t *testing.T) {
         gitCmd := GitCommands{
                executor: GitExecutorStub{
                    ExecuteStub: func(string) (string, error) {
                        return "master", nil
                    },
                },
            }

            result, err := gitCmd.CurrentBranchName()

            want := "master"
            if result != want {
                t.Errorf("got %s, want %s", result, want)
            }

            if err != nil {
                t.Errorf("got Error %s, want none", err)
            }
    })

    t.Run("Error", func(t *testing.T) {
        gitCmd := GitCommands{
            executor: GitExecutorStub{
                ExecuteStub: func(string) (string, error) {
                    return "", errors.New("git command error")
                },
            },
        }

        _, err := gitCmd.CurrentBranchName()

        if err == nil {
            t.Errorf("got none, want Error")
        }
    })
}
