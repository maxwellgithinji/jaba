package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/maxwellgithinji/jaba/pkg/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hi %s! Welcome to jaba programming language\n", user.Username)
	fmt.Println("Enter the jaba program below:")
	repl.Run(os.Stdin, os.Stdout)

}
