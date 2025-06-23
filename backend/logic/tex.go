package logic

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// DotToTex runs dot2tex CLI on the dot string and returns the TikZ code or error.
func DotToTex(dot string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "dot2tex", "--crop", "-ftikz")
	cmd.Stdin = bytes.NewBufferString(dot)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("dot2tex failed: %v, stderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

// TikZToSVG converts TikZ code to SVG using pdflatex and pdf2svg
func TikZToSVG(tikz string) (string, error) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "tikz2svg")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			// Log the error but don't fail the function since cleanup is not critical
			fmt.Printf("Warning: failed to remove temp directory %s: %v\n", tmpDir, err)
		}
	}()

	// Ensure TikZ has proper document structure
	if !bytes.Contains([]byte(tikz), []byte("\\documentclass")) {
		tikz = fmt.Sprintf(`\documentclass[border=10pt]{standalone}
\usepackage{tikz}
\usetikzlibrary{arrows,automata,positioning}
\begin{document}
%s
\end{document}`, tikz)
	}

	// Write TikZ to file
	texFile := filepath.Join(tmpDir, "input.tex")
	if err := os.WriteFile(texFile, []byte(tikz), 0644); err != nil {
		return "", err
	}

	// Compile with pdflatex
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "pdflatex",
		"-interaction=nonstopmode",
		"-halt-on-error",
		"-output-directory", tmpDir,
		texFile)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdflatex failed: %v, stderr: %s", err, stderr.String())
	}

	// Convert PDF to SVG using pdf2svg
	pdfFile := filepath.Join(tmpDir, "input.pdf")
	svgFile := filepath.Join(tmpDir, "output.svg")

	cmd = exec.CommandContext(ctx, "pdf2svg", pdfFile, svgFile)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdf2svg failed: %v, stderr: %s", err, stderr.String())
	}

	// Read SVG
	svg, err := os.ReadFile(svgFile)
	if err != nil {
		return "", err
	}

	return string(svg), nil
}
