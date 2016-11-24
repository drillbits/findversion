package main

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

func findversion(cli *CLI, pkgroot, prefix string) (string, int) {
	v, err := findByFile(pkgroot)
	if err != nil {
		return "", cli.Fatalf("%v", err)
	} else if v != "" {
		return v, ExitCodeOK
	}

	if !isGitRepo(pkgroot) {
		return "", cli.Fatalf("FAILED: not a Git repo; must put a VERSION file in %s", pkgroot)
	}

	branch, err := currentBranch(pkgroot)
	if err != nil {
		return "", cli.Fatalf("%v", err)
	}

	tag := "devel"
	if strings.HasPrefix(branch, prefix) {
		tag = findBranchClosestTag(pkgroot, branch)
		if tag == "" {
			tag = branch
		}
	}

	hd, err := getHashAndDate(pkgroot)
	if err != nil {
		return "", cli.Fatalf("%v", err)
	}
	v = tag + hd

	return v, ExitCodeOK
}

func findByFile(pkgroot string) (string, error) {
	path := fmtPath("%s/VERSION", pkgroot)
	if isFile(path) {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		s := strings.Trim(string(b), " \t\r\n")
		if s != "" {
			return s, nil
		}
	}
	return "", nil
}

func isGitRepo(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	gitDir := chomp(string(out))
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(dir, gitDir)
	}

	return isDir(gitDir)
}

func currentBranch(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return chomp(string(out)), nil
}

func gitLogOnlyRefNames(dir, branch string) ([]byte, error) {
	cmd := exec.Command("git", "log", "--decorate=full", "--format=format:%D", "master.."+branch)
	cmd.Dir = dir
	return cmd.Output()
}

func determineClosestTag(b []byte) string {
	for _, line := range strings.Split(string(b), "\n") {
		for _, name := range strings.Split(line, ",") {
			name = strings.Trim(name, " ")
			i := strings.Index(name, "refs/tags/")
			if i > 0 {
				return name[i+len("refs/tags/"):]
			}
		}
	}
	return ""
}

func findBranchClosestTag(dir, branch string) string {
	b, _ := gitLogOnlyRefNames(dir, branch)
	return determineClosestTag(b)
}

func getHashAndDate(dir string) (string, error) {
	cmd := exec.Command("git", "log", "-n", "1", "--format=format: +%h %cd", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return chomp(string(out)), nil
}
