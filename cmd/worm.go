package cmd

import (
	"github.com/JcKendo/worm/internal/command"
	"github.com/JcKendo/worm/internal/config"
	"github.com/JcKendo/worm/internal/history"
	"github.com/JcKendo/worm/internal/interactive"
	"github.com/JcKendo/worm/internal/ssh"
	"github.com/JcKendo/worm/internal/tsh"
	"github.com/JcKendo/worm/internal/workspace"
	"os"
)

func Main() {
	args := os.Args[1:]
	mode := config.TSHMode
	action, value := command.Which()
	switch action {
	case command.InteractiveHistory:
		args, mode = interactive.History()
	case command.InteractiveConfig:
		args, mode = interactive.Config("")
	case command.InteractiveConfigWithSearch:
		args, mode = interactive.Config(value)
	case command.ListHistory:
		history.Print()
		return
	case command.ListConfig:
		config.Print()
		return
	case command.ListWorkspace:
		workspace.Print()
		return
	case command.ActiveWorkspace:
		interactive.Active()
		return
	default:

	}
	history.AddHistoryFromArgs(args, mode)

	if mode == config.SSHMode {
		ssh.Run(args)
	} else {
		tsh.Run(args)
	}
}
