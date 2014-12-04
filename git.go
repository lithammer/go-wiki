package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Author  string
	Date    time.Time
	File    string
	Hash    string
	Subject string
}

func (c Commit) Diff() ([]byte, error) {
	return Diff(c.File, c.Hash)
}

func (c Commit) FileNoExt() string {
	return strings.TrimSuffix(c.File, filepath.Ext(c.File))
}

func Diff(filepath, hash string) ([]byte, error) {
	var out bytes.Buffer

	// diff := fmt.Sprintf("%s^..%s", hash, hash)

	git := exec.Command("git", "-C", options.Dir, "diff", hash+"^", hash, "--no-color", "--", filepath)

	// Prune diff stats from output with tail
	tail := exec.Command("tail", "-n", "+6")

	var err error
	tail.Stdin, err = git.StdoutPipe()
	if err != nil {
		log.Println("ERROR", err)
	}

	tail.Stdout = &out

	err = tail.Start()
	if err != nil {
		log.Println("ERROR", err)
	}

	err = git.Run()
	if err != nil {
		log.Println("ERROR", err)
	}

	err = tail.Wait()
	if err != nil {
		log.Println("ERROR", err)
	}

	return out.Bytes(), err

	// Remove stat part from diff, example:
	//
	// diff --git foo.md foo.md
	// index 74c109a..8c55778 100644
	// --- foo.md
	// +++ foo.md
	// @@ -1,3 +1,6 @@
	//  # Foo
	//  ## Add more foo!
	// +- this is a list
	// +- of foo stuff
	// separator := []byte("@@\n")
	// diff := bytes.SplitAfterN(out.Bytes(), separator, 2)[1]

	// return diff, err
}

func Commits(filename string, n int) ([]Commit, error) {
	var commits []Commit

	// abbreviated commit hash|author name|author date, strict ISO 8601 format|subject
	logFormat := "--pretty=%h|%an|%aI|%s"

	cmd := exec.Command("git", "-C", options.Dir, "log", "-n", strconv.Itoa(n), logFormat, filename)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("ERROR", err)
		return commits, err
	}

	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		log.Println("ERROR", err)
		return commits, err
	}

	out := bufio.NewScanner(stdout)
	for out.Scan() {
		fields := strings.Split(out.Text(), "|")

		commit := Commit{
			Author:  fields[1],
			File:    filename,
			Hash:    fields[0],
			Subject: fields[3],
		}

		commit.Date, err = time.Parse(time.RFC3339Nano, fields[2])
		if err != nil {
			log.Println("ERROR", err)
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// Check if a path contains a Git repository
func IsGitRepository(path string) bool {
	var out bytes.Buffer
	cmd := exec.Command("git", "-C", options.Dir, "rev-parse", "--is-inside-work-tree")
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Println("ERROR", err)
		return false
	}

	var val bool
	_, err = fmt.Sscanf(out.String(), "%t", &val)
	if err != nil {
		log.Println("ERROR", err)
		return false
	}

	return val
}
