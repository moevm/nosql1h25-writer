package main

import (
	"os"

	"github.com/moevm/nosql1h25-writer/backend/internal/app"
)

func main() {
	app := app.New(os.Getenv("CONFIG_PATH"))
	app.Start()
}
