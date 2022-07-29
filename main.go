package main

import (
	"fmt"
	"os"
)

func main() {
	lot := NewParkingLot()
	repl := NewREPL(lot)
	args := os.Args[1:]

	if len(args) == 0 {
		repl.Run(os.Stdin, os.Stdout)
	} else {
		fmt.Println(args[0])
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", args[0], err)
			os.Exit(1)
		}
		defer f.Close()
		repl.Run(f, os.Stdout)
	}
}
