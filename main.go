package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	email := flag.String("email", "", "Email address of the user. Leave empty to get the commits made by any author")
	folder := flag.String("folder", "", "Folder to Scan")
	flag.Parse()

	// get git folder path
	folderPath, err := getFolder(*folder)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(folderPath)
	// get stats for email
	GetStats(*email, folderPath)
}
