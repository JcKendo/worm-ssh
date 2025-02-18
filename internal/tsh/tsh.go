package tsh

import (
	"fmt"
	"github.com/JcKendo/worm-ssh/internal/config"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func GenerateCommandArgs(c config.SSHConfig) []string {
	user := "root"

	if c.User != "" {
		user = c.User
	}

	return strings.Split(fmt.Sprintf("ssh %s@ip=%s", user, c.Host), " ")
}

func Run(args []string) {
	args = slices.DeleteFunc(args, func(s string) bool { return s == "" })
	fmt.Println(args)
	cmd := exec.Command("tsh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
