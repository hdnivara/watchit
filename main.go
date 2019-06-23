package main

import (
	"fmt"

	"watchit/command"
)

func main() {
	c := command.Parse()
	fmt.Printf("%#v\n", c)
}
