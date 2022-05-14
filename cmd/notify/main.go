package main

import (
	"notifier/pkg/notifier"
)

func main() {
	notifier.Notify("http://127.0.0.1:4567", "Hello, world!")
}
