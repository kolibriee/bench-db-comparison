package main

import "github.com/kolibriee/bench-db-comparison/app"

const (
	ConfigPath = "./"
	FileName   = "config"
)

func main() {
	app.Run(ConfigPath, FileName)
}
