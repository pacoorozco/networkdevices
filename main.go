package main

import "github.com/pacoorozco/networkdevices/internal"

func main() {
	a := App{}
	a.Initialize()

	a.Run(":8010")
}
