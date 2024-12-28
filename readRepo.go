package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TODO: Make it so that it reads through all the subdirectories of a directory and maintains a list of all git repositories on the given folder

func getFolder(folder string) (string, error) {
	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}
	contents, err := f.ReadDir(0)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	for _, content := range contents {
		if content.IsDir() {
			if content.Name() == ".git" {
				path, err := filepath.Abs(strings.TrimSuffix(folder, "/"))
				if err != nil {
					log.Fatal(err)
				}
				return path, nil
			}
		}
	}
	return "", fmt.Errorf("Not a git repository")
}

func genDatesMap(email, repo string, dates map[int]int) map[int]int {
	var cmd *exec.Cmd
	if len(email) != 0 {
		cmd = exec.Command("git", "-C", repo, "log", "--author="+email)
	} else {
		cmd = exec.Command("git", "-C", repo, "log", "--author="+email)
	}
	res, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	chunks := strings.Split(string(res), "\ncommit")
	ref := "Mon Jan 2 15:04:05 2006 -0700"

	for _, chunk := range chunks {
		lines := strings.Split(chunk, "\n")
		curr := time.Now()
		curr = time.Date(curr.Year(), curr.Month(), curr.Day(), 0, 0, 0, 0, curr.Location())
		for _, line := range lines {
			if strings.HasPrefix(line, "Date:") {
				t, err := time.Parse(ref, strings.TrimPrefix(line, "Date:   "))
				if err != nil {
					log.Fatal(err)
				}
				t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
				index := int(math.Round((curr.Sub(t).Hours()) / 24))
				if index <= 183 {
					dates[index]++
				}
			}
		}
	}
	return dates
}
