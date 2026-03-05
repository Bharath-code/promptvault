package main

import (
	"fmt"
	"os"

	"github.com/Bharath-code/promptvault/internal/cmd"
	"github.com/Bharath-code/promptvault/internal/db"
)

var version = "0.1.0"

func main() {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := cmd.Execute(database); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
