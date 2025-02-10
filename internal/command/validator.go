package command

import (
	"log"
	"os/exec"
)

// Deprecated: CheckSSH is deprecated
func CheckSSH() {
	_, err1 := exec.LookPath("ssh")
	_, err2 := exec.LookPath("tsh")

	if err1 != nil && err2 != nil {
		log.Fatal("ssh && tsh is not installed")
	}

	if err1 != nil {
		log.Fatal("[W] ssh is not installed")
	}

	if err2 != nil {
		log.Fatal("[W] tsh is not installed")
	}
}
