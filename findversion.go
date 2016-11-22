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

	var precise bool
	if strings.HasPrefix(branch, prefix) {
		v, precise = findBranchClosestTag(pkgroot, branch)
	}

	if !precise {
		hd, err := getHashAndDate(pkgroot)
		if err != nil {
			return "", cli.Fatalf("%v", err)
		}
		v = branch + hd
	}

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

func findBranchClosestTag(dir, branch string) (tag string, precise bool) {
	cmd := exec.Command("git", "log", "--decorate=full", "--format=format:%d", "master.."+branch)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return
	}

	tag = branch
	for _, line := range strings.SplitAfter(string(out), "\n") {
		// Each line is either blank, or looks like
		//	  (tag: refs/tags/0.1.0, refs/remotes/origin/release-branch.0.1.0, refs/heads/release-branch.0.1.0)
		// We need to find an element starting with refs/tags/.
		i := strings.Index(line, " refs/tags/")
		if i < 0 {
			continue
		}
		i += len(" refs/tags/")
		// The tag name ends at a comma or paren (prefer the first).
		j := strings.Index(line[i:], ",")
		if j < 0 {
			j = strings.Index(line[i:], ")")
		}
		if j < 0 {
			continue // malformed line; ignore it
		}
		tag = line[i : i+j]
		if i == 0 {
			precise = true // tag denotes HEAD
		}
		break
	}
	return
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
