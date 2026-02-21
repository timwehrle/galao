package main

import (
	"context"
	"fmt"

	"github.com/timwehrle/galao"
)

func main() {
	app := galao.New()

	count := 0

	render := func() {
		if err := app.SetView(galao.VStack(
			galao.Text(fmt.Sprintf("Count: %d", count)),
			galao.Button("inc", "Increment"),
			galao.Button("dec", "Decrement"),
		)); err != nil {
			panic(err)
		}
	}

	app.OnEvent("inc", func(e galao.Event) {
		count++
		render()
	})

	app.OnEvent("dec", func(e galao.Event) {
		count--
		render()
	})

	if err := app.Run(context.Background(), func() {
		render()
	}); err != nil {
		panic(err)
	}
}
