package main

import "simplest-shortener/internal"

func main() {
	app := internal.NewApp()
	app.Run()
}
