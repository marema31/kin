//go:generate bin/pkger

package main

import (
	"github.com/marema31/kin/cmd"
	"github.com/markbates/pkger"
)

func main() {
	pkger.Include("/site")
	cmd.Execute()
}
