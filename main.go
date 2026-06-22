package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:embed scripts/entrypoint.sh
var scripts embed.FS

func main() {
	image := flag.String("image", "node:lts", "OCI image to run")
	repo := flag.String("repo", "", "Repo path to mount at /workspace (default: current dir)")
	flag.Parse()

	if *repo == "" {
		wd, err := os.Getwd()
		if err != nil {
			die("failed to get current working directory: %v", err)
		}
		*repo = wd
	}

	if err := checkSmolvm(); err != nil {
		die("smolvm check failed: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "claude-sandbox-*")
	if err != nil {
		die("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	scriptDir := filepath.Join(tmpDir, "scripts")
	if err := os.Mkdir(scriptDir, 0755); err != nil {
		die("failed to create scripts dir: %v", err)
	}

	entrypointSrc, err := scripts.Open("scripts/entrypoint.sh")
	if err != nil {
		die("failed to open embedded entrypoint: %v", err)
	}
	defer entrypointSrc.Close()

	entrypointDst, err := os.Create(filepath.Join(scriptDir, "entrypoint.sh"))
	if err != nil {
		die("failed to create entrypoint file: %v", err)
	}
	defer entrypointDst.Close()

	if _, err := io.Copy(entrypointDst, entrypointSrc); err != nil {
		die("failed to copy entrypoint: %v", err)
	}

	if err := os.Chmod(filepath.Join(scriptDir, "entrypoint.sh"), 0755); err != nil {
		die("failed to chmod entrypoint: %v", err)
	}

	cmd := exec.Command("smolvm", "machine", "run",
		"--net", "-it",
		"--image", *image,
		"-v", *repo+":/workspace",
		"-v", scriptDir+":/sandbox",
		"--", "/bin/sh", "/sandbox/entrypoint.sh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		die("smolvm execution failed: %v", err)
	}
}

func checkSmolvm() error {
	cmd := exec.Command("command", "-v", "smolvm")
	if runtime.GOOS != "windows" {
		cmd = exec.Command("sh", "-c", "command -v smolvm")
	}

	if err := cmd.Run(); err != nil {
		msg := "smolvm not found. Install with:\n"
		msg += "  curl -sSL https://smolmachines.com/install.sh | bash\n"
		msg += "See: https://github.com/smol-machines/smolvm"
		return fmt.Errorf(msg)
	}
	return nil
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
