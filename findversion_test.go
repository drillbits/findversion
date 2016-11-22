package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestFindByFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	pkgroot := filepath.Join(wd, "testdata", "v8")
	v, err := findByFile(pkgroot)
	if err != nil {
		t.Fatal(err)
	}

	if v != "v8!" {
		t.Errorf("expected: v8!, but got: %v", v)
	}

	pkgroot = filepath.Join(wd, "testdata", "vd")
	v, err = findByFile(pkgroot)
	if err != nil {
		t.Fatal(err)
	}

	if v != "" {
		t.Errorf("expected: empty, but got: %v", v)
	}

	pkgroot = filepath.Join(wd, "testdata")
	v, err = findByFile(pkgroot)
	if err != nil {
		t.Fatal(err)
	}

	if v != "" {
		t.Errorf("expected: empty, but got: %v", v)
	}
}

func TestIsGitRepo(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	exp := true
	got := isGitRepo(dir)
	if got != exp {
		t.Errorf("expected: %v, but got: %v", exp, got)
	}

	dir, err = ioutil.TempDir("", "TestIsGitRepo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)

	exp = false
	got = isGitRepo(dir)
	if got != exp {
		t.Errorf("expected: %v, but got: %v", exp, got)
	}
}

func TestCurrentBranch(t *testing.T) {
	dir, close, err := tempGitRepo("", "TestCurrentBranch")
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	cmd := exec.Command("git", "checkout", "-b", "test-branch")
	cmd.Dir = dir
	_, err = cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	branch, err := currentBranch(dir)
	if err != nil {
		t.Fatal(err)
	}

	if branch != "test-branch" {
		t.Errorf("expected: test-branch, but got: %v", branch)
	}
}

func TestFindBranchClosestTag(t *testing.T) {
	dir, close, err := tempGitRepo("", "TestFindBranchClosestTag")
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	rb := "release-branch.0.1.0"
	v := "v0.1.0"

	cmd := exec.Command("git", "checkout", "-b", rb)
	cmd.Dir = dir
	_, err = cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	tag, precise := findBranchClosestTag(dir, rb)
	if err != nil {
		t.Fatal(err)
	}
	if precise {
		t.Errorf("no precise expected: %v", precise)
	}
	if tag != rb {
		t.Errorf("expected: %s, but got: %v", rb, tag)
	}

	tempfile, err := tempGitCommit(dir, "TestFindBranchClosestTag")
	if err != nil {
		t.Fatal(err)
	} else if tempfile != "" {
		defer os.Remove(tempfile)
	}
	cmd = exec.Command("git", "tag", v)
	cmd.Dir = dir
	_, err = cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	tag, precise = findBranchClosestTag(dir, rb)
	if err != nil {
		t.Fatal(err)
	}
	if precise {
		t.Errorf("no precise expected: %v", precise)
	}
	if tag != v {
		t.Errorf("expected: %v, but got: %v", v, tag)
	}
}

func TestGetHashAndDate(t *testing.T) {
	dir, close, err := tempGitRepo("", "TestGetHashAndDate")
	if err != nil {
		t.Fatal(err)
	}
	defer close()
	exp := " +"
	ref := filepath.Join(dir, ".git/refs/heads/master")
	b, err := ioutil.ReadFile(ref)
	if err != nil {
		t.Fatal(err)
	}
	exp += string(b)[:7]
	fi, err := os.Stat(ref)
	if err != nil {
		t.Fatal(err)
	}
	exp += fi.ModTime().Format(" Mon Jan 02 15:04:05 2006 -0700")

	got, err := getHashAndDate(dir)
	if err != nil {
		t.Fatal(err)
	}

	if got != exp {
		t.Errorf("expected: %v, but got: %v", exp, got)
	}
}
