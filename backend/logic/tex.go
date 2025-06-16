package logic

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// DotToTex runs dot2tex CLI on the dot string and returns the TikZ code or error.
func DotToTex(dot string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "dot2tex", "-ftikz")
	cmd.Stdin = bytes.NewBufferString(dot)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
