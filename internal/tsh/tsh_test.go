package tsh

import (
	"github.com/JcKendo/worm/internal/config"
	"testing"
)

func TestGenerateCommandArgs(t *testing.T) {
	c := config.SSHConfig{
		User: "testuser",
		Host: "testhost",
	}
	expected := []string{"ssh", "testuser@ip=testhost"}
	args := GenerateCommandArgs(c)
	if !equal(args, expected) {
		t.Errorf("expected %v, got %v", expected, args)
	}

}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
