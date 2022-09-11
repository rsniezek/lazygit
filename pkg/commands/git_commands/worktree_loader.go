package git_commands

import (
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type WorktreeLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewWorktreeLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *WorktreeLoader {
	return &WorktreeLoader{
		Common: common,
		cmd:    cmd,
	}
}

func (self *WorktreeLoader) GetWorktrees() ([]*models.Worktree, error) {
	cmdArgs := NewGitCmd("worktree").Arg("list", "--porcelain", "-z").ToArgv()
	worktreesOutput, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	splitLines := strings.Split(worktreesOutput, "\x00")

	var worktrees []*models.Worktree
	var currentWorktree *models.Worktree
	for _, splitLine := range splitLines {
		if len(splitLine) == 0 && currentWorktree != nil {
			worktrees = append(worktrees, currentWorktree)
			currentWorktree = nil
			continue
		}
		if strings.HasPrefix(splitLine, "worktree ") {
			path := strings.SplitN(splitLine, " ", 2)[1]
			currentWorktree = &models.Worktree{
				Id:   len(worktrees),
				Path: path,
			}
		} else if strings.HasPrefix(splitLine, "branch ") {
			branch := strings.SplitN(splitLine, " ", 2)[1]
			currentWorktree.Branch = filepath.Base(branch)
		}
	}

	return worktrees, nil
}