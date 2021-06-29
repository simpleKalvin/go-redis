package main

import "go_redis/src/cli"

func main() {
	userInput := cli.NewRedis()
	userInput.GET()
}
