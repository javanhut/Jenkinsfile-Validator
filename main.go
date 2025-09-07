package main

import (
	"fmt"

	"github.com/javanhut/Jenkinsfile-Validator/cli"
)

func main() {
	print := fmt.Println
	print("Jenkinsfile Validator")
	cli.Execute()
}
