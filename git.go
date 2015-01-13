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

func (c Commit) HumanDate() string {
	return c.Date.Format("2006-01-02 15:04")
}

func Diff(file, hash string) ([]byte, error) {
	var out bytes.Buffer

	git := exec.Command("git", "-C", options.Dir, "show", "--oneline", "--no-color", hash, file)

	// Prune diff stats from output with tail
	tail := exec.Command("tail", "-n", "+8")

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
}

func Commits(filename string, n int) ([]Commit, error) {
	var commits []Commit

	// abbreviated commit hash|author name|author date, strict ISO 8601 format|subject
	logFormat := "--pretty=%h|%an|%at|%s"

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

		unix, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			log.Println("ERROR", err)
		}
		commit.Date = time.Unix(unix, 0)

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
