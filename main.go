package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// TODO: Add a mapping to read store all the commit dates to generate image of contributions
func getStats(email string, repo string) {
	cmd := exec.Command("git", "-C", repo, "log", "--author="+email)
	fmt.Println(cmd.String())
	res, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	chunks := strings.Split(string(res), "\ncommit")
	fmt.Println(chunks[0])
	for _, chunk := range chunks {
		lines := strings.Split(chunk, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Date:") {
				fmt.Println(line)
			}
		}
	}
}

func main() {
	email := flag.String("email", "", "Email address of the user")
	folder := flag.String("folder", "", "Folder to Scan")
	flag.Parse()

	// get git folder path
	folderPath, err := getFolder(*folder)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(folderPath)
	// get stats for email
	getStats(*email, folderPath)
}
