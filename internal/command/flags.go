package command

import (
	"os"
)

type Action int

const (
	PassThrough Action = iota
	InteractiveHistory
	InteractiveConfig
	InteractiveConfigWithSearch
	ListHistory
	ListConfig
	ListWorkspace
	ActiveWorkspace
)

func Which() (Action, string) {
	if len(os.Args) == 1 {
		return InteractiveHistory, ""
	}

	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--history":
			return ListHistory, ""
		case "--config":
			return ListConfig, ""
		case "--workspace":
			return ListWorkspace, ""
		case "--active":
			return ActiveWorkspace, ""
		case "-":
			return InteractiveConfig, ""
		}
	}

	if len(os.Args) == 3 {
		if os.Args[1] == "-" {
			return InteractiveConfigWithSearch, os.Args[2]
		}
	}

	return PassThrough, ""
}
