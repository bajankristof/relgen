package main

import (
	"github.com/bajankristof/relgen/cmd"
)

func main() {
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}
