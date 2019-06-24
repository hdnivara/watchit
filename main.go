package main

import (
	"fmt"
	"log"

	"watchit/command"
	"watchit/watch"
)

func handler(op watch.Op, fileName string) {
	fmt.Printf("op=%s file=%s", op.Name(), fileName)
}

func main() {
	c := command.Parse()
	fmt.Printf("%#v\n", c)

	w := watch.New(c.Dirs, c.Cmds, c.Regex, c.Recursive, handler)

	if err := w.Setup(); err != nil {
		log.Fatalln("setting up watch failed")
	}

	if err := w.Start(); err != nil {
		log.Fatalln("starting watch failed")
	}
}
