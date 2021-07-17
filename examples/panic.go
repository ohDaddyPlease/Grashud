package main

import (
	"errors"

	grashud "github.com/ohdaddyplease/Grashud"
)

func main() {
	g := grashud.New()
	g.HandleSignals()
	defer g.HandlePanic()

	g.Add(func() error {
		return errors.New("Pew Pew")
	})

	panic("oh no! panic!")
}
