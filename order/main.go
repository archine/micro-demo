package main

import (
	"micro-demo/order/base"
	_ "micro-demo/order/controller"
	"os"
)

//go:generate go run main.go ast
func main() {
	base.RunApplication(os.Args)
}
