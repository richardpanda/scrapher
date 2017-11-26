package main

import "fmt"

func printError(e <-chan error) {
	for err := range e {
		fmt.Println(err)
	}
}
