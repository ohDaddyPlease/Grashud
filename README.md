# Grashud
Graceful shutdown, based on callbacks and context

Example:
```go
package main

import (
	"errors"
	"time"

	grashud "github.com/ohdaddyplease/Grashud"
)

func main() {
	g, _, cancel := grashud.New(
        grashud.WithCancel, //must be WithCancel, WithDeadline, WithTimeout
        nil, //must be nil for WithCancel, time.Time for WithDeadline, time.Duration for WithTimeout
    )
	defer cancel() //important

	g.AddFunc(func() error {
		return errors.New("Pew Pew")
	})

	time.Sleep(2 * time.Second)
}

```