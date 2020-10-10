package main

import (
	"github.com/pacoorozco/networkdevices/internal/app"
)

func main() {
	a := app.App{}
	a.Initialize()

	a.Run(":8010")
}
