package main

import (
	"errors"
	"fmt"
	"time"

	grashud "github.com/ohdaddyplease/Grashud"
)

func main() {
	g := grashud.New()
	g.HandleSignals()
	defer g.HandlePanic()

	g.Add(func() error {
		return errors.New("Pew Pew")
	})

	fmt.Println("Push Ctrl+C")

	time.Sleep(5 * time.Second)
}
