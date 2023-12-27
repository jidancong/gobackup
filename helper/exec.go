package helper

import (
	"fmt"
	"log/slog"
	"os/exec"
)

type Helper struct {
	pwd string
}

func NewHelper(pwd string) *Helper {
	return &Helper{pwd: pwd}
}

func (h *Helper) Exec(command string, args ...string) (string, error) {
	c := exec.Command(command, args...)
	c.Dir = h.pwd
	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}
	slog.Debug("执行命令", command, args)
	if err != nil {
		return "", fmt.Errorf("error: %v %s", err, string(out))
	}
	return string(out), nil
}
