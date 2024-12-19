package main

import (
	"fmt"
	"os"

	repl "github.com/kev/repl"
)

func main() {
	kevBanner()
	repl.Start(os.Stdin, os.Stdout)
}

func kevBanner() {
	banner, err := os.ReadFile("kev-banner.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(banner))
}
