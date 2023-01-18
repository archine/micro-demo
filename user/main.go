package main

import (
	"micro-demo/user/base"
	_ "micro-demo/user/controller"
	"os"
)

//go:generate go run main.go ast
func main() {
	base.RunApplication(os.Args)
}
