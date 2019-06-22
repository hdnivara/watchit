package main

import (
	"fmt"

	"watchit/cmd"
)

func main() {
	c := cmd.Parse()
	fmt.Printf("%#v\n", c)
}
