package main

import (
	"github.com/ramessesii2/depviz/pkg/cmd/depviz"
)

func main() {
	err := depviz.Root().Execute()
	if err != nil {
		panic(err)
	}
}
