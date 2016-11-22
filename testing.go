package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func tempGitRepo(dir, prefix string) (string, func(), error) {
	gitDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return "", func() {}, err
	}

	f, err := ioutil.TempFile(gitDir, prefix)
	if err != nil {
		return "", func() {
			os.Remove(gitDir)
		}, nil
	}

	close := func() {
		os.Remove(gitDir)
		os.Remove(f.Name())
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = gitDir
	_, err = cmd.Output()
	if err != nil {
		return "", close, err
	}

	cmd = exec.Command("git", "add", "-A")
	cmd.Dir = gitDir
	_, err = cmd.Output()
	if err != nil {
		return "", close, err
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = gitDir
	_, err = cmd.Output()
	if err != nil {
		return "", close, err
	}

	return gitDir, close, nil
}

func tempGitCommit(dir, prefix string) (string, error) {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return "", err
	}
	name := f.Name()

	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = dir
	_, err = cmd.Output()
	if err != nil {
		return name, err
	}

	msg := fmt.Sprintf("test commit %s", name)
	cmd = exec.Command("git", "commit", "-m", msg)
	cmd.Dir = dir
	_, err = cmd.Output()
	if err != nil {
		return name, err
	}

	return name, err
}
