package main

import (
	"log"

	"github.com/demisang/codegen/cmd/codegen"
)

func main() {
	err := codegen.Run()
	if err != nil {
		log.Fatal(err)
	}
}
