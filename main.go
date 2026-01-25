package main

import (
	"fmt"
	"os"

	repl "github.com/kev/repl"
)

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "run":
			if len(os.Args) < 3 {
				fmt.Println("Usage: kev run <file>")
				os.Exit(1)
			}
			fileName := os.Args[2]
			repl.RunFile(fileName)
		default:
			fmt.Println("Usage: kev run <file>")
			os.Exit(1)
		}
	} else {
		kevBanner()
		repl.RunPrompt(os.Stdin, os.Stdout)
	}
}

func kevBanner() {
	banner, err := os.ReadFile("kev-banner.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(banner))
}
